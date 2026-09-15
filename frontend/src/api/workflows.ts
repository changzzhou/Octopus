import type { 
  WorkflowDetail, 
  WorkflowSummary, 
  SaveWorkflowRequest,
  Node,
  Edge,
  CanvasMeta,
  WorkflowStatus,
} from '../types/workflow'

/**
 * API base URL - configurable via VITE_API_BASE_URL environment variable.
 * Default: /api/v1 (proxied to backend at :8888 via Vite dev server)
 */
const API_BASE = import.meta.env.VITE_API_BASE_URL || '/api/v1'

// Re-export types
export type { WorkflowDetail, WorkflowSummary, SaveWorkflowRequest, WorkflowStatus }

// Backwards compatibility alias
export type Workflow = WorkflowSummary

/**
 * Mock mode toggle.
 * Set VITE_MOCK_API=true in .env.local or .env to enable mock data.
 * In development, you can also set window.__OCTOPUS_MOCK_API__ = true in browser console.
 */
function isMockEnabled(): boolean {
  if (typeof window !== 'undefined' && (window as any).__OCTOPUS_MOCK_API__ === true) {
    return true
  }
  return import.meta.env.VITE_MOCK_API === 'true'
}

function delay(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms))
}

// ========== Run/Step Types (aligned with BE-2/BE-3 contract) ==========

export interface WorkflowDefinitionSnapshot {
  nodes: Node[]
  edges: Edge[]
  entry_node_id?: string
  variables_schema?: string
}

export type RunStatus = 'pending' | 'running' | 'succeeded' | 'failed' | 'cancelled'
export type StepStatus = 'pending' | 'running' | 'succeeded' | 'failed' | 'skipped'

export interface Run {
  id: number
  workflow_id: number
  workflow_version: number
  definition_snapshot?: WorkflowDefinitionSnapshot
  status: RunStatus
  trigger_type: string
  started_at?: string
  finished_at?: string
  error_message?: string
  created_at: string
  updated_at: string
}

export interface Step {
  id: number
  run_id: number
  node_id: string
  node_type: string
  node_config?: string
  status: StepStatus
  input_data?: string
  output_data?: string
  error_message?: string
  started_at?: string
  finished_at?: string
  created_at: string
  updated_at: string
}

export interface TriggerRunResponse {
  run_id: number
}

export interface GetRunResponse {
  run: Run
}

export interface GetRunStepsResponse {
  steps: Step[]
}

// ========== SSE Event Types (aligned with BE-3 contract) ==========

export interface SSEEventEnvelope {
  event_id: string
  event_type: string
  run_id: number
  workflow_id: number
  occurred_at: string
  sequence?: number
  payload: unknown
}

export interface RunStatusChangedPayload {
  from_status: string
  to_status: string
  reason?: string
}

export interface StepStatusChangedPayload {
  step_id: string
  from: string
  to: string
  error_summary?: string
}

export interface RunTerminalPayload {
  final_status: string
}

export interface HumanWaitingPayload {
  step_id: string
  prompt_summary?: string
  context_refs?: string[]
}

export const SSE_EVENT_TYPES = {
  RUN_STATUS_CHANGED: 'run.status_changed',
  STEP_STATUS_CHANGED: 'step.status_changed',
  RUN_TERMINAL: 'run.terminal',
  HUMAN_WAITING: 'human.waiting',
} as const

// ========== Mock Data ==========

// Mock data for workflows list
const mockWorkflows: WorkflowSummary[] = [
  {
    id: 1,
    name: 'Data Pipeline',
    description: 'ETL workflow for daily data processing',
    status: 'enabled',
    version: 3,
    created_at: '2026-09-01T10:00:00Z',
    updated_at: '2026-09-14T08:30:00Z',
  },
  {
    id: 2,
    name: 'User Onboarding',
    description: 'Automated user welcome and setup flow',
    status: 'draft',
    version: 1,
    created_at: '2026-09-05T14:00:00Z',
    updated_at: '2026-09-13T16:45:00Z',
  },
  {
    id: 3,
    name: 'Report Generation',
    description: 'Weekly sales report automation',
    status: 'disabled',
    version: 2,
    created_at: '2026-09-08T09:00:00Z',
    updated_at: '2026-09-12T11:20:00Z',
  },
  {
    id: 4,
    name: 'Notification Service',
    description: 'Multi-channel notification dispatch',
    status: 'enabled',
    version: 5,
    created_at: '2026-09-10T11:00:00Z',
    updated_at: '2026-09-14T07:00:00Z',
  },
]

