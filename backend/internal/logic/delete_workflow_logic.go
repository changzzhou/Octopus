package logic

import (
	"context"

	"backend/internal/errorx"
	"backend/internal/model"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkflowLogic {
	return &DeleteWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteWorkflowLogic) DeleteWorkflow(req *types.DeleteWorkflowReq) (resp *types.DeleteWorkflowResp, err error) {
	workflow, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errorx.NewNotFoundError("workflow not found")
		}
		l.Logger.Errorf("find workflow failed: %v", err)
		return nil, errorx.NewInternalError("failed to find workflow")
	}

	if workflow.Status != "draft" {
		return nil, errorx.ErrDeleteNonDraft
	}

	err = l.svcCtx.WorkflowModel.Delete(l.ctx, uint64(req.Id))
	if err != nil {
		l.Logger.Errorf("delete workflow failed: %v", err)
		return nil, errorx.NewInternalError("failed to delete workflow")
	}

	return &types.DeleteWorkflowResp{}, nil
}
