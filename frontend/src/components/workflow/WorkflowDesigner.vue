<script setup lang="ts">
import { ref, markRaw, type ComponentPublicInstance } from 'vue'
import { VueFlow, type NodeMouseEvent, type VueFlowStore } from '@vue-flow/core'
import { Controls } from '@vue-flow/controls'
import { MiniMap } from '@vue-flow/minimap'
import { Background, BackgroundVariant } from '@vue-flow/background'

import ScriptNode from './ScriptNode.vue'
import HttpNode from './HttpNode.vue'
import HumanNode from './HumanNode.vue'
import NodePalette from './NodePalette.vue'
import ConfigSidebar from './ConfigSidebar.vue'

import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Save, Loader2, Paintbrush } from '@lucide/vue'

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

const vfInstance = ref<VueFlowStore | null>(null)
const vueFlowRef = ref<ComponentPublicInstance | null>(null)
const nodeCount = ref(0)

function updateNodeCount() {
  if (!vfInstance.value) {
    nodeCount.value = 0
    return
  }
  const nodesRef = vfInstance.value.getNodes as any
  const nodes = nodesRef?.value ?? nodesRef
  nodeCount.value = Array.isArray(nodes) ? nodes.length : 0
}

function handleInit(instance: VueFlowStore) {
  vfInstance.value = instance
  
  if (props.initialNodes?.length) {
    const vfNodes = props.initialNodes.map(fromBeNode)
    instance.setNodes(vfNodes)
  }
  if (props.initialEdges?.length) {
    const vfEdges = props.initialEdges.map(fromBeEdge)
    instance.setEdges(vfEdges)
  }
  
  updateNodeCount()
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
  event.stopPropagation()
  event.preventDefault()
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'move'
  }
}

function onDrop(event: DragEvent) {
  event.preventDefault()
  event.stopPropagation()
  
  if (!vfInstance.value || !isFlowReady.value) return
  
  const type = event.dataTransfer?.getData('application/vueflow-nodetype') as NodeType
  if (!type) return
  
  const el = vueFlowRef.value?.$el as HTMLElement | undefined
  if (!el) return
  
  const bounds = el.getBoundingClientRect()
  
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
  updateNodeCount()
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
  updateNodeCount()
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
  <div class="flex flex-1 flex-col overflow-hidden">
    <!-- Toolbar -->
    <div class="flex h-12 flex-shrink-0 items-center justify-between border-b border-zinc-200 bg-white px-4 dark:border-zinc-800 dark:bg-zinc-950">
      <div class="flex items-center gap-3">
        <h2 class="text-sm font-medium text-zinc-900 dark:text-zinc-50">{{ workflowName }}</h2>
        <Badge variant="outline" class="font-mono text-xs">v{{ workflowVersion }}</Badge>
        <Badge v-if="isDirty" variant="warning" class="text-xs">Unsaved</Badge>
      </div>
      <div class="flex items-center gap-2">
        <Button 
          size="sm"
          :disabled="!isDirty || saving"
          @click="saveWorkflow"
        >
          <Loader2 v-if="saving" class="h-4 w-4 animate-spin" />
          <Save v-else class="h-4 w-4" />
          {{ saving ? 'Saving...' : 'Save' }}
        </Button>
      </div>
    </div>
    
    <!-- Designer Content -->
    <div class="flex flex-1 gap-4 overflow-hidden p-4">
      <NodePalette class="flex-shrink-0" />
      
      <!-- Canvas -->
      <div 
        class="relative flex-1 overflow-hidden rounded-lg border border-zinc-200 bg-white dark:border-zinc-800 dark:bg-zinc-950"
        @drop.capture="onDrop"
        @dragover.capture="onDragOver"
      >
        <VueFlow
          ref="vueFlowRef"
          :id="flowId"
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
        
        <!-- Empty State -->
        <div v-if="nodeCount === 0 && isFlowReady" class="pointer-events-none absolute inset-0 flex flex-col items-center justify-center">
          <Paintbrush class="h-12 w-12 text-zinc-300 dark:text-zinc-700" />
          <p class="mt-4 text-sm font-medium text-zinc-700 dark:text-zinc-300">Start Building Your Workflow</p>
          <p class="mt-1 text-xs text-zinc-500 dark:text-zinc-500">Drag nodes from the palette to get started</p>
        </div>
        
        <!-- Loading State -->
        <div v-if="!isFlowReady" class="absolute inset-0 flex flex-col items-center justify-center bg-white/90 dark:bg-zinc-950/90">
          <Loader2 class="h-8 w-8 animate-spin text-zinc-400" />
          <p class="mt-3 text-sm text-zinc-500">Initializing canvas...</p>
        </div>
      </div>
      
      <ConfigSidebar
        class="flex-shrink-0"
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
:deep(.vue-flow__minimap) {
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e4e4e7;
}

:deep(.vue-flow__controls) {
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e4e4e7;
}

:deep(.vue-flow__edge-path) {
  stroke-width: 2;
}
</style>
