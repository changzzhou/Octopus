<script setup lang="ts">
import { ref, markRaw, computed, type ComponentPublicInstance } from 'vue'
import { VueFlow, type NodeMouseEvent, type VueFlowStore } from '@vue-flow/core'
import { Controls } from '@vue-flow/controls'
import { MiniMap } from '@vue-flow/minimap'
import { Background, BackgroundVariant } from '@vue-flow/background'

import ScriptNode from './ScriptNode.vue'
import HttpNode from './HttpNode.vue'
import HumanNode from './HumanNode.vue'
import NodePalette from './NodePalette.vue'
import ConfigSidebar from './ConfigSidebar.vue'

import type { 
  Node,
  Edge,
  NodeType,
  WorkflowNodeData,
  EdgeOutlet,
  CanvasMeta,
} from '../../types/workflow'
import { 
  createDefaultNodeConfig, 
  NODE_TYPE_CONFIGS,
  toBeNode,
  toBeEdge,
  fromBeNode,
  fromBeEdge,
} from '../../types/workflow'

const props = defineProps<{
  initialNodes?: Node[]
  initialEdges?: Edge[]
  initialCanvasMeta?: CanvasMeta
  workflowId: number
  workflowName: string
  workflowVersion: number
}>()

const emit = defineEmits<{
  (e: 'save', data: { nodes: Node[]; edges: Edge[]; canvas_meta?: CanvasMeta }): void
  (e: 'change'): void
}>()

const nodeTypes: Record<string, any> = {
  script: markRaw(ScriptNode),
  http: markRaw(HttpNode),
  human: markRaw(HumanNode),
}

const selectedNode = ref<any>(null)
const isDirty = ref(false)
const saving = ref(false)
const isFlowReady = ref(false)

const flowId = `workflow-${props.workflowId}`

// Store the VueFlow instance received from @init event
const vfInstance = ref<VueFlowStore | null>(null)

// Ref to VueFlow component for getBoundingClientRect
const vueFlowRef = ref<ComponentPublicInstance | null>(null)

// Compute initial nodes/edges from props  
const initialVfNodes = computed(() => 
  props.initialNodes?.length 
    ? props.initialNodes.map(fromBeNode) 
    : []
)

const initialVfEdges = computed(() =>
  props.initialEdges?.length 
    ? props.initialEdges.map(fromBeEdge) 
    : []
)

const nodeCount = computed(() => {
  if (!vfInstance.value) return 0
  // getNodes is a computed ref in VueFlow store
  const nodesRef = vfInstance.value.getNodes as any
  const nodes = nodesRef?.value ?? nodesRef
  return Array.isArray(nodes) ? nodes.length : 0
})

// Handle VueFlow init event - receive the store instance
function handleInit(instance: VueFlowStore) {
  vfInstance.value = instance
  isFlowReady.value = true
}

function markDirty() {
  isDirty.value = true
  emit('change')
}

function handleConnect(connection: any) {
  if (!vfInstance.value || !isFlowReady.value) return
  
  const outlet: EdgeOutlet = connection.sourceHandle === 'failure' ? 'failure' : 'success'
  
  const newEdge: any = {
    id: `e-${connection.source}-${connection.sourceHandle}-${connection.target}-${Date.now()}`,
    source: connection.source,
    target: connection.target,
    sourceHandle: connection.sourceHandle,
    targetHandle: connection.targetHandle,
    data: { outlet },
    style: {
      stroke: outlet === 'success' ? '#10b981' : '#ef4444',
      strokeWidth: 2,
    },
    animated: outlet === 'failure',
    label: outlet === 'success' ? '✓' : '✗',
    labelStyle: { fill: outlet === 'success' ? '#10b981' : '#ef4444', fontWeight: 600 },
    labelBgStyle: { fill: 'white', fillOpacity: 0.9 },
    labelBgPadding: [4, 4] as [number, number],
  }
  
  vfInstance.value.addEdges([newEdge])
  markDirty()
}

function onDragOver(event: DragEvent) {
  event.preventDefault()
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'move'
  }
}

