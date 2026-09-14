// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"encoding/json"

	"backend/internal/errorx"
	"backend/internal/model"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRunLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get run details
func NewGetRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRunLogic {
	return &GetRunLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRunLogic) GetRun(req *types.GetRunReq) (resp *types.GetRunResp, err error) {
	run, err := l.svcCtx.RunModel.FindOne(l.ctx, uint64(req.RunId))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errorx.NewCodeError(404, "run not found")
		}
		return nil, err
	}

	var def *types.WorkflowDefinitionSnapshot
	if run.DefinitionSnapshot != "" {
		var parsedDef types.WorkflowDefinitionSnapshot
		if err := json.Unmarshal([]byte(run.DefinitionSnapshot), &parsedDef); err == nil {
			def = &parsedDef
		}
	}

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

	return &types.GetRunResp{
		Run: types.Run{
			Id:                 int64(run.Id),
			WorkflowId:         int64(run.WorkflowId),
			WorkflowVersion:    int(run.WorkflowVersion),
			DefinitionSnapshot: def,
			Status:             run.Status,
			TriggerType:        run.TriggerType,
			StartedAt:          startedAt,
			FinishedAt:         finishedAt,
			ErrorMessage:       errorMessage,
			CreatedAt:          run.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:          run.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}
