// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package handler

import (
	"fmt"
	"net/http"
	"time"

	"backend/internal/sse"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRunEventsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRunEventsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		run, err := svcCtx.RunModel.FindOne(r.Context(), uint64(req.RunId))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		hub := sse.GetHub()
		sub := hub.Subscribe(int64(run.Id))
		defer hub.Unsubscribe(sub)

		logx.Infof("SSE client subscribed to run %d", run.Id)

		keepAlive := time.NewTicker(15 * time.Second)
		defer keepAlive.Stop()

		for {
			select {
			case <-r.Context().Done():
				logx.Infof("SSE client disconnected from run %d", run.Id)
				return

			case event := <-sub.EventCh:
				data, err := event.ToSSEData()
				if err != nil {
					logx.Errorf("Failed to serialize event: %v", err)
					continue
				}
				fmt.Fprintf(w, "event: %s\n", event.EventType)
				fmt.Fprintf(w, "id: %s\n", event.EventId)
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()

			case <-keepAlive.C:
				fmt.Fprintf(w, ": keepalive\n\n")
				flusher.Flush()
			}
		}
	}
}
