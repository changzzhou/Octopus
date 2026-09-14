<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { listWorkflows, createWorkflow, type WorkflowSummary, type WorkflowStatus } from '../api/workflows'

const router = useRouter()
const workflows = ref<WorkflowSummary[]>([])
const total = ref(0)
const loading = ref(true)
const error = ref<string | null>(null)

const showCreateModal = ref(false)
const newWorkflowName = ref('')
const newWorkflowDescription = ref('')
const creating = ref(false)
const createError = ref<string | null>(null)

const statusStyles: Record<WorkflowStatus, string> = {
  draft: 'bg-slate-100 text-slate-600 dark:bg-slate-700 dark:text-slate-300',
  enabled: 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300',
  disabled: 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300',
}

function getStatusStyle(status: WorkflowStatus): string {
  return statusStyles[status] || statusStyles.draft
}

function formatUpdatedTime(dateStr: string) {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

async function loadWorkflows() {
  loading.value = true
  error.value = null
  try {
    const data = await listWorkflows()
    workflows.value = data.workflows
    total.value = data.total
  } catch (e) {
    error.value = 'Failed to load workflows'
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  if (!newWorkflowName.value.trim()) {
    createError.value = 'Name is required'
    return
  }
  
  creating.value = true
  createError.value = null
  try {
    const result = await createWorkflow(newWorkflowName.value.trim(), newWorkflowDescription.value.trim() || undefined)
    showCreateModal.value = false
    newWorkflowName.value = ''
    newWorkflowDescription.value = ''
    router.push(`/workflows/${result.id}`)
  } catch (e) {
    createError.value = 'Failed to create workflow'
  } finally {
    creating.value = false
  }
}

function openCreateModal() {
  showCreateModal.value = true
  createError.value = null
}

function closeCreateModal() {
  showCreateModal.value = false
  newWorkflowName.value = ''
  newWorkflowDescription.value = ''
  createError.value = null
}

onMounted(loadWorkflows)
</script>

<template>
  <div class="min-h-screen bg-slate-100 dark:bg-slate-900">
    <header class="bg-white dark:bg-slate-800 shadow">
      <div class="container mx-auto px-4 py-4 flex items-center justify-between">
        <router-link to="/" class="text-2xl font-bold text-purple-600 dark:text-purple-400">
          🐙 Octopus
        </router-link>
        <nav>
          <router-link 
            to="/workflows" 
            class="text-slate-600 dark:text-slate-300 hover:text-purple-600 dark:hover:text-purple-400"
          >
            Workflows
          </router-link>
        </nav>
      </div>
    </header>

    <main class="container mx-auto px-4 py-8">
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-3xl font-bold text-slate-900 dark:text-white">
          Workflows
        </h1>
        <button 
          class="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white font-medium rounded-lg transition-colors"
          @click="openCreateModal"
        >
          + New Workflow
        </button>
      </div>

      <div v-if="loading" class="text-center py-12 text-slate-500">
        Loading workflows...
      </div>

      <div v-else-if="error" class="text-center py-12 text-red-500">
        {{ error }}
      </div>

      <div v-else-if="workflows.length === 0" class="text-center py-12">
        <div class="text-6xl mb-4">📋</div>
        <h2 class="text-xl font-medium text-slate-700 dark:text-slate-300 mb-2">
          No workflows yet
        </h2>
        <p class="text-slate-500 dark:text-slate-400 mb-4">
          Create your first workflow to get started.
        </p>
        <button 
          class="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white font-medium rounded-lg transition-colors"
          @click="openCreateModal"
        >
          + Create Workflow
        </button>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full bg-white dark:bg-slate-800 rounded-xl shadow overflow-hidden">
          <thead class="bg-slate-50 dark:bg-slate-700">
            <tr>
              <th class="px-6 py-4 text-left text-sm font-semibold text-slate-700 dark:text-slate-200">Name</th>
              <th class="px-6 py-4 text-left text-sm font-semibold text-slate-700 dark:text-slate-200">Status</th>
              <th class="px-6 py-4 text-left text-sm font-semibold text-slate-700 dark:text-slate-200">Updated At</th>
              <th class="px-6 py-4 text-right text-sm font-semibold text-slate-700 dark:text-slate-200">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 dark:divide-slate-700">
            <tr 
              v-for="workflow in workflows" 
              :key="workflow.id"
              class="hover:bg-slate-50 dark:hover:bg-slate-700/50 transition-colors"
            >
              <td class="px-6 py-4">
                <router-link 
                  :to="`/workflows/${workflow.id}`"
                  class="text-slate-900 dark:text-white font-medium hover:text-purple-600 dark:hover:text-purple-400"
                >
                  {{ workflow.name }}
                </router-link>
                <p v-if="workflow.description" class="text-sm text-slate-500 dark:text-slate-400 mt-1">
                  {{ workflow.description }}
                </p>
              </td>
              <td class="px-6 py-4">
                <span 
                  :class="[getStatusStyle(workflow.status), 'px-2 py-1 rounded-full text-xs font-medium']"
                >
                  {{ workflow.status }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-slate-500 dark:text-slate-400">
                {{ formatUpdatedTime(workflow.updated_at) }}
              </td>
              <td class="px-6 py-4 text-right">
                <router-link 
                  :to="`/workflows/${workflow.id}/design`"
                  class="inline-flex items-center px-3 py-1 text-sm bg-purple-100 hover:bg-purple-200 dark:bg-purple-900 dark:hover:bg-purple-800 text-purple-700 dark:text-purple-300 rounded-lg transition-colors"
                >
                  Design
                </router-link>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="total > 0" class="mt-6 text-center text-slate-500">
        Total: {{ total }} workflow(s)
      </div>
    </main>

    <!-- Create Workflow Modal -->
    <div 
      v-if="showCreateModal" 
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
      @click.self="closeCreateModal"
    >
      <div class="bg-white dark:bg-slate-800 rounded-xl shadow-xl p-6 w-full max-w-md mx-4">
        <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4">
          Create New Workflow
        </h2>
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              Name *
            </label>
            <input 
              v-model="newWorkflowName"
              type="text"
              placeholder="Enter workflow name"
              class="w-full px-3 py-2 border border-slate-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-700 text-slate-900 dark:text-white focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              @keyup.enter="handleCreate"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              Description
            </label>
            <textarea 
              v-model="newWorkflowDescription"
              placeholder="Enter description (optional)"
              rows="3"
              class="w-full px-3 py-2 border border-slate-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-700 text-slate-900 dark:text-white focus:ring-2 focus:ring-purple-500 focus:border-transparent resize-none"
            />
          </div>
          
          <div v-if="createError" class="text-sm text-red-500">
            {{ createError }}
          </div>
        </div>
        
        <div class="flex justify-end gap-3 mt-6">
          <button 
            @click="closeCreateModal"
            class="px-4 py-2 text-slate-600 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200 font-medium transition-colors"
          >
            Cancel
          </button>
          <button 
            @click="handleCreate"
            :disabled="creating"
            class="px-4 py-2 bg-purple-600 hover:bg-purple-700 disabled:bg-purple-400 text-white font-medium rounded-lg transition-colors"
          >
            {{ creating ? 'Creating...' : 'Create' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
