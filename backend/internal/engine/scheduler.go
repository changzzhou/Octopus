package engine

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/model"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// Scheduler is the workflow state machine scheduler
type Scheduler struct {
	runsModel  model.WorkflowRunsModel
	stepsModel model.WorkflowRunStepsModel
	redis      *redis.Redis
}

// NewScheduler creates a new scheduler
func NewScheduler(runsModel model.WorkflowRunsModel, stepsModel model.WorkflowRunStepsModel, redis *redis.Redis) *Scheduler {
	return &Scheduler{
		runsModel:  runsModel,
		stepsModel: stepsModel,
		redis:      redis,
	}
}

// DefinitionSnapshot matches types.WorkflowDefinitionSnapshot
type DefinitionSnapshot struct {
	Nodes           []types.Node `json:"nodes"`
	Edges           []types.Edge `json:"edges"`
	EntryNodeId     string       `json:"entry_node_id,omitempty"`
	VariablesSchema string       `json:"variables_schema,omitempty"`
}

// StartRun initiates a workflow run
func (s *Scheduler) StartRun(ctx context.Context, run *model.WorkflowRuns) error {
	var def DefinitionSnapshot
	if err := json.Unmarshal([]byte(run.DefinitionSnapshot), &def); err != nil {
		return fmt.Errorf("failed to parse definition: %w", err)
	}

	if len(def.Nodes) == 0 {
		run.Status = RunStatusSucceeded
		run.FinishedAt = sql.NullTime{Time: time.Now(), Valid: true}
		return s.runsModel.Update(ctx, run)
	}

	for _, node := range def.Nodes {
		step := &model.WorkflowRunSteps{
			RunId:    run.Id,
			NodeId:   node.Id,
			NodeType: node.Type,
			NodeConfig: sql.NullString{
				String: node.Config,
				Valid:  node.Config != "",
			},
			Status: StepStatusPending,
		}
		if _, err := s.stepsModel.Insert(ctx, step); err != nil {
			return fmt.Errorf("failed to create step for node %s: %w", node.Id, err)
		}
	}

	run.Status = RunStatusRunning
	run.StartedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if err := s.runsModel.Update(ctx, run); err != nil {
		return fmt.Errorf("failed to update run status: %w", err)
	}

	return s.advanceRun(ctx, run.Id, run.WorkflowId, &def)
}

// advanceRun advances the run based on DAG topology
func (s *Scheduler) advanceRun(ctx context.Context, runId, workflowId uint64, def *DefinitionSnapshot) error {
	steps, err := s.stepsModel.FindByRunId(ctx, runId)
	if err != nil {
		return fmt.Errorf("failed to get steps: %w", err)
	}

	stepMap := make(map[string]*model.WorkflowRunSteps)
	for _, step := range steps {
		stepMap[step.NodeId] = step
	}

	readyNodes := s.findReadyNodes(def, stepMap)
	s.markUnreachableAsSkipped(ctx, def, stepMap)

	if len(readyNodes) == 0 {
		return s.checkRunCompletion(ctx, runId, steps)
	}

	for _, nodeId := range readyNodes {
		node := s.findNode(def, nodeId)
		if node == nil {
			continue
		}

		step := stepMap[nodeId]
		if err := s.executeNode(ctx, runId, workflowId, step, node, def); err != nil {
			logx.Errorf("failed to execute node %s: %v", nodeId, err)
		}
	}

	return nil
}

// findReadyNodes finds nodes that are pending and have all dependencies satisfied
func (s *Scheduler) findReadyNodes(def *DefinitionSnapshot, stepMap map[string]*model.WorkflowRunSteps) []string {
	var readyNodes []string

	for _, node := range def.Nodes {
		step := stepMap[node.Id]
		if step == nil || step.Status != StepStatusPending {
			continue
		}

		incomingEdges := s.findIncomingEdges(def, node.Id)

		if len(incomingEdges) == 0 {
			readyNodes = append(readyNodes, node.Id)
			continue
		}

		allSatisfied := true
		for _, edge := range incomingEdges {
			sourceStep := stepMap[edge.Source]
			if sourceStep == nil {
				allSatisfied = false
				break
			}

			outlet := edge.Outlet
			if outlet == "" {
				outlet = OutletSuccess
			}

			if outlet == OutletSuccess && sourceStep.Status != StepStatusSucceeded {
				allSatisfied = false
				break
			}
			if outlet == OutletFailure && sourceStep.Status != StepStatusFailed {
				allSatisfied = false
				break
			}
		}

		if allSatisfied {
			readyNodes = append(readyNodes, node.Id)
		}
	}

	return readyNodes
}

