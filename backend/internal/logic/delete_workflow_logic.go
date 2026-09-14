// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete workflow
func NewDeleteWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkflowLogic {
	return &DeleteWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteWorkflowLogic) DeleteWorkflow(req *types.DeleteWorkflowReq) (resp *types.DeleteWorkflowResp, err error) {
	// Stub: no-op - full implementation in BE-1
	return &types.DeleteWorkflowResp{}, nil
}
