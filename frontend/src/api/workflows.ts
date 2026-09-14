import type { 
  WorkflowDetail, 
  WorkflowSummary, 
  SaveWorkflowRequest,
  Node,
  Edge,
  CanvasMeta,
} from '../types/workflow'

const API_BASE = '/api/v1'

// Re-export types for backward compatibility
export type { WorkflowDetail, WorkflowSummary, SaveWorkflowRequest }

export interface ListWorkflowsResponse {
  total: number
  workflows: WorkflowSummary[]
}

export async function listWorkflows(page = 1, pageSize = 20, status?: string): Promise<ListWorkflowsResponse> {
  let url = `${API_BASE}/workflows?page=${page}&page_size=${pageSize}`
  if (status) {
    url += `&status=${status}`
  }
  const res = await fetch(url)
  if (!res.ok) throw new Error('Failed to fetch workflows')
  return res.json()
}

export async function getWorkflow(id: number): Promise<{ workflow: WorkflowSummary }> {
  const res = await fetch(`${API_BASE}/workflows/${id}`)
  if (!res.ok) throw new Error('Failed to fetch workflow')
  return res.json()
}

export async function getWorkflowDetail(id: number): Promise<{ workflow: WorkflowDetail }> {
  const res = await fetch(`${API_BASE}/workflows/${id}`)
  if (!res.ok) throw new Error('Failed to fetch workflow')
  const data = await res.json()
  
  // BE-1 returns WorkflowDetail with nodes/edges directly (not canvas wrapper)
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

export async function createWorkflow(data: CreateWorkflowRequest): Promise<{ id: number; version: number; status: string }> {
  const res = await fetch(`${API_BASE}/workflows`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) throw new Error('Failed to create workflow')
  return res.json()
}

export async function deleteWorkflow(id: number): Promise<void> {
  const res = await fetch(`${API_BASE}/workflows/${id}`, { method: 'DELETE' })
  if (!res.ok) throw new Error('Failed to delete workflow')
}

export async function enableWorkflow(id: number): Promise<{ status: string; version: number }> {
  const res = await fetch(`${API_BASE}/workflows/${id}/enable`, { method: 'POST' })
  if (!res.ok) throw new Error('Failed to enable workflow')
  return res.json()
}

export async function disableWorkflow(id: number): Promise<{ status: string }> {
  const res = await fetch(`${API_BASE}/workflows/${id}/disable`, { method: 'POST' })
  if (!res.ok) throw new Error('Failed to disable workflow')
  return res.json()
}

export async function checkHealth(): Promise<{ status: string; version: string }> {
  const res = await fetch(`${API_BASE}/health`)
  if (!res.ok) throw new Error('Health check failed')
  return res.json()
}
