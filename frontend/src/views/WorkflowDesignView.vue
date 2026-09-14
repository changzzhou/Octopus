<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, onBeforeRouteLeave } from 'vue-router'
import { getWorkflowDetail, saveWorkflow, type WorkflowDetail, type SaveWorkflowRequest } from '../api/workflows'
import type { Node, Edge, CanvasMeta } from '../types/workflow'
import WorkflowDesigner from '../components/workflow/WorkflowDesigner.vue'

const route = useRoute()

const workflow = ref<WorkflowDetail | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const saveStatus = ref<'idle' | 'saving' | 'saved' | 'error'>('idle')
const saveError = ref<string | null>(null)
const designerRef = ref<InstanceType<typeof WorkflowDesigner> | null>(null)

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
  <div class="workflow-design-page">
    <header class="page-header">
      <div class="header-left">
        <router-link to="/" class="logo">
          🐙 Octopus
        </router-link>
        <span class="header-divider">/</span>
        <router-link to="/workflows" class="breadcrumb">
          Workflows
        </router-link>
        <span class="header-divider">/</span>
        <router-link 
          v-if="workflow" 
          :to="`/workflows/${workflow.id}`" 
          class="breadcrumb"
        >
          {{ workflow.name }}
        </router-link>
        <span v-else class="breadcrumb">Loading...</span>
        <span class="header-divider">/</span>
        <span class="current-page">Design</span>
        <span v-if="workflow" class="status-badge" :class="workflow.status">
          {{ workflow.status }}
        </span>
      </div>
      <nav class="header-right">
        <div v-if="saveStatus === 'saved'" class="save-indicator saved">
          ✓ Saved
        </div>
        <div v-else-if="saveStatus === 'error'" class="save-indicator error" :title="saveError || undefined">
          ⚠ Error saving
        </div>
        <router-link 
          v-if="workflow"
          :to="`/workflows/${workflow.id}`"
          class="back-link"
        >
          ← Back to Details
        </router-link>
      </nav>
    </header>

    <main class="page-content">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Loading workflow...</p>
      </div>

      <div v-else-if="error" class="error-state">
        <div class="error-icon">⚠️</div>
        <p>{{ error }}</p>
        <button class="btn btn-secondary" @click="loadWorkflow">
          Try Again
        </button>
      </div>

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

<style scoped>
.workflow-design-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f1f5f9;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  height: 56px;
  background: white;
  border-bottom: 1px solid #e2e8f0;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo {
  font-size: 20px;
  font-weight: 700;
  color: #7c3aed;
  text-decoration: none;
}

.header-divider {
  color: #cbd5e1;
}

.breadcrumb {
  color: #64748b;
  text-decoration: none;
  font-size: 14px;
}

.breadcrumb:hover {
  color: #7c3aed;
}

.current-page {
  font-size: 14px;
  font-weight: 500;
  color: #1e293b;
}

.status-badge {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  margin-left: 8px;
}

.status-badge.draft {
  background: #fef3c7;
  color: #92400e;
}

.status-badge.enabled {
  background: #d1fae5;
  color: #065f46;
}

.status-badge.disabled {
  background: #fee2e2;
  color: #991b1b;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.save-indicator {
  font-size: 13px;
  padding: 4px 10px;
  border-radius: 12px;
}

.save-indicator.saved {
  background: #d1fae5;
  color: #065f46;
}

.save-indicator.error {
  background: #fee2e2;
  color: #991b1b;
  cursor: help;
}

.back-link {
  color: #64748b;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.2s;
}

.back-link:hover {
  color: #7c3aed;
}

.page-content {
  flex: 1;
  overflow: hidden;
}

.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 16px;
  color: #64748b;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #e2e8f0;
  border-top-color: #7c3aed;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-icon {
  font-size: 48px;
}

.error-state p {
  color: #dc2626;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary {
  background: #e2e8f0;
  color: #475569;
}

.btn-secondary:hover {
  background: #cbd5e1;
}
</style>
