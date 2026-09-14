package logic

import (
	"context"

	"backend/internal/errorx"
	"backend/internal/model"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DisableWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDisableWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DisableWorkflowLogic {
	return &DisableWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DisableWorkflowLogic) DisableWorkflow(req *types.DisableWorkflowReq) (resp *types.DisableWorkflowResp, err error) {
	workflow, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errorx.NewNotFoundError("workflow not found")
		}
		l.Logger.Errorf("find workflow failed: %v", err)
		return nil, errorx.NewInternalError("failed to find workflow")
	}

	if workflow.Status != "enabled" {
		return nil, errorx.NewBadRequestError("can only disable enabled workflows; current status: " + workflow.Status)
	}

	err = l.svcCtx.WorkflowModel.UpdateStatus(l.ctx, uint64(req.Id), "disabled")
	if err != nil {
		l.Logger.Errorf("disable workflow failed: %v", err)
		return nil, errorx.NewInternalError("failed to disable workflow")
	}

	return &types.DisableWorkflowResp{
		Status: "disabled",
	}, nil
}
