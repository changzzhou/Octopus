export type NodeType = 'script' | 'http' | 'human'

export type EdgeOutlet = 'success' | 'failure'

export interface WorkflowNodeData {
  label: string
  type: NodeType
  config: ScriptNodeConfig | HttpNodeConfig | HumanNodeConfig
}

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

export interface WorkflowNode {
  id: string
  type: NodeType
  position: { x: number; y: number }
  data: WorkflowNodeData
}

export interface WorkflowEdgeData {
  outlet: EdgeOutlet
}

export interface WorkflowEdge {
  id: string
  source: string
  target: string
  sourceHandle?: string | null
  targetHandle?: string | null
  data?: WorkflowEdgeData
  style?: Record<string, any>
  animated?: boolean
  label?: string
  labelStyle?: Record<string, any>
  labelBgStyle?: Record<string, any>
  labelBgPadding?: [number, number]
}

export interface WorkflowCanvas {
  nodes: WorkflowNode[]
  edges: WorkflowEdge[]
}

export interface WorkflowDetail {
  id: number
  name: string
  description?: string
  status: number
  version: number
  canvas?: WorkflowCanvas
  created_at: string
  updated_at: string
}

export interface SaveWorkflowRequest {
  name?: string
  description?: string
  version: number
  canvas: WorkflowCanvas
}

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

export function createDefaultNodeConfig(type: NodeType): WorkflowNodeData['config'] {
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