// Mock detail data (with nodes/edges)
const mockWorkflowDetails: Record<number, WorkflowDetail> = {
  1: {
    id: 1,
    name: 'Data Pipeline',
    description: 'ETL workflow for daily data processing',
    status: 'enabled',
    version: 3,
    nodes: [
      { id: 'script-1', type: 'script', name: 'Extract Data', position: { x: 100, y: 100 }, config: '{"language":"python","code":"# Extract data","timeout":30000}' },
      { id: 'http-1', type: 'http', name: 'Transform API', position: { x: 350, y: 100 }, config: '{"method":"POST","url":"https://api.example.com/transform","timeout":30000}' },
    ],
    edges: [
      { id: 'e-1', source: 'script-1', target: 'http-1', outlet: 'success' },
    ],
    created_at: '2026-09-01T10:00:00Z',
    updated_at: '2026-09-14T08:30:00Z',
  },
  2: {
    id: 2,
    name: 'User Onboarding',
    description: 'Automated user welcome and setup flow',
    status: 'draft',
    version: 1,
    nodes: [],
    edges: [],
    created_at: '2026-09-05T14:00:00Z',
    updated_at: '2026-09-13T16:45:00Z',
  },
  3: {
    id: 3,
    name: 'Report Generation',
    description: 'Weekly sales report automation',
    status: 'disabled',
    version: 2,
    nodes: [
      { id: 'script-1', type: 'script', name: 'Generate Report', position: { x: 100, y: 100 }, config: '{"language":"javascript","code":"// Generate report","timeout":60000}' },
      { id: 'human-1', type: 'human', name: 'Review Report', position: { x: 350, y: 100 }, config: '{"instructions":"Please review the generated report","assignee":"reviewer@example.com"}' },
    ],
    edges: [
      { id: 'e-1', source: 'script-1', target: 'human-1', outlet: 'success' },
    ],
    created_at: '2026-09-08T09:00:00Z',
    updated_at: '2026-09-12T11:20:00Z',
  },
  4: {
    id: 4,
    name: 'Notification Service',
    description: 'Multi-channel notification dispatch',
    status: 'enabled',
    version: 5,
    nodes: [],
    edges: [],
    created_at: '2026-09-10T11:00:00Z',
    updated_at: '2026-09-14T07:00:00Z',
  },
}

let mockIdCounter = 100
let mockRunIdCounter = 1000
let mockStepIdCounter = 10000

