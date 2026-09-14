package sse

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
)

type EventEnvelope struct {
	EventId    string      `json:"event_id"`
	EventType  string      `json:"event_type"`
	RunId      int64       `json:"run_id"`
	WorkflowId int64       `json:"workflow_id"`
	OccurredAt string      `json:"occurred_at"`
	Sequence   int64       `json:"sequence,omitempty"`
	Payload    interface{} `json:"payload"`
}

type RunStatusChangedPayload struct {
	FromStatus string `json:"from_status"`
	ToStatus   string `json:"to_status"`
	Reason     string `json:"reason,omitempty"`
}

type StepStatusChangedPayload struct {
	StepId       string `json:"step_id"`
	From         string `json:"from"`
	To           string `json:"to"`
	ErrorSummary string `json:"error_summary,omitempty"`
}

type RunTerminalPayload struct {
	FinalStatus string `json:"final_status"`
}

type HumanWaitingPayload struct {
	StepId        string   `json:"step_id"`
	PromptSummary string   `json:"prompt_summary,omitempty"`
	ContextRefs   []string `json:"context_refs,omitempty"`
}

const (
	EventRunStatusChanged  = "run.status_changed"
	EventStepStatusChanged = "step.status_changed"
	EventRunTerminal       = "run.terminal"
	EventHumanWaiting      = "human.waiting"
)

type Subscriber struct {
	RunId   int64
	EventCh chan *EventEnvelope
	Done    chan struct{}
}

type Hub struct {
	mu          sync.RWMutex
	subscribers map[int64]map[*Subscriber]struct{}
	sequences   map[int64]int64
}

var globalHub *Hub
var once sync.Once

func GetHub() *Hub {
	once.Do(func() {
		globalHub = &Hub{
			subscribers: make(map[int64]map[*Subscriber]struct{}),
			sequences:   make(map[int64]int64),
		}
	})
	return globalHub
}

func (h *Hub) Subscribe(runId int64) *Subscriber {
	h.mu.Lock()
	defer h.mu.Unlock()

	sub := &Subscriber{
		RunId:   runId,
		EventCh: make(chan *EventEnvelope, 100),
		Done:    make(chan struct{}),
	}

	if h.subscribers[runId] == nil {
		h.subscribers[runId] = make(map[*Subscriber]struct{})
	}
	h.subscribers[runId][sub] = struct{}{}

	return sub
}

func (h *Hub) Unsubscribe(sub *Subscriber) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if subs, ok := h.subscribers[sub.RunId]; ok {
		delete(subs, sub)
		if len(subs) == 0 {
			delete(h.subscribers, sub.RunId)
		}
	}
	close(sub.Done)
}

func (h *Hub) Publish(runId, workflowId int64, eventType string, payload interface{}) {
	h.mu.Lock()
	h.sequences[runId]++
	seq := h.sequences[runId]
	subs := h.subscribers[runId]
	h.mu.Unlock()

	if len(subs) == 0 {
		return
	}

	envelope := &EventEnvelope{
		EventId:    uuid.New().String(),
		EventType:  eventType,
		RunId:      runId,
		WorkflowId: workflowId,
		OccurredAt: time.Now().UTC().Format(time.RFC3339Nano),
		Sequence:   seq,
		Payload:    payload,
	}

	for sub := range subs {
		select {
		case sub.EventCh <- envelope:
		default:
		}
	}
}

func (e *EventEnvelope) ToSSEData() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
