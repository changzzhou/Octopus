package engine

import (
	"context"
	"time"

	"backend/internal/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// ResultPoller polls for task results from Redis
type ResultPoller struct {
	scheduler *Scheduler
	redis     *redis.Redis
	stopCh    chan struct{}
}

// NewResultPoller creates a new result poller
func NewResultPoller(runsModel model.WorkflowRunsModel, stepsModel model.WorkflowRunStepsModel, redis *redis.Redis) *ResultPoller {
	return &ResultPoller{
		scheduler: NewScheduler(runsModel, stepsModel, redis),
		redis:     redis,
		stopCh:    make(chan struct{}),
	}
}

// Start starts the result poller in a goroutine
func (p *ResultPoller) Start(ctx context.Context) {
	go func() {
		logx.Info("Result poller started")
		for {
			select {
			case <-p.stopCh:
				logx.Info("Result poller stopped")
				return
			case <-ctx.Done():
				logx.Info("Result poller context cancelled")
				return
			default:
				p.pollOnce(ctx)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the result poller
func (p *ResultPoller) Stop() {
	close(p.stopCh)
}

// pollOnce polls for one result
func (p *ResultPoller) pollOnce(ctx context.Context) {
	result, err := p.redis.Lpop(ResultQueueKey)
	if err != nil {
		if err.Error() != "redis: nil" {
			logx.Errorf("Failed to pop result from queue: %v", err)
		}
		return
	}

	if result == "" {
		return
	}

	resp, err := ParseTaskResponse(result)
	if err != nil {
		logx.Errorf("Failed to parse task result: %v", err)
		return
	}

	logx.Infof("Received task result: task_id=%s, status=%s", resp.TaskId, resp.Status)

	if err := p.scheduler.HandleTaskResult(ctx, resp); err != nil {
		logx.Errorf("Failed to handle task result: %v", err)
	}
}

// GetScheduler returns the underlying scheduler
func (p *ResultPoller) GetScheduler() *Scheduler {
	return p.scheduler
}
