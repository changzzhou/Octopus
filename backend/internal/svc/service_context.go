package svc

import (
	"backend/internal/config"
	"backend/internal/engine"
	"backend/internal/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config            config.Config
	WorkflowModel     model.WorkflowsModel
	RunModel          model.WorkflowRunsModel
	StepModel         model.WorkflowRunStepsModel
	RedisClient       *redis.Redis
	Scheduler         *engine.Scheduler
	ResultPoller      *engine.ResultPoller
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := c.MustNewMySqlConn()
	redisClient := redis.MustNewRedis(c.Redis)
	runModel := model.NewWorkflowRunsModel(conn)
	stepModel := model.NewWorkflowRunStepsModel(conn)

	scheduler := engine.NewScheduler(runModel, stepModel, redisClient)
	resultPoller := engine.NewResultPoller(runModel, stepModel, redisClient)

	return &ServiceContext{
		Config:        c,
		WorkflowModel: model.NewWorkflowsModel(conn),
		RunModel:      runModel,
		StepModel:     stepModel,
		RedisClient:   redisClient,
		Scheduler:     scheduler,
		ResultPoller:  resultPoller,
	}
}