function onDrop(event: DragEvent) {
  event.preventDefault()
  
  if (!vfInstance.value || !isFlowReady.value) return
  
  const type = event.dataTransfer?.getData('application/vueflow-nodetype') as NodeType
  if (!type) return
  
  // Get VueFlow element bounds
  const el = vueFlowRef.value?.$el as HTMLElement | undefined
  if (!el) return
  
  const bounds = el.getBoundingClientRect()
  
  // Convert screen coordinates to flow coordinates
  const position = vfInstance.value.screenToFlowCoordinate({
    x: event.clientX - bounds.left,
    y: event.clientY - bounds.top,
  })
  
  const config = NODE_TYPE_CONFIGS[type]
  const nodeId = `${type}-${Date.now()}`
  
  const newNode: any = {
    id: nodeId,
    type,
    position,
    data: {
      label: config.label,
      type,
      config: createDefaultNodeConfig(type),
    },
  }
  
  vfInstance.value.addNodes([newNode])
  markDirty()
}

function onNodeClick(event: NodeMouseEvent) {
  selectedNode.value = event.node
}

function onPaneClick() {
  selectedNode.value = null
}

function onNodeDragStop() {
  markDirty()
}

function updateNodeData(nodeId: string, data: WorkflowNodeData) {
  if (!vfInstance.value) return
  
  const nodesRef = vfInstance.value.getNodes as any
  const nodes = nodesRef?.value ?? nodesRef
  if (!Array.isArray(nodes)) return
  
  const nodeIndex = nodes.findIndex((n: any) => n.id === nodeId)
  if (nodeIndex !== -1) {
    const updatedNodes = nodes.map((n: any) => 
      n.id === nodeId ? { ...n, data } : n
    )
    vfInstance.value.setNodes(updatedNodes)
    
    if (selectedNode.value?.id === nodeId) {
      selectedNode.value = { ...selectedNode.value, data }
    }
    markDirty()
  }
}

function deleteNode(nodeId: string) {
  if (!vfInstance.value) return
  vfInstance.value.removeNodes([nodeId])
  selectedNode.value = null
  markDirty()
}

function closeSidebar() {
  selectedNode.value = null
}

async function saveWorkflow() {
  if (!vfInstance.value) return
  
  saving.value = true
  try {
    const nodesRef = vfInstance.value.getNodes as any
    const edgesRef = vfInstance.value.getEdges as any
    const vfNodes = nodesRef?.value ?? nodesRef
    const vfEdges = edgesRef?.value ?? edgesRef
    
    const beNodes: Node[] = (vfNodes || []).map((n: any) => toBeNode({
      id: n.id,
      type: n.type || 'script',
      position: n.position,
      data: n.data as WorkflowNodeData,
    }))
    
    const beEdges: Edge[] = (vfEdges || []).map((e: any) => toBeEdge({
      id: e.id,
      source: e.source,
      target: e.target,
      sourceHandle: e.sourceHandle,
      data: e.data,
    }))
    
    const viewport = vfInstance.value.getViewport()
    const canvasMeta: CanvasMeta = {
      viewport_x: viewport.x,
      viewport_y: viewport.y,
      zoom: viewport.zoom,
    }
    
    emit('save', { nodes: beNodes, edges: beEdges, canvas_meta: canvasMeta })
  } finally {
    saving.value = false
  }
}

defineExpose({
  getCanvas: () => {
    if (!vfInstance.value) return { nodes: [], edges: [] }
    
    const nodesRef = vfInstance.value.getNodes as any
    const edgesRef = vfInstance.value.getEdges as any
    const vfNodes = nodesRef?.value ?? nodesRef
    const vfEdges = edgesRef?.value ?? edgesRef
    
    return {
      nodes: (vfNodes || []).map((n: any) => toBeNode({
        id: n.id,
        type: n.type || 'script',
        position: n.position,
        data: n.data as WorkflowNodeData,
      })),
      edges: (vfEdges || []).map((e: any) => toBeEdge({
        id: e.id,
        source: e.source,
        target: e.target,
        sourceHandle: e.sourceHandle,
        data: e.data,
      })),
    }
  },
  isDirty: () => isDirty.value,
  setClean: () => { isDirty.value = false },
})
</script>

