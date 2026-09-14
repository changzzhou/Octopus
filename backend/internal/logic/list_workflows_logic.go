// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWorkflowsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// List all workflows
func NewListWorkflowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkflowsLogic {
	return &ListWorkflowsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListWorkflowsLogic) ListWorkflows(req *types.ListWorkflowsReq) (resp *types.ListWorkflowsResp, err error) {
	// Stub: returns empty list - full implementation in BE-1
	return &types.ListWorkflowsResp{
		Total:     0,
		Workflows: []types.Workflow{},
	}, nil
}
