package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkflowRunStepsModel = (*customWorkflowRunStepsModel)(nil)

type (
	// WorkflowRunStepsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkflowRunStepsModel.
	WorkflowRunStepsModel interface {
		workflowRunStepsModel
		withSession(session sqlx.Session) WorkflowRunStepsModel
		FindByRunId(ctx context.Context, runId uint64) ([]*WorkflowRunSteps, error)
	}

	customWorkflowRunStepsModel struct {
		*defaultWorkflowRunStepsModel
	}
)

// NewWorkflowRunStepsModel returns a model for the database table.
func NewWorkflowRunStepsModel(conn sqlx.SqlConn) WorkflowRunStepsModel {
	return &customWorkflowRunStepsModel{
		defaultWorkflowRunStepsModel: newWorkflowRunStepsModel(conn),
	}
}

func (m *customWorkflowRunStepsModel) withSession(session sqlx.Session) WorkflowRunStepsModel {
	return NewWorkflowRunStepsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customWorkflowRunStepsModel) FindByRunId(ctx context.Context, runId uint64) ([]*WorkflowRunSteps, error) {
	query := fmt.Sprintf("select %s from %s where `run_id` = ? order by `id` asc", workflowRunStepsRows, m.table)
	var resp []*WorkflowRunSteps
	err := m.conn.QueryRowsCtx(ctx, &resp, query, runId)
	return resp, err
}
