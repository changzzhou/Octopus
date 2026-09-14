const API_BASE = '/api/v1'

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

export type WorkflowStatus = 'draft' | 'enabled' | 'disabled'

export interface Workflow {
  id: number
  name: string
  description?: string
  status: WorkflowStatus
  created_at: string
  updated_at: string
}

export interface ListWorkflowsResponse {
  total: number
  workflows: Workflow[]
}

const mockWorkflows: Workflow[] = [
  {
    id: 1,
    name: 'Data Pipeline',
    description: 'ETL workflow for daily data processing',
    status: 'enabled',
    created_at: '2026-09-01T10:00:00Z',
    updated_at: '2026-09-14T08:30:00Z',
  },
  {
    id: 2,
    name: 'User Onboarding',
    description: 'Automated user welcome and setup flow',
    status: 'draft',
    created_at: '2026-09-05T14:00:00Z',
    updated_at: '2026-09-13T16:45:00Z',
  },
  {
    id: 3,
    name: 'Report Generation',
    description: 'Weekly sales report automation',
    status: 'disabled',
    created_at: '2026-09-08T09:00:00Z',
    updated_at: '2026-09-12T11:20:00Z',
  },
  {
    id: 4,
    name: 'Notification Service',
    description: 'Multi-channel notification dispatch',
    status: 'enabled',
    created_at: '2026-09-10T11:00:00Z',
    updated_at: '2026-09-14T07:00:00Z',
  },
]

let mockIdCounter = 100

export async function listWorkflows(page = 1, pageSize = 20): Promise<ListWorkflowsResponse> {
  if (isMockEnabled()) {
    await delay(300)
    const start = (page - 1) * pageSize
    const end = start + pageSize
    return {
      total: mockWorkflows.length,
      workflows: mockWorkflows.slice(start, end),
    }
  }
  
  const res = await fetch(`${API_BASE}/workflows?page=${page}&page_size=${pageSize}`)
  if (!res.ok) throw new Error('Failed to fetch workflows')
  return res.json()
}

export async function getWorkflow(id: number): Promise<{ workflow: Workflow }> {
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

export async function createWorkflow(name: string, description?: string): Promise<{ id: number }> {
  if (isMockEnabled()) {
    await delay(300)
    const newWorkflow: Workflow = {
      id: ++mockIdCounter,
      name,
      description: description || '',
      status: 'draft',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }
    mockWorkflows.unshift(newWorkflow)
    return { id: newWorkflow.id }
  }
  
  const res = await fetch(`${API_BASE}/workflows`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, description }),
  })
  if (!res.ok) throw new Error('Failed to create workflow')
  return res.json()
}

export async function updateWorkflow(id: number, data: { name?: string; description?: string; status?: WorkflowStatus }): Promise<void> {
  if (isMockEnabled()) {
    await delay(200)
    const workflow = mockWorkflows.find(w => w.id === id)
    if (!workflow) throw new Error('Workflow not found')
    Object.assign(workflow, data, { updated_at: new Date().toISOString() })
    return
  }
  
  const res = await fetch(`${API_BASE}/workflows/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) throw new Error('Failed to update workflow')
}

export async function deleteWorkflow(id: number): Promise<void> {
  if (isMockEnabled()) {
    await delay(200)
    const index = mockWorkflows.findIndex(w => w.id === id)
    if (index !== -1) mockWorkflows.splice(index, 1)
    return
  }
  
  const res = await fetch(`${API_BASE}/workflows/${id}`, { method: 'DELETE' })
  if (!res.ok) throw new Error('Failed to delete workflow')
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

function delay(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms))
}
