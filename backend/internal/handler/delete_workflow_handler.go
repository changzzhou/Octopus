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

// Delete workflow
func DeleteWorkflowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteWorkflowReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewDeleteWorkflowLogic(r.Context(), svcCtx)
		resp, err := l.DeleteWorkflow(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
