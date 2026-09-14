<script setup lang="ts">
import { ref, onMounted, markRaw, nextTick, computed } from 'vue'
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

const flowId = `workflow-${props.workflowId}`

const nodes = ref<any[]>([])
const edges = ref<any[]>([])

const nodeCount = computed(() => nodes.value.length)

const vfInstance = ref<VueFlowStore | null>(null)

function onInit(instance: VueFlowStore) {
  vfInstance.value = instance
}

onMounted(async () => {
  await nextTick()
  
  if (props.initialNodes?.length) {
    nodes.value = props.initialNodes.map(fromBeNode)
  }
  
  if (props.initialEdges?.length) {
    edges.value = props.initialEdges.map(fromBeEdge)
  }
})

function markDirty() {
  isDirty.value = true
  emit('change')
}

function handleConnect(connection: any) {
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
  
  edges.value = [...edges.value, newEdge]
  markDirty()
}

function onDrop(event: DragEvent) {
  const type = event.dataTransfer?.getData('application/vueflow-nodetype') as NodeType
  if (!type) return
  
  event.preventDefault()
  
  if (!vfInstance.value) return
  
  const position = vfInstance.value.screenToFlowCoordinate({
    x: event.clientX,
    y: event.clientY,
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
  nextTick(() => {
    markDirty()
  })
}

function onDragOver(event: DragEvent) {
  event.preventDefault()
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'move'
  }
}

function onNodeClick(event: NodeMouseEvent) {
  selectedNode.value = event.node
}

function onPaneClick() {
  selectedNode.value = null
}

function updateNodeData(nodeId: string, data: WorkflowNodeData) {
  const nodeIndex = nodes.value.findIndex(n => n.id === nodeId)
  if (nodeIndex !== -1) {
    const updatedNodes = [...nodes.value]
    updatedNodes[nodeIndex] = {
      ...updatedNodes[nodeIndex],
      data,
    }
    nodes.value = updatedNodes
    
    if (selectedNode.value?.id === nodeId) {
      selectedNode.value = { ...selectedNode.value, data }
    }
    markDirty()
  }
}

function deleteNode(nodeId: string) {
  nodes.value = nodes.value.filter(n => n.id !== nodeId)
  edges.value = edges.value.filter(e => e.source !== nodeId && e.target !== nodeId)
  selectedNode.value = null
  markDirty()
}

function closeSidebar() {
  selectedNode.value = null
}

async function saveWorkflow() {
  saving.value = true
  try {
    const vfNodes = nodes.value
    const vfEdges = edges.value
    
    const beNodes: Node[] = vfNodes.map(n => toBeNode({
      id: n.id,
      type: n.type || 'script',
      position: n.position,
      data: n.data as WorkflowNodeData,
    }))
    
    const beEdges: Edge[] = vfEdges.map(e => toBeEdge({
      id: e.id,
      source: e.source,
      target: e.target,
      sourceHandle: e.sourceHandle,
      data: e.data,
    }))
    
    let canvasMeta: CanvasMeta | undefined
    if (vfInstance.value) {
      const viewport = vfInstance.value.getViewport()
      canvasMeta = {
        viewport_x: viewport.x,
        viewport_y: viewport.y,
        zoom: viewport.zoom,
      }
    }
    
    emit('save', { nodes: beNodes, edges: beEdges, canvas_meta: canvasMeta })
  } finally {
    saving.value = false
  }
}

defineExpose({
  getCanvas: () => {
    const vfNodes = nodes.value
    const vfEdges = edges.value
    return {
      nodes: vfNodes.map(n => toBeNode({
        id: n.id,
        type: n.type || 'script',
        position: n.position,
        data: n.data as WorkflowNodeData,
      })),
      edges: vfEdges.map(e => toBeEdge({
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
      
      <div 
        class="designer-canvas"
        @drop="onDrop"
        @dragover="onDragOver"
      >
        <VueFlow
          :id="flowId"
          v-model:nodes="nodes"
          v-model:edges="edges"
          :node-types="nodeTypes"
          :default-viewport="{ zoom: 1, x: 100, y: 100 }"
          :min-zoom="0.25"
          :max-zoom="2"
          @node-click="onNodeClick"
          @pane-click="onPaneClick"
          @node-drag-stop="markDirty"
          @connect="handleConnect"
          @vue-flow-init="onInit"
        >
          <Background :variant="BackgroundVariant.Dots" :gap="20" :size="1" />
          <Controls position="bottom-left" />
          <MiniMap position="bottom-right" />
        </VueFlow>
        
        <div v-if="nodeCount === 0" class="empty-canvas">
          <div class="empty-icon">🎨</div>
          <p class="empty-title">Start Building Your Workflow</p>
          <p class="empty-hint">Drag nodes from the palette to get started</p>
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
  height: 100%;
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
}

.designer-sidebar {
  flex-shrink: 0;
}

.empty-canvas {
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
