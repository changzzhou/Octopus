# BE-3: SSE Realtime Events

This document describes the Server-Sent Events (SSE) endpoint for subscribing to workflow run events in real-time.

## API Endpoint

```
GET /api/v1/runs/:runId/events
```

Returns an SSE stream of events for the specified run.

## Event Types

### 1. `run.status_changed`

Emitted when a run's status changes (e.g., pending → running, running → succeeded).

```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440000",
  "event_type": "run.status_changed",
  "run_id": 123,
  "workflow_id": 456,
  "occurred_at": "2026-09-14T15:30:00.123456789Z",
  "sequence": 1,
  "payload": {
    "from_status": "pending",
    "to_status": "running",
    "reason": ""
  }
}
```

### 2. `step.status_changed`

Emitted when a step's status changes.

```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440001",
  "event_type": "step.status_changed",
  "run_id": 123,
  "workflow_id": 456,
  "occurred_at": "2026-09-14T15:30:01.123456789Z",
  "sequence": 2,
  "payload": {
    "step_id": "node_abc123",
    "from": "pending",
    "to": "running",
    "error_summary": ""
  }
}
```

### 3. `run.terminal`

Emitted when a run reaches a terminal state (succeeded, failed, cancelled).

```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440002",
  "event_type": "run.terminal",
  "run_id": 123,
  "workflow_id": 456,
  "occurred_at": "2026-09-14T15:30:05.123456789Z",
  "sequence": 5,
  "payload": {
    "final_status": "succeeded"
  }
}
```

### 4. `human.waiting`

Emitted when a step enters `waiting_human` status, indicating human input is required.

```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440003",
  "event_type": "human.waiting",
  "run_id": 123,
  "workflow_id": 456,
  "occurred_at": "2026-09-14T15:30:02.123456789Z",
  "sequence": 3,
  "payload": {
    "step_id": "human_node_xyz",
    "prompt_summary": "",
    "context_refs": null
  }
}
```

## Common Envelope

All events share a common envelope structure:

| Field        | Type    | Description                                        |
|--------------|---------|---------------------------------------------------|
| `event_id`   | string  | Unique UUID for this event                         |
| `event_type` | string  | Type of event (see above)                          |
| `run_id`     | int64   | ID of the workflow run                             |
| `workflow_id`| int64   | ID of the workflow definition                      |
| `occurred_at`| string  | ISO 8601 timestamp with nanoseconds                |
| `sequence`   | int64   | Monotonically increasing sequence number per run   |
| `payload`    | object  | Event-specific payload                             |

## Reconnection Strategy

SSE events are **not** the sole source of truth. When reconnecting:

1. First, fetch the current run state: `GET /api/v1/runs/:runId`
2. Then, fetch all steps: `GET /api/v1/runs/:runId/steps`
3. Finally, establish SSE connection to receive future updates

This ensures you have the complete current state before receiving incremental updates.

## Acceptance Test Steps

### Prerequisites

1. Backend server running on port 8888
2. MySQL and Redis running
3. A workflow created and enabled

### Test Procedure

**Terminal 1 - Subscribe to SSE events:**
```bash
# First, create a workflow if needed
curl -X POST http://localhost:8888/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{"name":"Test SSE Workflow","nodes":[{"id":"n1","type":"script","name":"Script Node","position":{"x":0,"y":0},"config":"{}"}],"edges":[]}'

# Note the workflow ID from response (e.g., 1)

# Trigger a run to get a run ID
curl -X POST http://localhost:8888/api/v1/workflows/1/runs

# Note the run_id from response (e.g., 1)

# Subscribe to SSE events for the run
curl -N http://localhost:8888/api/v1/runs/1/events
```

**Terminal 2 - Trigger another run and observe events:**
```bash
# Trigger a new run
curl -X POST http://localhost:8888/api/v1/workflows/1/runs

# Use the returned run_id in Terminal 1's curl command
```

**Expected SSE Output:**
```
event: run.status_changed
id: 550e8400-e29b-41d4-a716-446655440000
data: {"event_id":"550e8400-...","event_type":"run.status_changed",...}

event: step.status_changed
id: 550e8400-e29b-41d4-a716-446655440001
data: {"event_id":"550e8400-...","event_type":"step.status_changed",...}

event: run.terminal
id: 550e8400-e29b-41d4-a716-446655440002
data: {"event_id":"550e8400-...","event_type":"run.terminal",...}
```

### Verify Keepalive

If no events occur within 15 seconds, you should see:
```
: keepalive

```

## Implementation Notes

- Events are published in-memory via a hub/subscriber pattern
- Events are not persisted; SSE is for real-time notifications only
- Clients should always fetch state via REST APIs for the source of truth
- The `sequence` field can help detect missed events during reconnection