// markUnreachableAsSkipped marks nodes that can never be reached as skipped
func (s *Scheduler) markUnreachableAsSkipped(ctx context.Context, def *DefinitionSnapshot, stepMap map[string]*model.WorkflowRunSteps) {
	for _, node := range def.Nodes {
		step := stepMap[node.Id]
		if step == nil || step.Status != StepStatusPending {
			continue
		}

		incomingEdges := s.findIncomingEdges(def, node.Id)
		if len(incomingEdges) == 0 {
			continue
		}

		allUnreachable := true
		for _, edge := range incomingEdges {
			sourceStep := stepMap[edge.Source]
			if sourceStep == nil {
				allUnreachable = false
				break
			}

			outlet := edge.Outlet
			if outlet == "" {
				outlet = OutletSuccess
			}

			sourceTerminal := sourceStep.Status == StepStatusSucceeded ||
				sourceStep.Status == StepStatusFailed ||
				sourceStep.Status == StepStatusSkipped

			if !sourceTerminal {
				allUnreachable = false
				break
			}

			canBeReached := (outlet == OutletSuccess && sourceStep.Status == StepStatusSucceeded) ||
				(outlet == OutletFailure && sourceStep.Status == StepStatusFailed)

			if canBeReached {
				allUnreachable = false
				break
			}
		}

		if allUnreachable {
			step.Status = StepStatusSkipped
			step.FinishedAt = sql.NullTime{Time: time.Now(), Valid: true}
			if err := s.stepsModel.Update(ctx, step); err != nil {
				logx.Errorf("failed to mark step %s as skipped: %v", step.NodeId, err)
			} else {
				logx.Infof("Marked step %s as skipped (unreachable)", step.NodeId)
			}
		}
	}
}

// findIncomingEdges finds edges targeting the given node
func (s *Scheduler) findIncomingEdges(def *DefinitionSnapshot, nodeId string) []types.Edge {
	var edges []types.Edge
	for _, edge := range def.Edges {
		if edge.Target == nodeId {
			edges = append(edges, edge)
		}
	}
	return edges
}

// findNode finds a node by ID
func (s *Scheduler) findNode(def *DefinitionSnapshot, nodeId string) *types.Node {
	for i := range def.Nodes {
		if def.Nodes[i].Id == nodeId {
			return &def.Nodes[i]
		}
	}
	return nil
}

// executeNode executes a single node
func (s *Scheduler) executeNode(ctx context.Context, runId, workflowId uint64, step *model.WorkflowRunSteps, node *types.Node, def *DefinitionSnapshot) error {
	step.Status = StepStatusRunning
	step.StartedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if err := s.stepsModel.Update(ctx, step); err != nil {
		return fmt.Errorf("failed to update step status: %w", err)
	}

	if IsWorkerDispatchType(node.Type) {
		return s.dispatchToWorker(ctx, runId, workflowId, step, node)
	}

	step.Status = StepStatusSucceeded
	step.FinishedAt = sql.NullTime{Time: time.Now(), Valid: true}
	step.OutputData = sql.NullString{String: `{"message":"in-scheduler execution completed"}`, Valid: true}
	if err := s.stepsModel.Update(ctx, step); err != nil {
		return fmt.Errorf("failed to update step completion: %w", err)
	}

	run, err := s.runsModel.FindOne(ctx, runId)
	if err != nil {
		return fmt.Errorf("failed to get run: %w", err)
	}

	var newDef DefinitionSnapshot
	if err := json.Unmarshal([]byte(run.DefinitionSnapshot), &newDef); err != nil {
		return fmt.Errorf("failed to parse definition: %w", err)
	}

	return s.advanceRun(ctx, runId, workflowId, &newDef)
}

