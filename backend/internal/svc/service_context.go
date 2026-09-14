package svc

import (
	"backend/internal/config"
	"backend/internal/model"
)

type ServiceContext struct {
	Config        config.Config
	WorkflowModel model.WorkflowsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := c.MustNewMySqlConn()
	return &ServiceContext{
		Config:        c,
		WorkflowModel: model.NewWorkflowsModel(conn),
	}
}
