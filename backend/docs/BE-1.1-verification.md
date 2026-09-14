# BE-1.1 Verification Steps

This document provides verification steps for the two fixes introduced in BE-1.1:
1. **disabled → enabled** status transition support
2. **UNREACHABLE_NODE** false negative fix for nodes with incoming edges

## Prerequisites

Start the backend service:
```bash
cd backend
go run workflow.go -f etc/workflow-api.yaml
```

Ensure MySQL is running with the workflow schema from `backend/sql/workflows.sql`.

---

## Fix 1: disabled → enabled Transition

### Test Case 1: Create → Enable → Disable → Re-enable (Success)

```bash
# Step 1: Create a workflow (status = draft)
curl -X POST http://localhost:8888/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Enable Disable Loop",
    "nodes": [
      {"id": "n1", "type": "script", "name": "Entry Node", "position": {"x": 0, "y": 0}}
    ],
    "edges": [],
    "entry_node_id": "n1"
  }'
# Response: {"id": 1, "version": 1, "status": "draft"}

# Step 2: Enable the workflow (draft → enabled)
curl -X POST http://localhost:8888/api/v1/workflows/1/enable
# Response: {"status": "enabled", "version": 1}

# Step 3: Disable the workflow (enabled → disabled)
curl -X POST http://localhost:8888/api/v1/workflows/1/disable
# Response: {"status": "disabled"}

# Step 4: Re-enable the workflow (disabled → enabled) - THIS IS THE FIX
curl -X POST http://localhost:8888/api/v1/workflows/1/enable
# Expected Response: {"status": "enabled", "version": 1}
# Before fix: {"code": 400, "message": "can only enable draft workflows; current status: disabled"}
```

### Test Case 2: Illegal Transitions (Failure Expected)

```bash
# Try to enable an already enabled workflow
curl -X POST http://localhost:8888/api/v1/workflows/1/enable
# Expected: {"code": 400, "message": "can only enable draft or disabled workflows; current status: enabled"}

# Try to disable a draft workflow
curl -X POST http://localhost:8888/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{"name": "Draft Test", "nodes": [{"id": "n1", "type": "script", "name": "N1", "position": {"x": 0, "y": 0}}], "entry_node_id": "n1"}'
# Note the returned ID, e.g., 2

curl -X POST http://localhost:8888/api/v1/workflows/2/disable
# Expected: {"code": 400, "message": "can only disable enabled workflows; current status: draft"}

# Try to delete a non-draft (enabled or disabled) workflow
curl -X DELETE http://localhost:8888/api/v1/workflows/1
# Expected: {"code": 400, "message": "can only delete draft workflows"}
```

### Valid Status Transitions Summary

| Current Status | Enable | Disable | Delete |
|----------------|--------|---------|--------|
| draft          | ✓ (→ enabled) | ✗ | ✓ |
| enabled        | ✗ | ✓ (→ disabled) | ✗ |
| disabled       | ✓ (→ enabled) | ✗ | ✗ |

---

## Fix 2: UNREACHABLE_NODE False Negative

### Bug Description

Previously, the validator only flagged nodes as `UNREACHABLE_NODE` if they had **no incoming edges**. However, a node can have incoming edges but still be unreachable from the `entry_node_id` if the edges originate from other unreachable nodes.

### Example Graph Reproducing the Bug

```
Graph Structure:
  entry_node_id: "A"

  A ──success──> B       (reachable from A)
  
  C ──success──> D       (C and D are isolated subgraph, both have edges but unreachable from A)
```

In this graph:
- Node A is the entry point (reachable)
- Node B is reachable from A via edge
- Node C has no incoming edges (was correctly flagged before)
- Node D has an incoming edge from C, but C is unreachable, so D is also unreachable
  - **Before fix**: D was NOT flagged (false negative) because it has incoming edges
  - **After fix**: D IS correctly flagged as UNREACHABLE_NODE

### Test Case: UNREACHABLE_NODE with Incoming Edges

