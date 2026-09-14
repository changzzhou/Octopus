package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkflowsModel = (*customWorkflowsModel)(nil)

type (
	WorkflowsModel interface {
		workflowsModel
		withSession(session sqlx.Session) WorkflowsModel
		FindByPage(ctx context.Context, page, pageSize int, status string) ([]*Workflows, error)
		Count(ctx context.Context, status string) (int64, error)
		UpdateWithVersion(ctx context.Context, data *Workflows, oldVersion uint64) (sql.Result, error)
		UpdateStatus(ctx context.Context, id uint64, status string) error
	}

	customWorkflowsModel struct {
		*defaultWorkflowsModel
	}
)

func NewWorkflowsModel(conn sqlx.SqlConn) WorkflowsModel {
	return &customWorkflowsModel{
		defaultWorkflowsModel: newWorkflowsModel(conn),
	}
}

func (m *customWorkflowsModel) withSession(session sqlx.Session) WorkflowsModel {
	return NewWorkflowsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customWorkflowsModel) FindByPage(ctx context.Context, page, pageSize int, status string) ([]*Workflows, error) {
	offset := (page - 1) * pageSize
	var workflows []*Workflows

	if status != "" {
		query := fmt.Sprintf("select %s from %s where `status` = ? order by `updated_at` desc limit ? offset ?", workflowsRows, m.table)
		err := m.conn.QueryRowsCtx(ctx, &workflows, query, status, pageSize, offset)
		return workflows, err
	}

	query := fmt.Sprintf("select %s from %s order by `updated_at` desc limit ? offset ?", workflowsRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &workflows, query, pageSize, offset)
	return workflows, err
}

func (m *customWorkflowsModel) Count(ctx context.Context, status string) (int64, error) {
	var count int64

	if status != "" {
		query := fmt.Sprintf("select count(*) from %s where `status` = ?", m.table)
		err := m.conn.QueryRowCtx(ctx, &count, query, status)
		return count, err
	}

	query := fmt.Sprintf("select count(*) from %s", m.table)
	err := m.conn.QueryRowCtx(ctx, &count, query)
	return count, err
}

func (m *customWorkflowsModel) UpdateWithVersion(ctx context.Context, data *Workflows, oldVersion uint64) (sql.Result, error) {
	fields := []string{
		"`name` = ?",
		"`description` = ?",
		"`nodes` = ?",
		"`edges` = ?",
		"`entry_node_id` = ?",
		"`variables_schema` = ?",
		"`canvas_meta` = ?",
		"`updated_by` = ?",
		"`version` = ?",
	}
	query := fmt.Sprintf("update %s set %s where `id` = ? and `version` = ?", m.table, strings.Join(fields, ", "))
	return m.conn.ExecCtx(ctx, query, data.Name, data.Description, data.Nodes, data.Edges, data.EntryNodeId, data.VariablesSchema, data.CanvasMeta, data.UpdatedBy, data.Version, data.Id, oldVersion)
}

func (m *customWorkflowsModel) UpdateStatus(ctx context.Context, id uint64, status string) error {
	query := fmt.Sprintf("update %s set `status` = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, status, id)
	return err
}