const mockRuns: Map<number, Run[]> = new Map([
  [1, [
    {
      id: 1,
      workflow_id: 1,
      workflow_version: 3,
      definition_snapshot: {
        nodes: [
          { id: 'start', type: 'script', name: 'Initialize', position: { x: 100, y: 100 } },
          { id: 'fetch', type: 'http', name: 'Fetch Data', position: { x: 250, y: 100 } },
          { id: 'transform', type: 'script', name: 'Transform', position: { x: 400, y: 100 } },
          { id: 'notify', type: 'http', name: 'Send Notification', position: { x: 550, y: 100 } },
        ],
        edges: [
          { id: 'e1', source: 'start', target: 'fetch', outlet: 'success' },
          { id: 'e2', source: 'fetch', target: 'transform', outlet: 'success' },
          { id: 'e3', source: 'transform', target: 'notify', outlet: 'success' },
        ],
        entry_node_id: 'start',
      },
      status: 'succeeded',
      trigger_type: 'manual',
      started_at: '2026-09-14T08:00:00Z',
      finished_at: '2026-09-14T08:05:32Z',
      created_at: '2026-09-14T08:00:00Z',
      updated_at: '2026-09-14T08:05:32Z',
    },
    {
      id: 2,
      workflow_id: 1,
      workflow_version: 3,
      definition_snapshot: {
        nodes: [
          { id: 'start', type: 'script', name: 'Initialize', position: { x: 100, y: 100 } },
          { id: 'fetch', type: 'http', name: 'Fetch Data', position: { x: 250, y: 100 } },
          { id: 'transform', type: 'script', name: 'Transform', position: { x: 400, y: 100 } },
          { id: 'notify', type: 'http', name: 'Send Notification', position: { x: 550, y: 100 } },
        ],
        edges: [
          { id: 'e1', source: 'start', target: 'fetch', outlet: 'success' },
          { id: 'e2', source: 'fetch', target: 'transform', outlet: 'success' },
          { id: 'e3', source: 'transform', target: 'notify', outlet: 'success' },
        ],
        entry_node_id: 'start',
      },
      status: 'running',
      trigger_type: 'manual',
      started_at: '2026-09-14T10:30:00Z',
      created_at: '2026-09-14T10:30:00Z',
      updated_at: '2026-09-14T10:32:15Z',
    },
  ]],
  [2, [
    {
      id: 3,
      workflow_id: 2,
      workflow_version: 1,
      definition_snapshot: {
        nodes: [
          { id: 'welcome', type: 'http', name: 'Send Welcome Email', position: { x: 100, y: 100 } },
          { id: 'setup', type: 'script', name: 'Setup Account', position: { x: 300, y: 100 } },
        ],
        edges: [
          { id: 'e1', source: 'welcome', target: 'setup', outlet: 'success' },
        ],
        entry_node_id: 'welcome',
      },
      status: 'failed',
      trigger_type: 'trial',
      started_at: '2026-09-13T16:00:00Z',
      finished_at: '2026-09-13T16:01:45Z',
      error_message: 'Step "Setup Account" failed: invalid user configuration',
      created_at: '2026-09-13T16:00:00Z',
      updated_at: '2026-09-13T16:01:45Z',
    },
  ]],
])

const mockSteps: Map<number, Step[]> = new Map([
  [1, [
    { id: 1, run_id: 1, node_id: 'start', node_type: 'script', status: 'succeeded', started_at: '2026-09-14T08:00:00Z', finished_at: '2026-09-14T08:00:05Z', created_at: '2026-09-14T08:00:00Z', updated_at: '2026-09-14T08:00:05Z' },
    { id: 2, run_id: 1, node_id: 'fetch', node_type: 'http', status: 'succeeded', started_at: '2026-09-14T08:00:05Z', finished_at: '2026-09-14T08:02:30Z', created_at: '2026-09-14T08:00:05Z', updated_at: '2026-09-14T08:02:30Z' },
    { id: 3, run_id: 1, node_id: 'transform', node_type: 'script', status: 'succeeded', started_at: '2026-09-14T08:02:30Z', finished_at: '2026-09-14T08:04:00Z', created_at: '2026-09-14T08:02:30Z', updated_at: '2026-09-14T08:04:00Z' },
    { id: 4, run_id: 1, node_id: 'notify', node_type: 'http', status: 'succeeded', started_at: '2026-09-14T08:04:00Z', finished_at: '2026-09-14T08:05:32Z', created_at: '2026-09-14T08:04:00Z', updated_at: '2026-09-14T08:05:32Z' },
  ]],
  [2, [
    { id: 5, run_id: 2, node_id: 'start', node_type: 'script', status: 'succeeded', started_at: '2026-09-14T10:30:00Z', finished_at: '2026-09-14T10:30:08Z', created_at: '2026-09-14T10:30:00Z', updated_at: '2026-09-14T10:30:08Z' },
    { id: 6, run_id: 2, node_id: 'fetch', node_type: 'http', status: 'running', started_at: '2026-09-14T10:30:08Z', created_at: '2026-09-14T10:30:08Z', updated_at: '2026-09-14T10:32:15Z' },
    { id: 7, run_id: 2, node_id: 'transform', node_type: 'script', status: 'pending', created_at: '2026-09-14T10:30:00Z', updated_at: '2026-09-14T10:30:00Z' },
    { id: 8, run_id: 2, node_id: 'notify', node_type: 'http', status: 'pending', created_at: '2026-09-14T10:30:00Z', updated_at: '2026-09-14T10:30:00Z' },
  ]],
  [3, [
    { id: 9, run_id: 3, node_id: 'welcome', node_type: 'http', status: 'succeeded', started_at: '2026-09-13T16:00:00Z', finished_at: '2026-09-13T16:00:30Z', created_at: '2026-09-13T16:00:00Z', updated_at: '2026-09-13T16:00:30Z' },
    { id: 10, run_id: 3, node_id: 'setup', node_type: 'script', status: 'failed', error_message: 'invalid user configuration', started_at: '2026-09-13T16:00:30Z', finished_at: '2026-09-13T16:01:45Z', created_at: '2026-09-13T16:00:30Z', updated_at: '2026-09-13T16:01:45Z' },
  ]],
])

