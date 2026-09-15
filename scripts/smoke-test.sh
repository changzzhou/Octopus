#!/bin/bash
# Octopus E2E Smoke Test Script
# Tests the complete workflow: Definition → Run → SSE

set -e

BASE_URL="${BASE_URL:-http://localhost:8888}"
OUTPUT_DIR="${OUTPUT_DIR:-/tmp/octopus-smoke-test}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

mkdir -p "$OUTPUT_DIR"

log_step() {
    echo -e "\n${YELLOW}=== $1 ===${NC}"
}

log_pass() {
    echo -e "${GREEN}[PASS]${NC} $1"
}

log_fail() {
    echo -e "${RED}[FAIL]${NC} $1"
}

check_json() {
    local response="$1"
    local expected_key="$2"
    
    if echo "$response" | python3 -c "import sys,json; json.load(sys.stdin)['$expected_key']" 2>/dev/null; then
        return 0
    fi
    return 1
}

# Test 1: Health Check
log_step "1. Health Check"
HEALTH=$(curl -s "$BASE_URL/api/v1/health")
echo "$HEALTH" | tee "$OUTPUT_DIR/01-health.json"

if echo "$HEALTH" | grep -q '"status":"ok"'; then
    log_pass "Health check passed"
else
    log_fail "Health check failed"
    exit 1
fi

# Test 2: Create Workflow
log_step "2. Create Workflow"
CREATE_RESP=$(curl -s -X POST "$BASE_URL/api/v1/workflows" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Smoke Test Workflow",
    "description": "Automated smoke test workflow",
    "nodes": [
      {"id": "node_1", "type": "script", "name": "Script Node", "position": {"x": 100, "y": 100}},
      {"id": "node_2", "type": "http", "name": "HTTP Node", "position": {"x": 300, "y": 100}}
    ],
    "edges": [
      {"id": "edge_1", "source": "node_1", "target": "node_2", "outlet": "success"}
    ],
    "entry_node_id": "node_1"
  }')
echo "$CREATE_RESP" | tee "$OUTPUT_DIR/02-create-workflow.json"

WORKFLOW_ID=$(echo "$CREATE_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])" 2>/dev/null || echo "")
if [ -n "$WORKFLOW_ID" ]; then
    log_pass "Created workflow ID: $WORKFLOW_ID"
else
    log_fail "Failed to create workflow"
    exit 1
fi

# Test 3: Get Workflow
log_step "3. Get Workflow"
GET_RESP=$(curl -s "$BASE_URL/api/v1/workflows/$WORKFLOW_ID")
echo "$GET_RESP" | python3 -m json.tool | tee "$OUTPUT_DIR/03-get-workflow.json"

if echo "$GET_RESP" | grep -q '"status":"draft"'; then
    log_pass "Workflow retrieved, status: draft"
else
    log_fail "Failed to get workflow"
    exit 1
fi

# Test 4: Validate Workflow
log_step "4. Validate Workflow"
VALIDATE_RESP=$(curl -s -X POST "$BASE_URL/api/v1/workflows/$WORKFLOW_ID/validate")
echo "$VALIDATE_RESP" | tee "$OUTPUT_DIR/04-validate-workflow.json"

if echo "$VALIDATE_RESP" | grep -q '"valid":true'; then
    log_pass "Workflow validation passed"
else
    log_fail "Workflow validation failed"
    echo "$VALIDATE_RESP"
    exit 1
fi

# Test 5: Enable Workflow
log_step "5. Enable Workflow"
ENABLE_RESP=$(curl -s -X POST "$BASE_URL/api/v1/workflows/$WORKFLOW_ID/enable")
echo "$ENABLE_RESP" | tee "$OUTPUT_DIR/05-enable-workflow.json"

if echo "$ENABLE_RESP" | grep -q '"status":"enabled"'; then
    log_pass "Workflow enabled"
else
    log_fail "Failed to enable workflow"
    exit 1
fi

# Test 6: Trigger Run
log_step "6. Trigger Run"
RUN_RESP=$(curl -s -X POST "$BASE_URL/api/v1/workflows/$WORKFLOW_ID/runs")
echo "$RUN_RESP" | tee "$OUTPUT_DIR/06-trigger-run.json"

RUN_ID=$(echo "$RUN_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['run_id'])" 2>/dev/null || echo "")
if [ -n "$RUN_ID" ]; then
    log_pass "Triggered run ID: $RUN_ID"