<template>
  <div class="workflow-designer">
    <div class="designer-toolbar">
      <div class="toolbar-left">
        <h2 class="workflow-title">{{ workflowName }}</h2>
        <span class="version-badge">v{{ workflowVersion }}</span>
        <span v-if="isDirty" class="unsaved-badge">Unsaved</span>
      </div>
      <div class="toolbar-right">
        <button 
          class="btn btn-primary"
          :disabled="!isDirty || saving"
          @click="saveWorkflow"
        >
          <span v-if="saving">Saving...</span>
          <span v-else>💾 Save</span>
        </button>
      </div>
    </div>
    
    <div class="designer-content">
      <NodePalette class="designer-palette" />
      
      <!-- Drop zone wrapper -->
      <div 
        class="designer-canvas"
        @drop.capture="onDrop"
        @dragover.capture="onDragOver"
      >
        <!-- VueFlow provides nodes/edges -->
        <VueFlow
          ref="vueFlowRef"
          :id="flowId"
          :nodes="initialVfNodes"
          :edges="initialVfEdges"
          :node-types="nodeTypes"
          :default-viewport="{ zoom: 1, x: 100, y: 100 }"
          :min-zoom="0.25"
          :max-zoom="2"
          @drop="onDrop"
          @dragover="onDragOver"
          @node-click="onNodeClick"
          @pane-click="onPaneClick"
          @node-drag-stop="onNodeDragStop"
          @connect="handleConnect"
          @init="handleInit"
        >
          <Background :variant="BackgroundVariant.Dots" :gap="20" :size="1" />
          <Controls position="bottom-left" />
          <MiniMap position="bottom-right" />
        </VueFlow>
        
        <div v-if="nodeCount === 0 && isFlowReady" class="empty-canvas">
          <div class="empty-icon">🎨</div>
          <p class="empty-title">Start Building Your Workflow</p>
          <p class="empty-hint">Drag nodes from the palette to get started</p>
        </div>
        
        <div v-if="!isFlowReady" class="loading-canvas">
          <div class="spinner"></div>
          <p>Initializing canvas...</p>
        </div>
      </div>
      
      <ConfigSidebar
        class="designer-sidebar"
        :node="selectedNode"
        @update="updateNodeData"
        @delete="deleteNode"
        @close="closeSidebar"
      />
    </div>
  </div>
</template>

<style>
@import '@vue-flow/core/dist/style.css';
@import '@vue-flow/core/dist/theme-default.css';
@import '@vue-flow/controls/dist/style.css';
@import '@vue-flow/minimap/dist/style.css';
</style>

<style scoped>
.workflow-designer {
  display: flex;
  flex-direction: column;
  /* Use flex: 1 to properly fill parent flex container */
  flex: 1;
  min-height: 0; /* Allow flex shrinking */
  background: #f1f5f9;
}

.designer-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: white;
  border-bottom: 1px solid #e2e8f0;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.workflow-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.version-badge {
  background: #e0e7ff;
  color: #3730a3;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.unsaved-badge {
  background: #fef3c7;
  color: #92400e;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #7c3aed;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #6d28d9;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.designer-content {
  flex: 1;
  display: flex;
  gap: 16px;
  padding: 16px;
  overflow: hidden;
  min-height: 0; /* Allow flex shrinking */
}

.designer-palette {
  flex-shrink: 0;
}

.designer-canvas {
  flex: 1;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  position: relative;
  overflow: hidden;
  min-height: 0; /* Allow flex shrinking */
  min-width: 0; /* Allow flex shrinking */
}

/* Ensure VueFlow fills its container */
.designer-canvas :deep(.vue-flow) {
  position: absolute !important;
  inset: 0 !important;
}

.designer-sidebar {
  flex-shrink: 0;
}

.empty-canvas,
.loading-canvas {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  z-index: 1;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px;
}

.empty-hint {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
}

.loading-canvas {
  background: rgba(255, 255, 255, 0.9);
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #e2e8f0;
  border-top-color: #7c3aed;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 12px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

:deep(.vue-flow__minimap) {
  border-radius: 8px;
  overflow: hidden;
}

:deep(.vue-flow__controls) {
  border-radius: 8px;
  overflow: hidden;
}

:deep(.vue-flow__edge-path) {
  stroke-width: 2;
}
</style>
