// Workflow status - aligned with BE-1 contract
export type WorkflowStatus = 'draft' | 'enabled' | 'disabled'

// Node types - aligned with BE-1 (NO condition node per Appendix C)
export type NodeType = 'script' | 'http' | 'human'

// Edge outlet - aligned with BE-1
export type EdgeOutlet = 'success' | 'failure'

// Position - aligned with BE-1
export interface Position {
  x: number
  y: number
}

// Node - aligned with BE-1 contract
export interface Node {
  id: string
  type: NodeType
  name: string
  position: Position
  config?: string // JSON string (BE-1 contract)
  credential_ref?: string
}

// Edge - aligned with BE-1 contract
export interface Edge {
  id: string
  source: string
  target: string
  outlet: EdgeOutlet // success|failure (BE-1 contract)
  label?: string
}

// Canvas metadata - aligned with BE-1
export interface CanvasMeta {
  viewport_x?: number
  viewport_y?: number
  zoom?: number
}

// WorkflowDetail - aligned with BE-1 GET /api/v1/workflows/:id response
export interface WorkflowDetail {
  id: number
  name: string
  description?: string
  status: WorkflowStatus // draft|enabled|disabled (string enum, not number)
  version: number
  nodes: Node[]
  edges: Edge[]
  entry_node_id?: string
  variables_schema?: string
  canvas_meta?: CanvasMeta
  created_by?: string
  updated_by?: string
  created_at: string
  updated_at: string
}

// WorkflowSummary - for list view
export interface WorkflowSummary {
  id: number
  name: string
  description?: string
  status: WorkflowStatus
  version: number
  created_at: string
  updated_at: string
}

// SaveWorkflowRequest - aligned with BE-1 PUT /api/v1/workflows/:id
export interface SaveWorkflowRequest {
  name?: string
  description?: string
  nodes?: Node[]
  edges?: Edge[]
  entry_node_id?: string
  variables_schema?: string
  canvas_meta?: CanvasMeta
  version: number // For optimistic concurrency
}

// ========== FE-only types for Vue Flow canvas rendering ==========

// Node config types for UI forms (stored as JSON string in BE)
export interface ScriptNodeConfig {
  language: 'javascript' | 'python' | 'shell'
  code: string
  timeout?: number
}

export interface HttpNodeConfig {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  url: string
  headers?: Record<string, string>
  body?: string
  timeout?: number
}

export interface HumanNodeConfig {
  assignee?: string
  instructions: string
  formFields?: FormField[]
}

export interface FormField {
  name: string
  type: 'text' | 'number' | 'boolean' | 'select'
  label: string
  required?: boolean
  options?: string[]
}

export type NodeConfig = ScriptNodeConfig | HttpNodeConfig | HumanNodeConfig

// Vue Flow node data (internal FE representation)
export interface WorkflowNodeData {
  label: string
  type: NodeType
  config: NodeConfig
}

// Vue Flow edge data (internal FE representation)
export interface WorkflowEdgeData {
  outlet: EdgeOutlet
}

// Node type UI config
export const NODE_TYPE_CONFIGS: Record<NodeType, { label: string; icon: string; color: string }> = {
  script: {
    label: 'Script',
    icon: '📜',
    color: '#10b981',
  },
  http: {
    label: 'HTTP Request',
    icon: '🌐',
    color: '#3b82f6',
  },
  human: {
    label: 'Human Task',
    icon: '👤',
    color: '#f59e0b',
  },
}

export function createDefaultNodeConfig(type: NodeType): NodeConfig {
  switch (type) {
    case 'script':
      return {
        language: 'javascript',
        code: '// Your code here\nreturn { success: true };',
        timeout: 30000,
      }
    case 'http':
      return {
        method: 'GET',
        url: 'https://api.example.com',
        headers: {},
        timeout: 30000,
      }
    case 'human':
      return {
        instructions: 'Please review and approve',
        formFields: [],
      }
  }
}

// Convert FE Vue Flow node to BE Node format
export function toBeNode(vfNode: { id: string; type: string; position: Position; data: WorkflowNodeData }): Node {
  return {
    id: vfNode.id,
    type: vfNode.type as NodeType,
    name: vfNode.data.label,
    position: vfNode.position,
    config: JSON.stringify(vfNode.data.config),
  }
}

// Convert BE Node to FE Vue Flow node format
export function fromBeNode(beNode: Node): { id: string; type: string; position: Position; data: WorkflowNodeData } {
  let config: NodeConfig
  try {
    config = beNode.config ? JSON.parse(beNode.config) : createDefaultNodeConfig(beNode.type as NodeType)
  } catch {
    config = createDefaultNodeConfig(beNode.type as NodeType)
  }
  
  return {
    id: beNode.id,
    type: beNode.type,
    position: beNode.position,
    data: {
      label: beNode.name,
      type: beNode.type as NodeType,
      config,
    },
  }
}

// Convert FE Vue Flow edge to BE Edge format
export function toBeEdge(vfEdge: { id: string; source: string; target: string; sourceHandle?: string | null; data?: WorkflowEdgeData }): Edge {
  const outlet: EdgeOutlet = vfEdge.sourceHandle === 'failure' ? 'failure' : 'success'
  return {
    id: vfEdge.id,
    source: vfEdge.source,
    target: vfEdge.target,
    outlet: vfEdge.data?.outlet || outlet,
  }
}

// Convert BE Edge to FE Vue Flow edge format
export function fromBeEdge(beEdge: Edge): any {
  const isFailure = beEdge.outlet === 'failure'
  return {
    id: beEdge.id,
    source: beEdge.source,
    target: beEdge.target,
    sourceHandle: beEdge.outlet,
    data: { outlet: beEdge.outlet },
    style: {
      stroke: isFailure ? '#ef4444' : '#10b981',
      strokeWidth: 2,
    },
    animated: isFailure,
    label: isFailure ? '✗' : '✓',
    labelStyle: { fill: isFailure ? '#ef4444' : '#10b981', fontWeight: 600 },
    labelBgStyle: { fill: 'white', fillOpacity: 0.9 },
    labelBgPadding: [4, 4] as [number, number],
  }
}