// dispatchToWorker dispatches a task to the worker queue
func (s *Scheduler) dispatchToWorker(ctx context.Context, runId, workflowId uint64, step *model.WorkflowRunSteps, node *types.Node) error {
	var payload interface{}
	if step.NodeConfig.Valid && step.NodeConfig.String != "" {
		_ = json.Unmarshal([]byte(step.NodeConfig.String), &payload)
	}

	if node.Type == NodeTypeHuman {
		step.Status = StepStatusWaitingHuman
		if err := s.stepsModel.Update(ctx, step); err != nil {
			return fmt.Errorf("failed to update step to waiting_human: %w", err)
		}
		logx.Infof("Step %d (node %s) is waiting for human input", step.Id, step.NodeId)
		return nil
	}

	task := NewTaskRequest(node.Type, payload, workflowId, runId, step.NodeId)

	taskJSON, err := task.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize task: %w", err)
	}

	_, err = s.redis.Rpush(TaskQueueKey, taskJSON)
	if err != nil {
		return fmt.Errorf("failed to push task to queue: %w", err)
	}

	logx.Infof("Dispatched task %s to worker queue for step %s", task.TaskId, step.NodeId)
	return nil
}

// checkRunCompletion checks if the run is complete
func (s *Scheduler) checkRunCompletion(ctx context.Context, runId uint64, steps []*model.WorkflowRunSteps) error {
	allComplete := true
	anyFailed := false
	anyWaiting := false

	for _, step := range steps {
		switch step.Status {
		case StepStatusPending, StepStatusRunning:
			allComplete = false
		case StepStatusFailed:
			anyFailed = true
		case StepStatusWaitingHuman:
			anyWaiting = true
			allComplete = false
		}
	}

	if !allComplete {
		return nil
	}

	run, err := s.runsModel.FindOne(ctx, runId)
	if err != nil {
		return fmt.Errorf("failed to get run: %w", err)
	}

	if anyFailed {
		run.Status = RunStatusFailed
	} else if anyWaiting {
		return nil
	} else {
		run.Status = RunStatusSucceeded
	}
	run.FinishedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return s.runsModel.Update(ctx, run)
}

// HandleTaskResult processes a task result from the worker
func (s *Scheduler) HandleTaskResult(ctx context.Context, resp *TaskResponse) error {
	runId, err := parseUint64(resp.Context.ExecutionId)
	if err != nil {
		return fmt.Errorf("invalid execution_id: %w", err)
	}

	workflowId, err := parseUint64(resp.Context.WorkflowId)
	if err != nil {
		return fmt.Errorf("invalid workflow_id: %w", err)
	}

	step, err := s.stepsModel.FindOneByRunIdNodeId(ctx, runId, resp.Context.NodeId)
	if err != nil {
		return fmt.Errorf("failed to find step: %w", err)
	}

	step.FinishedAt = sql.NullTime{Time: time.Now(), Valid: true}

	if resp.Status == "success" {
		step.Status = StepStatusSucceeded
		if resp.Result != nil {
			resultJSON, _ := json.Marshal(resp.Result)
			step.OutputData = sql.NullString{String: string(resultJSON), Valid: true}
		}
	} else {
		step.Status = StepStatusFailed
		step.ErrorMessage = sql.NullString{String: resp.Error, Valid: true}
	}

	if err := s.stepsModel.Update(ctx, step); err != nil {
		return fmt.Errorf("failed to update step: %w", err)
	}

	run, err := s.runsModel.FindOne(ctx, runId)
	if err != nil {
		return fmt.Errorf("failed to get run: %w", err)
	}

	var def DefinitionSnapshot
	if err := json.Unmarshal([]byte(run.DefinitionSnapshot), &def); err != nil {
		return fmt.Errorf("failed to parse definition: %w", err)
	}

	return s.advanceRun(ctx, runId, workflowId, &def)
}

func parseUint64(s string) (uint64, error) {
	var result uint64
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
