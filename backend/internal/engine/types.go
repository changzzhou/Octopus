package engine

// Node types that require worker dispatch
const (
	NodeTypeScript = "script"
	NodeTypeHTTP   = "http"
	NodeTypeHuman  = "human"
)

// Step statuses
const (
	StepStatusPending      = "pending"
	StepStatusRunning      = "running"
	StepStatusSucceeded    = "succeeded"
	StepStatusFailed       = "failed"
	StepStatusWaitingHuman = "waiting_human"
	StepStatusSkipped      = "skipped"
)

// Run statuses
const (
	RunStatusPending   = "pending"
	RunStatusRunning   = "running"
	RunStatusSucceeded = "succeeded"
	RunStatusFailed    = "failed"
	RunStatusCancelled = "cancelled"
)

// Edge outlets
const (
	OutletSuccess = "success"
	OutletFailure = "failure"
)

// IsWorkerDispatchType returns true if the node type requires worker dispatch
func IsWorkerDispatchType(nodeType string) bool {
	switch nodeType {
	case NodeTypeScript, NodeTypeHTTP, NodeTypeHuman:
		return true
	default:
		return false
	}
}