```bash
# Create a workflow with an unreachable subgraph (nodes have edges but unreachable from entry)
curl -X POST http://localhost:8888/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Unreachable Subgraph Test",
    "nodes": [
      {"id": "A", "type": "script", "name": "Entry A", "position": {"x": 0, "y": 0}},
      {"id": "B", "type": "script", "name": "Node B", "position": {"x": 100, "y": 0}},
      {"id": "C", "type": "script", "name": "Isolated Start", "position": {"x": 0, "y": 100}},
      {"id": "D", "type": "script", "name": "Has Incoming Edge", "position": {"x": 100, "y": 100}}
    ],
    "edges": [
      {"id": "e1", "source": "A", "target": "B", "outlet": "success"},
      {"id": "e2", "source": "C", "target": "D", "outlet": "success"}
    ],
    "entry_node_id": "A"
  }'
# Response: {"id": 3, "version": 1, "status": "draft"}

# Validate the workflow
curl -X POST http://localhost:8888/api/v1/workflows/3/validate
# Expected Response (AFTER FIX - both C and D are flagged):
# {
#   "valid": false,
#   "errors": [
#     {"code": "UNREACHABLE_NODE", "message": "node is not reachable from entry node: C", "node_id": "C"},
#     {"code": "UNREACHABLE_NODE", "message": "node is not reachable from entry node: D", "node_id": "D"}
#   ]
# }
#
# Before fix, only C was flagged because D has incoming edge from C

# Try to enable (should fail validation)
curl -X POST http://localhost:8888/api/v1/workflows/3/enable
# Expected: {"code": 400, "message": "workflow validation failed: node is not reachable from entry node: C"}
```

### Additional Test: Circular Unreachable Subgraph

A more complex case where all nodes in the unreachable portion have incoming edges:

```bash
curl -X POST http://localhost:8888/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Circular Unreachable Test",
    "nodes": [
      {"id": "entry", "type": "script", "name": "Entry", "position": {"x": 0, "y": 0}},
      {"id": "X", "type": "script", "name": "Isolated X", "position": {"x": 100, "y": 100}},
      {"id": "Y", "type": "script", "name": "Isolated Y", "position": {"x": 200, "y": 100}},
      {"id": "Z", "type": "script", "name": "Isolated Z", "position": {"x": 150, "y": 200}}
    ],
    "edges": [
      {"id": "e1", "source": "X", "target": "Y", "outlet": "success"},
      {"id": "e2", "source": "Y", "target": "Z", "outlet": "success"},
      {"id": "e3", "source": "Z", "target": "X", "outlet": "failure"}
    ],
    "entry_node_id": "entry"
  }'
# Creates a workflow where:
# - entry is the entry node (isolated, reachable only as starting point)
# - X, Y, Z form a circular subgraph where all have incoming edges
# Note: This will also trigger CYCLE_DETECTED

curl -X POST http://localhost:8888/api/v1/workflows/4/validate
# Expected: Multiple errors including UNREACHABLE_NODE for X, Y, Z and CYCLE_DETECTED
```

---

## Code Changes Summary

### 1. `enable_workflow_logic.go`
Changed status check from:
```go
if workflow.Status != "draft" {
```
To:
```go
if workflow.Status != "draft" && workflow.Status != "disabled" {
```

### 2. `validate_workflow_logic.go`
Removed the nested condition that only flagged nodes without incoming edges:
```go
// Before (only flagged nodes with no incoming edges):
if !reachable[node.Id] && node.Id != entryNodeId {
    if len(incoming[node.Id]) == 0 && node.Id != entryNodeId {
        errors = append(...)
    }
}

// After (flags ALL unreachable nodes):
if !reachable[node.Id] && node.Id != entryNodeId {
    errors = append(...)
}
```

### 3. `workflow.api`
Updated documentation for enable endpoint:
```
@doc "Enable workflow (draft|disabled -> enabled)"
```
