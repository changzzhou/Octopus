package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkflowRunsModel = (*customWorkflowRunsModel)(nil)

type (
	// WorkflowRunsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkflowRunsModel.
	WorkflowRunsModel interface {
		workflowRunsModel
		withSession(session sqlx.Session) WorkflowRunsModel
		FindByWorkflowId(ctx context.Context, workflowId uint64) ([]*WorkflowRuns, error)
	}

	customWorkflowRunsModel struct {
		*defaultWorkflowRunsModel
	}
)

// NewWorkflowRunsModel returns a model for the database table.
func NewWorkflowRunsModel(conn sqlx.SqlConn) WorkflowRunsModel {
	return &customWorkflowRunsModel{
		defaultWorkflowRunsModel: newWorkflowRunsModel(conn),
	}
}

func (m *customWorkflowRunsModel) withSession(session sqlx.Session) WorkflowRunsModel {
	return NewWorkflowRunsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customWorkflowRunsModel) FindByWorkflowId(ctx context.Context, workflowId uint64) ([]*WorkflowRuns, error) {
	query := fmt.Sprintf("select %s from %s where `workflow_id` = ? order by `id` desc", workflowRunsRows, m.table)
	var resp []*WorkflowRuns
	err := m.conn.QueryRowsCtx(ctx, &resp, query, workflowId)
	return resp, err
}
