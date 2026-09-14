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

type CreateWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkflowLogic {
	return &CreateWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateWorkflowLogic) CreateWorkflow(req *types.CreateWorkflowReq) (resp *types.CreateWorkflowResp, err error) {
	nodesJSON, err := json.Marshal(req.Nodes)
	if err != nil {
		return nil, errorx.NewBadRequestError("invalid nodes format")
	}

	edgesJSON, err := json.Marshal(req.Edges)
	if err != nil {
		return nil, errorx.NewBadRequestError("invalid edges format")
	}

	var canvasMetaJSON []byte
	if req.CanvasMeta != nil {
		canvasMetaJSON, err = json.Marshal(req.CanvasMeta)
		if err != nil {
			return nil, errorx.NewBadRequestError("invalid canvas_meta format")
		}
	}

	workflow := &model.Workflows{
		Name:            req.Name,
		Description:     sql.NullString{String: req.Description, Valid: req.Description != ""},
		Status:          "draft",
		Version:         1,
		Nodes:           sql.NullString{String: string(nodesJSON), Valid: true},
		Edges:           sql.NullString{String: string(edgesJSON), Valid: true},
		EntryNodeId:     sql.NullString{String: req.EntryNodeId, Valid: req.EntryNodeId != ""},
		VariablesSchema: sql.NullString{String: req.VariablesSchema, Valid: req.VariablesSchema != ""},
		CanvasMeta:      sql.NullString{String: string(canvasMetaJSON), Valid: len(canvasMetaJSON) > 0},
	}

	result, err := l.svcCtx.WorkflowModel.Insert(l.ctx, workflow)
	if err != nil {
		l.Logger.Errorf("create workflow failed: %v", err)
		return nil, errorx.NewInternalError("failed to create workflow")
	}

	id, err := result.LastInsertId()
	if err != nil {
		l.Logger.Errorf("get last insert id failed: %v", err)
		return nil, errorx.NewInternalError("failed to get workflow id")
	}

	return &types.CreateWorkflowResp{
		Id:      id,
		Version: 1,
		Status:  "draft",
	}, nil
}
