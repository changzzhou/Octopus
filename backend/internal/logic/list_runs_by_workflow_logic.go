// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRunsByWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// List runs by workflow
func NewListRunsByWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRunsByWorkflowLogic {
	return &ListRunsByWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRunsByWorkflowLogic) ListRunsByWorkflow(req *types.ListRunsByWorkflowReq) (resp *types.ListRunsByWorkflowResp, err error) {
	runs, err := l.svcCtx.RunModel.FindByWorkflowId(l.ctx, uint64(req.Id))
	if err != nil {
		return nil, err
	}

	runSummaries := make([]types.RunSummary, 0, len(runs))
	for _, run := range runs {
		startedAt := ""
		if run.StartedAt.Valid {
			startedAt = run.StartedAt.Time.Format("2006-01-02T15:04:05Z")
		}

		finishedAt := ""
		if run.FinishedAt.Valid {
			finishedAt = run.FinishedAt.Time.Format("2006-01-02T15:04:05Z")
		}

		errorMessage := ""
		if run.ErrorMessage.Valid {
			errorMessage = run.ErrorMessage.String
		}

		runSummaries = append(runSummaries, types.RunSummary{
			Id:              int64(run.Id),
			WorkflowId:      int64(run.WorkflowId),
			WorkflowVersion: int(run.WorkflowVersion),
			Status:          run.Status,
			TriggerType:     run.TriggerType,
			StartedAt:       startedAt,
			FinishedAt:      finishedAt,
			ErrorMessage:    errorMessage,
			CreatedAt:       run.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:       run.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &types.ListRunsByWorkflowResp{
		Runs: runSummaries,
	}, nil
}
