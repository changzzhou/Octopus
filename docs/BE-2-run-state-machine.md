# BE-2: Run/State Machine + Worker Minimal Closed Loop

## Overview

This sprint implements a minimal runnable loop for workflow execution as **incremental** changes on top of BE-1:
1. Create a simple workflow with nodes and edges (using BE-1 CRUD)
2. POST to trigger a run
3. Go state machine advances by topology using Edge.outlet
4. Redis queue dispatches tasks to Python worker
5. Worker executes stub handlers and reports results
6. Steps reach succeeded/failed/skipped state and are queryable

## Incremental Changes (on top of BE-1)

This PR preserves all BE-1 functionality:
- Workflow CRUD (create/get/list/update/delete)
- Validate/Enable/Disable APIs
- String status enum (`draft|enabled|disabled`)
- goctl Spec-First layout (`internal/model/`, `internal/logic/`, etc.)

**Added in BE-2:**
- `workflow_runs` and `workflow_run_steps` tables
- Run/Step models via `goctl model mysql ddl` (without `-c`)
- Run APIs (trigger, get run, get steps)
- State machine engine (`internal/engine/`)
- Redis configuration and result poller

## Components

### 1. Database Schema (Added Tables)

```sql
-- workflow_runs: Run instances with frozen definition_snapshot
-- workflow_run_steps: Individual step execution records
```

#### Key Fields
- `workflow_runs.definition_snapshot` - Frozen workflow definition at run time
- `workflow_runs.workflow_version` - Version of workflow when run was triggered
- `workflow_run_steps.status` - pending/running/succeeded/failed/waiting_human/skipped

### 2. New API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/workflows/:id/runs` | Trigger a workflow run (draft allowed) |
| GET | `/api/v1/runs/:runId` | Get run details |
| GET | `/api/v1/runs/:runId/steps` | Get run steps |

### 3. State Machine

The Go scheduler (`internal/engine/scheduler.go`) implements DAG-based execution:

- **Topology-based advancement**: Nodes are executed when all their incoming edge conditions are satisfied
- **Edge.outlet routing**: Supports `success` and `failure` outlets
  - `success` outlet: Target node runs when source succeeds
  - `failure` outlet: Target node runs when source fails
- **Unreachable node skipping**: Nodes that cannot be reached due to outlet mismatch are marked as `skipped`
- **Worker dispatch**: `script`, `http`, and `human` node types are dispatched to workers
- **In-scheduler execution**: Other node types execute synchronously in the scheduler

### 4. Redis Queue Protocol

Tasks are dispatched via Redis lists:
- **Task Queue**: `octopus:tasks` - Tasks awaiting worker pickup
- **Result Queue**: `octopus:results` - Worker results for scheduler processing

#### Task Context
- `workflow_id`: The actual workflow ID (not run ID)
- `execution_id`: The run ID
- `node_id`: The node being executed

#### Task Request Format
```json
{
  "task_id": "uuid",
  "task_type": "script|http",
  "payload": {},
  "context": {
    "workflow_id": "1",
    "execution_id": "1",
    "node_id": "node_1",
    "attempt": 1
  },
  "created_at": "2024-01-01T00:00:00Z"
}
```

### 5. Python Worker

Located in `worker/`, the Python worker:
- Polls Redis task queue for tasks
- Dispatches to registered handlers by task_type
- Reports results back to result queue with context

Supported handlers:
- `script` - Script execution stub
- `http` - HTTP request stub
- `noop` - No-op for testing
- `echo` - Echo payload back

### 6. Human Node Support

Human nodes enter `waiting_human` status and pause execution:
- No task is dispatched to workers
- Run continues for other parallel branches
- Full resume UI will be implemented in future sprints

## Verification

### Prerequisites
```bash
# Start MySQL and Redis
docker compose up -d

# Apply schema (includes new run/step tables)
mysql -h localhost -u root -poctopus_root_pwd octopus < docker/mysql/init/00_schema.sql

# Start backend
cd backend && go run workflow.go

# Start worker (in another terminal)
cd worker && pip install redis python-dotenv && python main.py
```

### Test Flow

#### 1. Create a workflow with 2 nodes (using BE-1 API)
```bash
curl -X POST http://localhost:8888/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Workflow",
    "description": "A simple 2-node workflow",
    "nodes": [
      {"id": "node_1", "type": "script", "name": "Script Node", "position": {"x": 100, "y": 100}},
      {"id": "node_2", "type": "http", "name": "HTTP Node", "position": {"x": 300, "y": 100}}
    ],
    "edges": [
      {"id": "edge_1", "source": "node_1", "target": "node_2", "outlet": "success"}
    ]
  }'
```

#### 2. Trigger a run
```bash
curl -X POST http://localhost:8888/api/v1/workflows/1/runs
# Response: {"run_id": 1}
```

#### 3. Query run status
```bash
curl http://localhost:8888/api/v1/runs/1
```

#### 4. Query run steps
```bash
curl http://localhost:8888/api/v1/runs/1/steps
```

#### Expected Step Statuses
With worker running:
- Both steps should reach `succeeded` status
- Run status should be `succeeded`

Without worker:
- First step stays in `running` (task dispatched but not processed)
- Second step stays in `pending`

## Configuration

### Backend (`backend/etc/workflow-api.yaml`)
```yaml
Name: workflow-api
Host: 0.0.0.0
Port: 8888

MySQL:
  DataSource: root:octopus_root_pwd@tcp(localhost:3306)/octopus?charset=utf8mb4&parseTime=True&loc=Local

Redis:
  Host: localhost:6379
  Type: node
```

### Worker Environment
```bash
REDIS_URL=redis://localhost:6379
TASK_QUEUE=octopus:tasks
RESULT_QUEUE=octopus:results
```

## Out of Scope (Future Sprints)

- SSE for real-time updates
- Resume/Cancel/Continue-from-failure
- Credential full APIs
- Full human UI linkage
- Plugin market
