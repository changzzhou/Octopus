package logic

import (
	"context"

	"backend/internal/errorx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWorkflowsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListWorkflowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkflowsLogic {
	return &ListWorkflowsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListWorkflowsLogic) ListWorkflows(req *types.ListWorkflowsReq) (resp *types.ListWorkflowsResp, err error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	total, err := l.svcCtx.WorkflowModel.Count(l.ctx, req.Status)
	if err != nil {
		l.Logger.Errorf("count workflows failed: %v", err)
		return nil, errorx.NewInternalError("failed to count workflows")
	}

	workflows, err := l.svcCtx.WorkflowModel.FindByPage(l.ctx, page, pageSize, req.Status)
	if err != nil {
		l.Logger.Errorf("list workflows failed: %v", err)
		return nil, errorx.NewInternalError("failed to list workflows")
	}

	summaries := make([]types.WorkflowSummary, 0, len(workflows))
	for _, w := range workflows {
		summaries = append(summaries, types.WorkflowSummary{
			Id:          int64(w.Id),
			Name:        w.Name,
			Description: w.Description.String,
			Status:      w.Status,
			Version:     int(w.Version),
			CreatedAt:   w.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   w.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &types.ListWorkflowsResp{
		Total:     total,
		Workflows: summaries,
	}, nil
}
