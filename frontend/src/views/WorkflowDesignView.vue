<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, onBeforeRouteLeave } from 'vue-router'
import { getWorkflowDetail, saveWorkflow, type WorkflowDetail, type SaveWorkflowRequest } from '../api/workflows'
import type { Node, Edge, CanvasMeta } from '../types/workflow'
import WorkflowDesigner from '../components/workflow/WorkflowDesigner.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { 
  Workflow, 
  ArrowLeft, 
  Check, 
  AlertTriangle,
  Loader2,
} from '@lucide/vue'

const route = useRoute()

const workflow = ref<WorkflowDetail | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const saveStatus = ref<'idle' | 'saving' | 'saved' | 'error'>('idle')
const saveError = ref<string | null>(null)
const designerRef = ref<InstanceType<typeof WorkflowDesigner> | null>(null)

const statusVariantMap: Record<string, 'secondary' | 'success' | 'destructive'> = {
  draft: 'secondary',
  enabled: 'success',
  disabled: 'destructive',
}

onMounted(async () => {
  await loadWorkflow()
})

async function loadWorkflow() {
  loading.value = true
  error.value = null
  try {
    const id = Number(route.params.id)
    const data = await getWorkflowDetail(id)
    workflow.value = data.workflow
  } catch (e) {
    error.value = 'Failed to load workflow'
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function handleSave(data: { nodes: Node[]; edges: Edge[]; canvas_meta?: CanvasMeta }) {
  if (!workflow.value) return
  
  saveStatus.value = 'saving'
  saveError.value = null
  
  try {
    const request: SaveWorkflowRequest = {
      version: workflow.value.version,
      nodes: data.nodes,
      edges: data.edges,
      canvas_meta: data.canvas_meta,
    }
    
    const result = await saveWorkflow(workflow.value.id, request)
    
    workflow.value.version = result.version || workflow.value.version + 1
    workflow.value.nodes = data.nodes
    workflow.value.edges = data.edges
    workflow.value.canvas_meta = data.canvas_meta
    saveStatus.value = 'saved'
    designerRef.value?.setClean()
    
    setTimeout(() => {
      if (saveStatus.value === 'saved') {
        saveStatus.value = 'idle'
      }
    }, 2000)
  } catch (e) {
    saveStatus.value = 'error'
    saveError.value = e instanceof Error ? e.message : 'Failed to save workflow'
    console.error(e)
  }
}

function handleCanvasChange() {
  saveStatus.value = 'idle'
}

onBeforeRouteLeave((_to, _from, next) => {
  if (designerRef.value?.isDirty()) {
    const answer = window.confirm('You have unsaved changes. Are you sure you want to leave?')
    if (!answer) {
      next(false)
      return
    }
  }
  next()
})

function handleBeforeUnload(e: BeforeUnloadEvent) {
  if (designerRef.value?.isDirty()) {
    e.preventDefault()
    e.returnValue = ''
  }
}

onMounted(() => {
  window.addEventListener('beforeunload', handleBeforeUnload)
})

onBeforeUnmount(() => {
  window.removeEventListener('beforeunload', handleBeforeUnload)
})
</script>

<template>
  <div class="flex h-screen flex-col bg-zinc-100 dark:bg-zinc-900">
    <!-- Header -->
    <header class="flex h-14 flex-shrink-0 items-center justify-between border-b border-zinc-200 bg-white px-4 dark:border-zinc-800 dark:bg-zinc-950">
      <div class="flex items-center gap-2">
        <router-link to="/" class="flex items-center gap-2 font-semibold text-zinc-900 dark:text-zinc-50">
          <Workflow class="h-5 w-5" />
          <span>Octopus</span>
        </router-link>
        
        <Separator orientation="vertical" class="mx-2 h-5" />
        
        <router-link 
          to="/workflows" 
          class="text-sm text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-50"
        >
          Workflows
        </router-link>
        
        <span class="text-zinc-300 dark:text-zinc-700">/</span>
        
        <router-link 
          v-if="workflow" 
          :to="`/workflows/${workflow.id}`" 
          class="text-sm text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-50"
        >
          {{ workflow.name }}
        </router-link>
        <span v-else class="text-sm text-zinc-500">Loading...</span>
        
        <span class="text-zinc-300 dark:text-zinc-700">/</span>
        
        <span class="text-sm font-medium text-zinc-900 dark:text-zinc-50">Design</span>
        
        <Badge 
          v-if="workflow" 
          :variant="statusVariantMap[workflow.status] || 'secondary'"
          class="ml-2"
        >
          {{ workflow.status }}
        </Badge>
      </div>
      
      <div class="flex items-center gap-3">
        <!-- Save Status Indicator -->
        <div 
          v-if="saveStatus === 'saved'" 
          class="flex items-center gap-1 text-sm text-emerald-600 dark:text-emerald-400"
        >
          <Check class="h-4 w-4" />
          <span>Saved</span>
        </div>
        <div 
          v-else-if="saveStatus === 'error'" 
          class="flex items-center gap-1 text-sm text-red-600 dark:text-red-400"
          :title="saveError || undefined"
        >
          <AlertTriangle class="h-4 w-4" />
          <span>Error saving</span>
        </div>
        
        <Button
          v-if="workflow"
          variant="outline"
          size="sm"
          as="router-link"
          :to="`/workflows/${workflow.id}`"
        >
          <ArrowLeft class="h-4 w-4" />
          Back to Details
        </Button>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex flex-1 flex-col overflow-hidden">
      <!-- Loading State -->
      <div v-if="loading" class="flex flex-1 items-center justify-center">
        <Loader2 class="h-8 w-8 animate-spin text-zinc-400" />
        <span class="ml-3 text-zinc-500">Loading workflow...</span>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="flex flex-1 flex-col items-center justify-center gap-4">
        <AlertTriangle class="h-12 w-12 text-red-500" />
        <p class="text-red-600">{{ error }}</p>
        <Button variant="outline" @click="loadWorkflow">
          Try Again
        </Button>
      </div>

      <!-- Designer -->
      <WorkflowDesigner
        v-else-if="workflow"
        ref="designerRef"
        :workflow-id="workflow.id"
        :workflow-name="workflow.name"
        :workflow-version="workflow.version"
        :initial-nodes="workflow.nodes"
        :initial-edges="workflow.edges"
        :initial-canvas-meta="workflow.canvas_meta"
        @save="handleSave"
        @change="handleCanvasChange"
      />
    </main>
  </div>
</template>
