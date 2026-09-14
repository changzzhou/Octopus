// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"database/sql"
	"encoding/json"

	"backend/internal/errorx"
	"backend/internal/model"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TriggerRunLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Trigger a workflow run (draft allowed)
func NewTriggerRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerRunLogic {
	return &TriggerRunLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TriggerRunLogic) TriggerRun(req *types.TriggerRunReq) (resp *types.TriggerRunResp, err error) {
	workflow, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errorx.NewCodeError(404, "workflow not found")
		}
		return nil, err
	}

	snapshot := types.WorkflowDefinitionSnapshot{
		Nodes:           []types.Node{},
		Edges:           []types.Edge{},
		EntryNodeId:     "",
		VariablesSchema: "",
	}

	if workflow.Nodes.Valid && workflow.Nodes.String != "" {
		var nodes []types.Node
		if err := json.Unmarshal([]byte(workflow.Nodes.String), &nodes); err == nil {
			snapshot.Nodes = nodes
		}
	}

	if workflow.Edges.Valid && workflow.Edges.String != "" {
		var edges []types.Edge
		if err := json.Unmarshal([]byte(workflow.Edges.String), &edges); err == nil {
			snapshot.Edges = edges
		}
	}

	if workflow.EntryNodeId.Valid {
		snapshot.EntryNodeId = workflow.EntryNodeId.String
	}

	if workflow.VariablesSchema.Valid {
		snapshot.VariablesSchema = workflow.VariablesSchema.String
	}

	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		return nil, err
	}

	run := &model.WorkflowRuns{
		WorkflowId:         uint64(req.Id),
		WorkflowVersion:    workflow.Version,
		DefinitionSnapshot: string(snapshotJSON),
		Status:             "pending",
		TriggerType:        "manual",
		StartedAt:          sql.NullTime{},
		FinishedAt:         sql.NullTime{},
		ErrorMessage:       sql.NullString{},
	}

	result, err := l.svcCtx.RunModel.Insert(l.ctx, run)
	if err != nil {
		return nil, err
	}

	runId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	run.Id = uint64(runId)

	go func() {
		ctx := context.Background()
		if err := l.svcCtx.Scheduler.StartRun(ctx, run); err != nil {
			logx.Errorf("Failed to start run %d: %v", runId, err)
		}
	}()

	return &types.TriggerRunResp{
		RunId: runId,
	}, nil
}
