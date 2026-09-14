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

type EnableWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnableWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnableWorkflowLogic {
	return &EnableWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnableWorkflowLogic) EnableWorkflow(req *types.EnableWorkflowReq) (resp *types.EnableWorkflowResp, err error) {
	workflow, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errorx.NewNotFoundError("workflow not found")
		}
		l.Logger.Errorf("find workflow failed: %v", err)
		return nil, errorx.NewInternalError("failed to find workflow")
	}

	if workflow.Status != "draft" {
		return nil, errorx.NewBadRequestError("can only enable draft workflows; current status: " + workflow.Status)
	}

	var nodes []types.Node
	if workflow.Nodes.Valid && workflow.Nodes.String != "" {
		if err := json.Unmarshal([]byte(workflow.Nodes.String), &nodes); err != nil {
			return nil, errorx.NewBadRequestError("invalid nodes format")
		}
	}

	var edges []types.Edge
	if workflow.Edges.Valid && workflow.Edges.String != "" {
		if err := json.Unmarshal([]byte(workflow.Edges.String), &edges); err != nil {
			return nil, errorx.NewBadRequestError("invalid edges format")
		}
	}

	validationErrors := validateGraph(nodes, edges, workflow.EntryNodeId.String)
	if len(validationErrors) > 0 {
		return nil, errorx.NewBadRequestError("workflow validation failed: " + validationErrors[0].Message)
	}

	err = l.svcCtx.WorkflowModel.UpdateStatus(l.ctx, uint64(req.Id), "enabled")
	if err != nil {
		l.Logger.Errorf("enable workflow failed: %v", err)
		return nil, errorx.NewInternalError("failed to enable workflow")
	}

	return &types.EnableWorkflowResp{
		Status:  "enabled",
		Version: int(workflow.Version),
	}, nil
}
