// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package handler

import (
	"net/http"

	"backend/internal/logic"
	"backend/internal/svc"
	"backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Update workflow
func UpdateWorkflowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateWorkflowReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUpdateWorkflowLogic(r.Context(), svcCtx)
		resp, err := l.UpdateWorkflow(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
