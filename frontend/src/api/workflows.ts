import type { 
  WorkflowDetail, 
  WorkflowSummary, 
  SaveWorkflowRequest,
  Node,
  Edge,
  CanvasMeta,
  WorkflowStatus,
} from '../types/workflow'

const API_BASE = '/api/v1'

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
