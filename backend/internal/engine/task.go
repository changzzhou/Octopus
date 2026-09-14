package engine

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// TaskContext contains execution context for a task
type TaskContext struct {
	WorkflowId  string `json:"workflow_id"`
	ExecutionId string `json:"execution_id"`
	NodeId      string `json:"node_id"`
	Attempt     int    `json:"attempt"`
}

// TaskRequest is the unified task request sent to workers
type TaskRequest struct {
	TaskId    string      `json:"task_id"`
	TaskType  string      `json:"task_type"`
	Payload   interface{} `json:"payload"`
	Context   TaskContext `json:"context"`
	CreatedAt string      `json:"created_at"`
}

// TaskResponse is the unified response from workers
type TaskResponse struct {
	TaskId      string      `json:"task_id"`
	Status      string      `json:"status"` // "success" | "failure" | "pending"
	Result      interface{} `json:"result,omitempty"`
	Error       string      `json:"error,omitempty"`
	CompletedAt string      `json:"completed_at"`
	Context     TaskContext `json:"context,omitempty"`
}

// Redis queue keys
const (
	TaskQueueKey   = "octopus:tasks"
	ResultQueueKey = "octopus:results"
)

// NewTaskRequest creates a new task request
func NewTaskRequest(taskType string, payload interface{}, workflowId uint64, runId uint64, nodeId string) *TaskRequest {
	return &TaskRequest{
		TaskId:   uuid.New().String(),
		TaskType: taskType,
		Payload:  payload,
		Context: TaskContext{
			WorkflowId:  strconv.FormatUint(workflowId, 10),
			ExecutionId: strconv.FormatUint(runId, 10),
			NodeId:      nodeId,
			Attempt:     1,
		},
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// ToJSON serializes the task request to JSON
func (t *TaskRequest) ToJSON() (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ParseTaskResponse parses a JSON task response
func ParseTaskResponse(data string) (*TaskResponse, error) {
	var resp TaskResponse
	err := json.Unmarshal([]byte(data), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
