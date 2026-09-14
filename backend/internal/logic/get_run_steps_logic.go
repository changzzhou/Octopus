// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRunStepsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get run steps
func NewGetRunStepsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRunStepsLogic {
	return &GetRunStepsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRunStepsLogic) GetRunSteps(req *types.GetRunStepsReq) (resp *types.GetRunStepsResp, err error) {
	steps, err := l.svcCtx.StepModel.FindByRunId(l.ctx, uint64(req.RunId))
	if err != nil {
		return nil, err
	}

	var result []types.Step
	for _, step := range steps {
		nodeConfig := ""
		if step.NodeConfig.Valid {
			nodeConfig = step.NodeConfig.String
		}

		inputData := ""
		if step.InputData.Valid {
			inputData = step.InputData.String
		}

		outputData := ""
		if step.OutputData.Valid {
			outputData = step.OutputData.String
		}

		errorMessage := ""
		if step.ErrorMessage.Valid {
			errorMessage = step.ErrorMessage.String
		}

		startedAt := ""
		if step.StartedAt.Valid {
			startedAt = step.StartedAt.Time.Format("2006-01-02T15:04:05Z")
		}

		finishedAt := ""
		if step.FinishedAt.Valid {
			finishedAt = step.FinishedAt.Time.Format("2006-01-02T15:04:05Z")
		}

		result = append(result, types.Step{
			Id:           int64(step.Id),
			RunId:        int64(step.RunId),
			NodeId:       step.NodeId,
			NodeType:     step.NodeType,
			NodeConfig:   nodeConfig,
			Status:       step.Status,
			InputData:    inputData,
			OutputData:   outputData,
			ErrorMessage: errorMessage,
			StartedAt:    startedAt,
			FinishedAt:   finishedAt,
			CreatedAt:    step.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:    step.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &types.GetRunStepsResp{
		Steps: result,
	}, nil
}