// ========== Workflow List API ==========

export interface ListWorkflowsResponse {
  total: number
  workflows: WorkflowSummary[]
}

export async function listWorkflows(page = 1, pageSize = 20, status?: string): Promise<ListWorkflowsResponse> {
  if (isMockEnabled()) {
    await delay(300)
    let filtered = mockWorkflows
    if (status) {
      filtered = mockWorkflows.filter(w => w.status === status)
    }
    const start = (page - 1) * pageSize
    const end = start + pageSize
    return {
      total: filtered.length,
      workflows: filtered.slice(start, end),
    }
  }
  
  let url = `${API_BASE}/workflows?page=${page}&page_size=${pageSize}`
  if (status) {
    url += `&status=${status}`
  }
  const res = await fetch(url)
  if (!res.ok) throw new Error('Failed to fetch workflows')
  return res.json()
}

export async function getWorkflow(id: number): Promise<{ workflow: WorkflowSummary }> {
  if (isMockEnabled()) {
    await delay(200)
    const workflow = mockWorkflows.find(w => w.id === id)
    if (!workflow) throw new Error('Workflow not found')
    return { workflow }
  }
  
  const res = await fetch(`${API_BASE}/workflows/${id}`)
  if (!res.ok) throw new Error('Failed to fetch workflow')
  return res.json()
}

export async function getWorkflowDetail(id: number): Promise<{ workflow: WorkflowDetail }> {
  if (isMockEnabled()) {
    await delay(200)
    const workflow = mockWorkflowDetails[id]
    if (!workflow) {
      // Fallback to summary data if detail not found
      const summary = mockWorkflows.find(w => w.id === id)
      if (!summary) throw new Error('Workflow not found')
      return {
        workflow: {
          ...summary,
          nodes: [],
          edges: [],
        }
      }
    }
    return { workflow: { ...workflow } }
  }
  
  const res = await fetch(`${API_BASE}/workflows/${id}`)
  if (!res.ok) throw new Error('Failed to fetch workflow')
  const data = await res.json()
  
  // BE-1 returns WorkflowDetail with nodes/edges directly
  const workflow: WorkflowDetail = {
    id: data.workflow.id,
    name: data.workflow.name,
    description: data.workflow.description,
    status: data.workflow.status || 'draft',
    version: data.workflow.version || 1,
    nodes: data.workflow.nodes || [],
    edges: data.workflow.edges || [],
    entry_node_id: data.workflow.entry_node_id,
    variables_schema: data.workflow.variables_schema,
    canvas_meta: data.workflow.canvas_meta,
    created_by: data.workflow.created_by,
    updated_by: data.workflow.updated_by,
    created_at: data.workflow.created_at,
    updated_at: data.workflow.updated_at,
  }
  return { workflow }
}

