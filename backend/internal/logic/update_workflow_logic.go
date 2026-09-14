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

type UpdateWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkflowLogic {
	return &UpdateWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWorkflowLogic) UpdateWorkflow(req *types.UpdateWorkflowReq) (resp *types.UpdateWorkflowResp, err error) {
	existing, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errorx.NewNotFoundError("workflow not found")
		}
		l.Logger.Errorf("find workflow failed: %v", err)
		return nil, errorx.NewInternalError("failed to find workflow")
	}

	if req.Version != int(existing.Version) {
		return nil, errorx.ErrVersionConflict
	}

	name := existing.Name
	if req.Name != "" {
		name = req.Name
	}

	description := existing.Description
	if req.Description != "" {
		description = sql.NullString{String: req.Description, Valid: true}
	}

	var nodesJSON string
	if req.Nodes != nil {
		b, err := json.Marshal(req.Nodes)
		if err != nil {
			return nil, errorx.NewBadRequestError("invalid nodes format")
		}
		nodesJSON = string(b)
	} else {
		nodesJSON = existing.Nodes.String
	}

	var edgesJSON string
	if req.Edges != nil {
		b, err := json.Marshal(req.Edges)
		if err != nil {
			return nil, errorx.NewBadRequestError("invalid edges format")
		}
		edgesJSON = string(b)
	} else {
		edgesJSON = existing.Edges.String
	}

	entryNodeId := existing.EntryNodeId
	if req.EntryNodeId != "" {
		entryNodeId = sql.NullString{String: req.EntryNodeId, Valid: true}
	}

	variablesSchema := existing.VariablesSchema
	if req.VariablesSchema != "" {
		variablesSchema = sql.NullString{String: req.VariablesSchema, Valid: true}
	}

	canvasMeta := existing.CanvasMeta
	if req.CanvasMeta != nil {
		b, err := json.Marshal(req.CanvasMeta)
		if err != nil {
			return nil, errorx.NewBadRequestError("invalid canvas_meta format")
		}
		canvasMeta = sql.NullString{String: string(b), Valid: true}
	}

	newVersion := existing.Version + 1

	updated := &model.Workflows{
		Id:              existing.Id,
		Name:            name,
		Description:     description,
		Nodes:           sql.NullString{String: nodesJSON, Valid: true},
		Edges:           sql.NullString{String: edgesJSON, Valid: true},
		EntryNodeId:     entryNodeId,
		VariablesSchema: variablesSchema,
		CanvasMeta:      canvasMeta,
		Version:         newVersion,
	}

	result, err := l.svcCtx.WorkflowModel.UpdateWithVersion(l.ctx, updated, existing.Version)
	if err != nil {
		l.Logger.Errorf("update workflow failed: %v", err)
		return nil, errorx.NewInternalError("failed to update workflow")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		l.Logger.Errorf("get rows affected failed: %v", err)
		return nil, errorx.NewInternalError("failed to check update result")
	}

	if rowsAffected == 0 {
		return nil, errorx.ErrVersionConflict
	}

	return &types.UpdateWorkflowResp{
		Version: int(newVersion),
	}, nil
}
