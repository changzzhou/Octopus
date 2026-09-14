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

type GetWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkflowLogic {
	return &GetWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWorkflowLogic) GetWorkflow(req *types.GetWorkflowReq) (resp *types.GetWorkflowResp, err error) {
	workflow, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errorx.NewNotFoundError("workflow not found")
		}
		l.Logger.Errorf("get workflow failed: %v", err)
		return nil, errorx.NewInternalError("failed to get workflow")
	}

	var nodes []types.Node
	if workflow.Nodes.Valid && workflow.Nodes.String != "" {
		if err := json.Unmarshal([]byte(workflow.Nodes.String), &nodes); err != nil {
			l.Logger.Errorf("unmarshal nodes failed: %v", err)
			nodes = []types.Node{}
		}
	}
	if nodes == nil {
		nodes = []types.Node{}
	}

	var edges []types.Edge
	if workflow.Edges.Valid && workflow.Edges.String != "" {
		if err := json.Unmarshal([]byte(workflow.Edges.String), &edges); err != nil {
			l.Logger.Errorf("unmarshal edges failed: %v", err)
			edges = []types.Edge{}
		}
	}
	if edges == nil {
		edges = []types.Edge{}
	}

	var canvasMeta *types.CanvasMeta
	if workflow.CanvasMeta.Valid && workflow.CanvasMeta.String != "" {
		canvasMeta = &types.CanvasMeta{}
		if err := json.Unmarshal([]byte(workflow.CanvasMeta.String), canvasMeta); err != nil {
			l.Logger.Errorf("unmarshal canvas_meta failed: %v", err)
			canvasMeta = nil
		}
	}

	return &types.GetWorkflowResp{
		Workflow: types.WorkflowDetail{
			Id:              int64(workflow.Id),
			Name:            workflow.Name,
			Description:     workflow.Description.String,
			Status:          workflow.Status,
			Version:         int(workflow.Version),
			Nodes:           nodes,
			Edges:           edges,
			EntryNodeId:     workflow.EntryNodeId.String,
			VariablesSchema: workflow.VariablesSchema.String,
			CanvasMeta:      canvasMeta,
			CreatedBy:       workflow.CreatedBy.String,
			UpdatedBy:       workflow.UpdatedBy.String,
			CreatedAt:       workflow.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:       workflow.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}