export async function saveWorkflow(id: number, data: SaveWorkflowRequest): Promise<{ version: number }> {
  if (isMockEnabled()) {
    await delay(300)
    const workflow = mockWorkflowDetails[id]
    if (workflow) {
      workflow.nodes = data.nodes || []
      workflow.edges = data.edges || []
      workflow.canvas_meta = data.canvas_meta
      workflow.version = (workflow.version || 0) + 1
      workflow.updated_at = new Date().toISOString()
      
      // Update summary too
      const summary = mockWorkflows.find(w => w.id === id)
      if (summary) {
        summary.version = workflow.version
        summary.updated_at = workflow.updated_at
      }
      
      return { version: workflow.version }
    }
    throw new Error('Workflow not found')
  }
  
  const res = await fetch(`${API_BASE}/workflows/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error.message || 'Failed to save workflow')
  }
  return res.json()
}

export interface CreateWorkflowRequest {
  name: string
  description?: string
  nodes?: Node[]
  edges?: Edge[]
  entry_node_id?: string
  variables_schema?: string
  canvas_meta?: CanvasMeta
}

export async function createWorkflow(name: string, description?: string): Promise<{ id: number; version: number; status: string }> {
  if (isMockEnabled()) {
    await delay(300)
    const newId = ++mockIdCounter
    const now = new Date().toISOString()
    
    const newWorkflow: WorkflowSummary = {
      id: newId,
      name,
      description: description || '',
      status: 'draft',
      version: 1,
      created_at: now,
      updated_at: now,
    }
    mockWorkflows.unshift(newWorkflow)
    
    mockWorkflowDetails[newId] = {
      ...newWorkflow,
      nodes: [],
      edges: [],
    }
    
    return { id: newId, version: 1, status: 'draft' }
  }
  
  const res = await fetch(`${API_BASE}/workflows`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, description }),
  })
  if (!res.ok) throw new Error('Failed to create workflow')
  return res.json()
}

export async function deleteWorkflow(id: number): Promise<void> {
  if (isMockEnabled()) {
    await delay(200)
    const index = mockWorkflows.findIndex(w => w.id === id)
    if (index !== -1) mockWorkflows.splice(index, 1)
    delete mockWorkflowDetails[id]
    return
  }
  
  const res = await fetch(`${API_BASE}/workflows/${id}`, { method: 'DELETE' })
  if (!res.ok) throw new Error('Failed to delete workflow')
}

export async function enableWorkflow(id: number): Promise<{ status: string; version: number }> {
  if (isMockEnabled()) {
    await delay(200)
    const workflow = mockWorkflows.find(w => w.id === id)
    const detail = mockWorkflowDetails[id]
    if (workflow) {
      workflow.status = 'enabled'
      workflow.version = (workflow.version || 0) + 1
      workflow.updated_at = new Date().toISOString()
      if (detail) {
        detail.status = 'enabled'
        detail.version = workflow.version
        detail.updated_at = workflow.updated_at
      }
      return { status: 'enabled', version: workflow.version }
    }
    throw new Error('Workflow not found')
  }
  
  const res = await fetch(`${API_BASE}/workflows/${id}/enable`, { method: 'POST' })
  if (!res.ok) throw new Error('Failed to enable workflow')
  return res.json()
}

export async function disableWorkflow(id: number): Promise<{ status: string }> {
  if (isMockEnabled()) {
    await delay(200)
    const workflow = mockWorkflows.find(w => w.id === id)
    const detail = mockWorkflowDetails[id]
    if (workflow) {
      workflow.status = 'disabled'
      workflow.updated_at = new Date().toISOString()
      if (detail) {
        detail.status = 'disabled'
        detail.updated_at = workflow.updated_at
      }
      return { status: 'disabled' }
    }
    throw new Error('Workflow not found')
  }
  
  const res = await fetch(`${API_BASE}/workflows/${id}/disable`, { method: 'POST' })
  if (!res.ok) throw new Error('Failed to disable workflow')
  return res.json()
}

export async function checkHealth(): Promise<{ status: string; version: string }> {
  if (isMockEnabled()) {
    await delay(100)
    return { status: 'ok', version: 'v1.0.0-mock' }
  }
  
  const res = await fetch(`${API_BASE}/health`)
  if (!res.ok) throw new Error('Health check failed')
  return res.json()
}

// ========== Run/Step API (FE-3) ==========

/**
 * Trigger a workflow run (works for draft workflows too - trial run with frozen contract)
 */
export async function triggerRun(workflowId: number): Promise<TriggerRunResponse> {
  if (isMockEnabled()) {
    await delay(400)
    const workflow = mockWorkflows.find(w => w.id === workflowId)
    if (!workflow) throw new Error('Workflow not found')
    
    const newRunId = ++mockRunIdCounter
    const newRun: Run = {
      id: newRunId,
      workflow_id: workflowId,
      workflow_version: 1,
      definition_snapshot: {
        nodes: [
          { id: 'node1', type: 'script', name: 'Step 1', position: { x: 100, y: 100 } },
          { id: 'node2', type: 'http', name: 'Step 2', position: { x: 250, y: 100 } },
        ],
        edges: [
          { id: 'e1', source: 'node1', target: 'node2', outlet: 'success' },
        ],
        entry_node_id: 'node1',
      },
      status: 'pending',
      trigger_type: workflow.status === 'draft' ? 'trial' : 'manual',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }
    
    const existingRuns = mockRuns.get(workflowId) || []
    existingRuns.unshift(newRun)
    mockRuns.set(workflowId, existingRuns)
    
    const initialSteps: Step[] = [
      { id: ++mockStepIdCounter, run_id: newRunId, node_id: 'node1', node_type: 'script', status: 'pending', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
      { id: ++mockStepIdCounter, run_id: newRunId, node_id: 'node2', node_type: 'http', status: 'pending', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
    ]
    mockSteps.set(newRunId, initialSteps)
    
    simulateRunExecution(newRunId)
    
    return { run_id: newRunId }
  }
  
  const res = await fetch(`${API_BASE}/workflows/${workflowId}/runs`, { method: 'POST' })
  if (!res.ok) throw new Error('Failed to trigger run')
  return res.json()
}

/**
 * Simulate run execution for mock mode (updates run/step status over time)
 */
function simulateRunExecution(runId: number) {
  const runSteps = mockSteps.get(runId)
  if (!runSteps) return
  
  let allRuns: Run[] = []
  mockRuns.forEach(runs => allRuns = allRuns.concat(runs))
  const targetRun = allRuns.find(r => r.id === runId)
  if (!targetRun) return
  
  let currentStepIndex = 0
  
  function executeNextStep() {
    if (!runSteps || !targetRun) return
    
    if (currentStepIndex >= runSteps.length) {
      targetRun.status = 'succeeded'
      targetRun.finished_at = new Date().toISOString()
      targetRun.updated_at = new Date().toISOString()
      return
    }
    
    const step = runSteps[currentStepIndex]
    
    if (currentStepIndex === 0) {
      targetRun.status = 'running'
      targetRun.started_at = new Date().toISOString()
    }
    
    step.status = 'running'
    step.started_at = new Date().toISOString()
    step.updated_at = new Date().toISOString()
    targetRun.updated_at = new Date().toISOString()
    
    setTimeout(() => {
      if (!targetRun) return
      
      const shouldFail = Math.random() < 0.1
      if (shouldFail) {
        step.status = 'failed'
        step.error_message = 'Simulated failure for demonstration'
        step.finished_at = new Date().toISOString()
        step.updated_at = new Date().toISOString()
        
        targetRun.status = 'failed'
        targetRun.error_message = `Step "${step.node_id}" failed: ${step.error_message}`
        targetRun.finished_at = new Date().toISOString()
        targetRun.updated_at = new Date().toISOString()
        return
      }
      
      step.status = 'succeeded'
      step.finished_at = new Date().toISOString()
      step.updated_at = new Date().toISOString()
      
      currentStepIndex++
      executeNextStep()
    }, 1500 + Math.random() * 2000)
  }
  
  setTimeout(executeNextStep, 500)
}

/**
 * Get run details by run ID
 */
export async function getRun(runId: number): Promise<GetRunResponse> {
  if (isMockEnabled()) {
    await delay(150)
    let foundRun: Run | undefined
    mockRuns.forEach(runs => {
      const r = runs.find(run => run.id === runId)
      if (r) foundRun = r
    })
    if (!foundRun) throw new Error('Run not found')
    return { run: foundRun }
  }
  
  const res = await fetch(`${API_BASE}/runs/${runId}`)
  if (!res.ok) throw new Error('Failed to fetch run')
  return res.json()
}

/**
 * Get all steps for a run
 */
export async function getRunSteps(runId: number): Promise<GetRunStepsResponse> {
  if (isMockEnabled()) {
    await delay(150)
    const steps = mockSteps.get(runId) || []
    return { steps }
  }
  
  const res = await fetch(`${API_BASE}/runs/${runId}/steps`)
  if (!res.ok) throw new Error('Failed to fetch run steps')
  return res.json()
}

/**
 * Get runs for a workflow
 */
export async function getWorkflowRuns(workflowId: number): Promise<{ runs: Run[] }> {
  if (isMockEnabled()) {
    await delay(200)
    const runs = mockRuns.get(workflowId) || []
    return { runs }
  }
  
  const res = await fetch(`${API_BASE}/workflows/${workflowId}/runs`)
  if (!res.ok) throw new Error('Failed to fetch workflow runs')
  return res.json()
}

// ========== SSE Subscription (BE-3 integration) ==========

export type SSEEventHandler = (event: SSEEventEnvelope) => void

export interface SSESubscription {
  close: () => void
}

/**
 * Subscribe to run events via SSE (Server-Sent Events).
 * Aligned with BE-3 contract: GET /api/v1/runs/:runId/events
 * 
 * Event types:
 * - run.status_changed: Run status transition
 * - step.status_changed: Step status transition
 * - run.terminal: Run reached terminal state
 * - human.waiting: Human task waiting for input
 * 
 * @param runId - The run ID to subscribe to
 * @param onEvent - Callback for each event
 * @param onError - Optional error callback
 * @returns Subscription object with close() method
 */
export function subscribeRunEvents(
  runId: number,
  onEvent: SSEEventHandler,
  onError?: (error: Event) => void
): SSESubscription {
  if (isMockEnabled()) {
    return createMockSSESubscription(runId, onEvent)
  }
  
  const eventSource = new EventSource(`${API_BASE}/runs/${runId}/events`)
  
  const handleEvent = (e: MessageEvent) => {
    try {
      const envelope: SSEEventEnvelope = JSON.parse(e.data)
      onEvent(envelope)
    } catch (err) {
      console.error('Failed to parse SSE event:', err)
    }
  }
  
  eventSource.addEventListener(SSE_EVENT_TYPES.RUN_STATUS_CHANGED, handleEvent)
  eventSource.addEventListener(SSE_EVENT_TYPES.STEP_STATUS_CHANGED, handleEvent)
  eventSource.addEventListener(SSE_EVENT_TYPES.RUN_TERMINAL, handleEvent)
  eventSource.addEventListener(SSE_EVENT_TYPES.HUMAN_WAITING, handleEvent)
  
  eventSource.onerror = (err) => {
    console.error('SSE connection error:', err)
    onError?.(err)
  }
  
  return {
    close: () => {
      eventSource.close()
    }
  }
}

/**
 * Mock SSE subscription for development
 */
function createMockSSESubscription(runId: number, onEvent: SSEEventHandler): SSESubscription {
  let active = true
  let sequence = 0
  
  const interval = setInterval(() => {
    if (!active) return
    
    let foundRun: Run | undefined
    mockRuns.forEach(runs => {
      const r = runs.find(run => run.id === runId)
      if (r) foundRun = r
    })
    
    if (!foundRun) return
    
    const steps = mockSteps.get(runId) || []
    const runningStep = steps.find(s => s.status === 'running')
    
    if (runningStep) {
      onEvent({
        event_id: `mock-${Date.now()}`,
        event_type: SSE_EVENT_TYPES.STEP_STATUS_CHANGED,
        run_id: runId,
        workflow_id: foundRun.workflow_id,
        occurred_at: new Date().toISOString(),
        sequence: ++sequence,
        payload: {
          step_id: runningStep.node_id,
          from: 'pending',
          to: 'running',
        } as StepStatusChangedPayload,
      })
    }
    
    if (['succeeded', 'failed', 'cancelled'].includes(foundRun.status)) {
      onEvent({
        event_id: `mock-${Date.now()}`,
        event_type: SSE_EVENT_TYPES.RUN_TERMINAL,
        run_id: runId,
        workflow_id: foundRun.workflow_id,
        occurred_at: new Date().toISOString(),
        sequence: ++sequence,
        payload: {
          final_status: foundRun.status,
        } as RunTerminalPayload,
      })
      clearInterval(interval)
    }
  }, 2000)
  
  return {
    close: () => {
      active = false
      clearInterval(interval)
    }
  }
}

/**
 * Polling fallback for execution state updates.
 * Use this when SSE is not available or as a backup.
 * 
 * @param runId - The run ID to poll
 * @param onUpdate - Callback for status updates
 * @param intervalMs - Polling interval (default 2000ms)
 * @returns Cleanup function to stop polling
 */
export function pollRunStatus(
  runId: number,
  onUpdate: (run: Run, steps: Step[]) => void,
  intervalMs = 2000
): () => void {
  let active = true
  
  async function poll() {
    if (!active) return
    
    try {
      const [runRes, stepsRes] = await Promise.all([
        getRun(runId),
        getRunSteps(runId),
      ])
      
      if (active) {
        onUpdate(runRes.run, stepsRes.steps)
        
        const isTerminal = ['succeeded', 'failed', 'cancelled'].includes(runRes.run.status)
        if (!isTerminal) {
          setTimeout(poll, intervalMs)
        }
      }
    } catch (e) {
      console.error('Polling error:', e)
      if (active) {
        setTimeout(poll, intervalMs * 2)
      }
    }
  }
  
  poll()
  
  return () => {
    active = false
  }
}

/**
 * Live subscription manager: SSE-first with polling fallback.
 * 
 * Strategy:
 * 1. Attempt SSE connection to /api/v1/runs/:runId/events
 * 2. On SSE error, fall back to polling
 * 3. On SSE terminal event, close connection
 * 4. Always fetch full state on connect for consistency
 * 
 * @param runId - The run ID to subscribe to
 * @param onUpdate - Callback for state updates (run + steps)
 * @param options - Configuration options
 * @returns Cleanup function to stop all subscriptions
 */
export interface LiveSubscriptionOptions {
  pollIntervalMs?: number
  onConnectionChange?: (mode: 'sse' | 'polling' | 'disconnected') => void
}

export function subscribeRunLive(
  runId: number,
  onUpdate: (run: Run, steps: Step[]) => void,
  options: LiveSubscriptionOptions = {}
): () => void {
  const { pollIntervalMs = 2000, onConnectionChange } = options
  
  let active = true
  let sseSubscription: SSESubscription | null = null
  let pollCleanup: (() => void) | null = null
  let connectionMode: 'sse' | 'polling' | 'disconnected' = 'disconnected'
  
  function setMode(mode: 'sse' | 'polling' | 'disconnected') {
    if (connectionMode !== mode) {
      connectionMode = mode
      onConnectionChange?.(mode)
    }
  }
  
  async function fetchFullState(): Promise<{ run: Run; steps: Step[] } | null> {
    try {
      const [runRes, stepsRes] = await Promise.all([
        getRun(runId),
        getRunSteps(runId),
      ])
      return { run: runRes.run, steps: stepsRes.steps }
    } catch (e) {
      console.error('Failed to fetch run state:', e)
      return null
    }
  }
  
  function startPolling() {
    if (!active || pollCleanup) return
    
    setMode('polling')
    pollCleanup = pollRunStatus(runId, (run, steps) => {
      if (!active) return
      onUpdate(run, steps)
      
      if (['succeeded', 'failed', 'cancelled'].includes(run.status)) {
        cleanup()
      }
    }, pollIntervalMs)
  }
  
  function startSSE() {
    if (!active || isMockEnabled()) {
      startPolling()
      return
    }
    
    setMode('sse')
    
    sseSubscription = subscribeRunEvents(
      runId,
      async (event) => {
        if (!active) return
        
        const state = await fetchFullState()
        if (state && active) {
          onUpdate(state.run, state.steps)
        }
        
        if (event.event_type === SSE_EVENT_TYPES.RUN_TERMINAL) {
          cleanup()
        }
      },
      (_error) => {
        if (!active) return
        console.warn('SSE connection failed, falling back to polling')
        
        if (sseSubscription) {
          sseSubscription.close()
          sseSubscription = null
        }
        
        startPolling()
      }
    )
  }
  
  function cleanup() {
    active = false
    setMode('disconnected')
    
    if (sseSubscription) {
      sseSubscription.close()
      sseSubscription = null
    }
    
    if (pollCleanup) {
      pollCleanup()
      pollCleanup = null
    }
  }
  
  fetchFullState().then((state) => {
    if (!active) return
    
    if (state) {
      onUpdate(state.run, state.steps)
      
      if (['succeeded', 'failed', 'cancelled'].includes(state.run.status)) {
        cleanup()
        return
      }
    }
    
    startSSE()
  })
  
  return cleanup
}
