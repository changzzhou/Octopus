const API_BASE = '/api/v1'

export interface Workflow {
  id: number
  name: string
  description?: string
  status: number
  created_at: string
  updated_at: string
}

export interface ListWorkflowsResponse {
  total: number
  workflows: Workflow[]
}

export async function listWorkflows(page = 1, pageSize = 20): Promise<ListWorkflowsResponse> {
  const res = await fetch(`${API_BASE}/workflows?page=${page}&page_size=${pageSize}`)
  if (!res.ok) throw new Error('Failed to fetch workflows')
  return res.json()
}

export async function getWorkflow(id: number): Promise<{ workflow: Workflow }> {
  const res = await fetch(`${API_BASE}/workflows/${id}`)
  if (!res.ok) throw new Error('Failed to fetch workflow')
  return res.json()
}

export async function createWorkflow(name: string, description?: string): Promise<{ id: number }> {
  const res = await fetch(`${API_BASE}/workflows`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, description }),
  })
  if (!res.ok) throw new Error('Failed to create workflow')
  return res.json()
}

export async function deleteWorkflow(id: number): Promise<void> {
  const res = await fetch(`${API_BASE}/workflows/${id}`, { method: 'DELETE' })
  if (!res.ok) throw new Error('Failed to delete workflow')
}

export async function checkHealth(): Promise<{ status: string; version: string }> {
  const res = await fetch(`${API_BASE}/health`)
  if (!res.ok) throw new Error('Health check failed')
  return res.json()
}