else
    log_fail "Failed to trigger run"
    exit 1
fi

# Test 7: Wait and Check Run Status
log_step "7. Check Run Status (waiting for completion)"
MAX_WAIT=30
WAIT=0
while [ $WAIT -lt $MAX_WAIT ]; do
    RUN_STATUS=$(curl -s "$BASE_URL/api/v1/runs/$RUN_ID")
    STATUS=$(echo "$RUN_STATUS" | python3 -c "import sys,json; print(json.load(sys.stdin)['run']['status'])" 2>/dev/null || echo "unknown")
    
    if [ "$STATUS" = "succeeded" ] || [ "$STATUS" = "failed" ]; then
        break
    fi
    
    echo "  Status: $STATUS (waiting...)"
    sleep 1
    WAIT=$((WAIT + 1))
done

echo "$RUN_STATUS" | python3 -m json.tool | tee "$OUTPUT_DIR/07-run-status.json"

if [ "$STATUS" = "succeeded" ]; then
    log_pass "Run completed successfully"
else
    log_fail "Run did not succeed (status: $STATUS)"
    exit 1
fi

# Test 8: Check Run Steps
log_step "8. Check Run Steps"
STEPS_RESP=$(curl -s "$BASE_URL/api/v1/runs/$RUN_ID/steps")
echo "$STEPS_RESP" | python3 -m json.tool | tee "$OUTPUT_DIR/08-run-steps.json"

STEP_COUNT=$(echo "$STEPS_RESP" | python3 -c "import sys,json; print(len(json.load(sys.stdin)['steps']))" 2>/dev/null || echo "0")
SUCCEEDED_COUNT=$(echo "$STEPS_RESP" | python3 -c "import sys,json; print(len([s for s in json.load(sys.stdin)['steps'] if s['status']=='succeeded']))" 2>/dev/null || echo "0")

if [ "$STEP_COUNT" = "2" ] && [ "$SUCCEEDED_COUNT" = "2" ]; then
    log_pass "All 2 steps completed successfully"
else
    log_fail "Step check failed (total: $STEP_COUNT, succeeded: $SUCCEEDED_COUNT)"
    exit 1
fi

# Test 9: SSE Events (new run)
log_step "9. SSE Events Test"

# Trigger another run for SSE test
RUN_RESP2=$(curl -s -X POST "$BASE_URL/api/v1/workflows/$WORKFLOW_ID/runs")
RUN_ID2=$(echo "$RUN_RESP2" | python3 -c "import sys,json; print(json.load(sys.stdin)['run_id'])" 2>/dev/null || echo "")

if [ -z "$RUN_ID2" ]; then
    log_fail "Failed to trigger run for SSE test"
    exit 1
fi

echo "Triggered run $RUN_ID2 for SSE capture..."

# Capture SSE events with timeout
SSE_OUTPUT="$OUTPUT_DIR/09-sse-events.txt"
timeout 10 curl -N -s "$BASE_URL/api/v1/runs/$RUN_ID2/events" > "$SSE_OUTPUT" 2>&1 || true

echo "SSE Events captured:"
cat "$SSE_OUTPUT"
echo ""

# Check for expected events
if grep -q "run.status_changed" "$SSE_OUTPUT" && grep -q "run.terminal" "$SSE_OUTPUT"; then
    log_pass "SSE events captured: run.status_changed, run.terminal"
elif grep -q "step.status_changed" "$SSE_OUTPUT"; then
    log_pass "SSE events captured: step.status_changed (run may have completed before subscription)"
else
    log_fail "SSE events not captured as expected"
    exit 1
fi

# Test 10: Disable Workflow
log_step "10. Disable Workflow"
DISABLE_RESP=$(curl -s -X POST "$BASE_URL/api/v1/workflows/$WORKFLOW_ID/disable")
echo "$DISABLE_RESP" | tee "$OUTPUT_DIR/10-disable-workflow.json"

if echo "$DISABLE_RESP" | grep -q '"status":"disabled"'; then
    log_pass "Workflow disabled"
else
    log_fail "Failed to disable workflow"
    exit 1
fi

# Summary
log_step "SUMMARY"
echo -e "${GREEN}"
echo "=========================================="
echo "  ALL SMOKE TESTS PASSED!"
echo "=========================================="
echo -e "${NC}"
echo ""
echo "Test artifacts saved to: $OUTPUT_DIR"
echo ""
echo "Files:"
ls -la "$OUTPUT_DIR"
