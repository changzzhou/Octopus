<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getWorkflowDetail, getWorkflowRuns, triggerRun, type WorkflowDetail, type Run } from '../api/workflows'

const route = useRoute()
const router = useRouter()
const workflow = ref<WorkflowDetail | null>(null)
const runs = ref<Run[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const triggering = ref(false)
const triggerError = ref<string | null>(null)

async function loadData() {
  try {
    const id = Number(route.params.id)
    const [workflowData, runsData] = await Promise.all([
      getWorkflowDetail(id),
      getWorkflowRuns(id),
    ])
    workflow.value = workflowData.workflow
    runs.value = runsData.runs
  } catch (e) {
    error.value = 'Failed to load workflow'
  } finally {
    loading.value = false
  }
}

async function handleTriggerRun() {
  if (!workflow.value || triggering.value) return
  
  triggering.value = true
  triggerError.value = null
  
  try {
    const result = await triggerRun(workflow.value.id)
    router.push(`/workflows/${workflow.value.id}/executions/${result.run_id}`)
  } catch (e) {
    triggerError.value = 'Failed to trigger run'
    console.error(e)
  } finally {
    triggering.value = false
  }
}

function viewRun(runId: number) {
  if (!workflow.value) return
  router.push(`/workflows/${workflow.value.id}/executions/${runId}`)
}

function formatTime(dateStr: string) {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const runStatusStyles: Record<string, string> = {
  pending: 'bg-slate-100 text-slate-600 dark:bg-slate-700 dark:text-slate-300',
  running: 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300',
  succeeded: 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300',
  failed: 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300',
  cancelled: 'bg-orange-100 text-orange-700 dark:bg-orange-900 dark:text-orange-300',
}

function getRunStatusStyle(status: string) {
  return runStatusStyles[status] || runStatusStyles.pending
}

onMounted(loadData)
</script>

<template>
  <div class="min-h-screen bg-slate-100 dark:bg-slate-900">
    <header class="bg-white dark:bg-slate-800 shadow">
      <div class="container mx-auto px-4 py-4 flex items-center justify-between">
        <router-link to="/" class="text-2xl font-bold text-purple-600 dark:text-purple-400">
          🐙 Octopus
        </router-link>
        <nav class="space-x-4">
          <router-link 
            to="/workflows" 
            class="text-slate-600 dark:text-slate-300 hover:text-purple-600 dark:hover:text-purple-400"
          >
            ← Back to Workflows
          </router-link>
        </nav>
      </div>
    </header>

    <main class="container mx-auto px-4 py-8">
      <div v-if="loading" class="text-center py-12 text-slate-500">
        Loading workflow...
      </div>

      <div v-else-if="error" class="text-center py-12 text-red-500">
        {{ error }}
      </div>

      <div v-else-if="workflow" class="max-w-4xl mx-auto">
        <div class="bg-white dark:bg-slate-800 rounded-xl shadow-lg p-8">
          <div class="flex items-center justify-between mb-4">
            <h1 class="text-3xl font-bold text-slate-900 dark:text-white">
              {{ workflow.name }}
            </h1>
            <span class="px-3 py-1 rounded-full text-sm font-medium capitalize"
              :class="{
                'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200': workflow.status === 'draft',
                'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200': workflow.status === 'enabled',
                'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200': workflow.status === 'disabled',
              }"
            >
              {{ workflow.status }}
            </span>
          </div>
          
          <p class="text-slate-600 dark:text-slate-400 mb-6">
            {{ workflow.description || 'No description' }}
          </p>

          <div class="grid grid-cols-2 gap-4 text-sm mb-8">
            <div class="bg-slate-50 dark:bg-slate-700 rounded-lg p-4">
              <div class="text-slate-500 dark:text-slate-400">ID</div>
              <div class="font-mono text-slate-900 dark:text-white">{{ workflow.id }}</div>
            </div>
            <div class="bg-slate-50 dark:bg-slate-700 rounded-lg p-4">
              <div class="text-slate-500 dark:text-slate-400">Version</div>
              <div class="font-medium text-slate-900 dark:text-white">v{{ workflow.version }}</div>
            </div>
            <div class="bg-slate-50 dark:bg-slate-700 rounded-lg p-4">
              <div class="text-slate-500 dark:text-slate-400">Created</div>
              <div class="text-slate-900 dark:text-white">{{ workflow.created_at }}</div>
            </div>
            <div class="bg-slate-50 dark:bg-slate-700 rounded-lg p-4">
              <div class="text-slate-500 dark:text-slate-400">Updated</div>
              <div class="text-slate-900 dark:text-white">{{ workflow.updated_at }}</div>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex flex-wrap gap-3">
            <router-link
              :to="`/workflows/${workflow.id}/design`"
              class="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white font-medium rounded-lg text-center transition-colors flex items-center gap-2"
            >
              🎨 Open Designer
            </router-link>
            
            <button 
              @click="handleTriggerRun"
              :disabled="triggering"
              class="px-4 py-2 bg-green-600 hover:bg-green-700 disabled:bg-green-400 text-white font-medium rounded-lg transition-colors flex items-center gap-2"
            >
              <span v-if="triggering" class="animate-spin">⏳</span>
              <span v-else>▶️</span>
              {{ workflow.status === 'draft' ? 'Trial Run' : 'Run Workflow' }}
            </button>
          </div>

          <div v-if="triggerError" class="mt-3 text-sm text-red-500">
            {{ triggerError }}
          </div>

          <!-- Workflow Definition Summary -->
          <div v-if="workflow.nodes && workflow.nodes.length > 0" class="mt-8">
            <h2 class="text-lg font-semibold text-slate-900 dark:text-white mb-4">
              Workflow Definition
            </h2>
            <div class="text-sm text-slate-600 dark:text-slate-400">
              <p>{{ workflow.nodes.length }} node(s), {{ workflow.edges?.length || 0 }} edge(s)</p>
            </div>
          </div>

          <!-- Recent Runs -->
          <div class="mt-8">
            <h3 class="text-lg font-semibold text-slate-900 dark:text-white mb-4">
              Recent Executions
            </h3>
            
            <div v-if="runs.length === 0" class="p-6 bg-slate-50 dark:bg-slate-700 rounded-xl text-center">
              <div class="text-4xl mb-2">📭</div>
              <p class="text-slate-500 dark:text-slate-400">No executions yet</p>
              <p class="text-sm text-slate-400 dark:text-slate-500 mt-1">
                Click the button above to {{ workflow.status === 'draft' ? 'trial run' : 'run' }} this workflow
              </p>
            </div>
            
            <div v-else class="space-y-3">
              <div 
                v-for="run in runs.slice(0, 5)" 
                :key="run.id"
                @click="viewRun(run.id)"
                class="p-4 bg-slate-50 dark:bg-slate-700 rounded-lg cursor-pointer hover:bg-slate-100 dark:hover:bg-slate-600 transition-colors"
              >
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-3">
                    <span 
                      :class="[
                        'px-2 py-1 rounded-full text-xs font-medium',
                        getRunStatusStyle(run.status)
                      ]"
                    >
                      {{ run.status }}
                    </span>
                    <span class="font-medium text-slate-900 dark:text-white">
                      Run #{{ run.id }}
                    </span>
                    <span 
                      v-if="run.trigger_type === 'trial'"
                      class="px-2 py-0.5 bg-amber-100 dark:bg-amber-900 text-amber-700 dark:text-amber-300 rounded text-xs"
                    >
                      Trial
                    </span>
                  </div>
                  <div class="text-sm text-slate-500 dark:text-slate-400">
                    {{ formatTime(run.created_at) }}
                  </div>
                </div>
                <div v-if="run.error_message" class="mt-2 text-sm text-red-500 truncate">
                  {{ run.error_message }}
                </div>
              </div>
              
              <div v-if="runs.length > 5" class="text-center text-sm text-slate-500 dark:text-slate-400 py-2">
                {{ runs.length - 5 }} more execution(s)...
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
