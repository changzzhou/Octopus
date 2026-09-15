<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getWorkflowDetail, getWorkflowRuns, triggerRun, type WorkflowDetail, type Run } from '../api/workflows'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { 
  Workflow, 
  ArrowLeft, 
  Play, 
  Pencil, 
  Loader2,
  AlertCircle,
  Clock,
  GitBranch,
  Calendar,
  Hash,
  Inbox,
} from '@lucide/vue'

const route = useRoute()
const router = useRouter()
const workflow = ref<WorkflowDetail | null>(null)
const runs = ref<Run[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const triggering = ref(false)
const triggerError = ref<string | null>(null)

const statusVariantMap: Record<string, 'secondary' | 'success' | 'destructive'> = {
  draft: 'secondary',
  enabled: 'success',
  disabled: 'destructive',
}

const runStatusVariantMap: Record<string, 'secondary' | 'success' | 'destructive' | 'info' | 'warning'> = {
  pending: 'secondary',
  running: 'info',
  succeeded: 'success',
  failed: 'destructive',
  cancelled: 'warning',
}

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

onMounted(loadData)
</script>

<template>
  <div class="min-h-screen bg-zinc-50 dark:bg-zinc-950">
    <!-- Header -->
    <header class="sticky top-0 z-40 border-b border-zinc-200 bg-white/95 backdrop-blur supports-[backdrop-filter]:bg-white/60 dark:border-zinc-800 dark:bg-zinc-950/95 dark:supports-[backdrop-filter]:bg-zinc-950/60">
      <div class="container mx-auto flex h-14 items-center justify-between px-4">
        <router-link to="/" class="flex items-center gap-2 font-semibold text-zinc-900 dark:text-zinc-50">
          <Workflow class="h-5 w-5" />
          <span>Octopus</span>
        </router-link>
        <nav>
          <router-link 
            to="/workflows" 
            class="flex items-center gap-1 text-sm text-zinc-600 transition-colors hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-50"
          >
            <ArrowLeft class="h-4 w-4" />
            Back to Workflows
          </router-link>
        </nav>
      </div>
    </header>

    <main class="container mx-auto px-4 py-8">
      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-16">
        <Loader2 class="h-6 w-6 animate-spin text-zinc-400" />
        <span class="ml-2 text-sm text-zinc-500">Loading workflow...</span>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="flex flex-col items-center justify-center py-16">
        <AlertCircle class="h-10 w-10 text-red-500" />
        <p class="mt-2 text-sm text-red-600">{{ error }}</p>
      </div>

      <!-- Workflow Detail -->
      <div v-else-if="workflow" class="mx-auto max-w-3xl space-y-6">
        <!-- Main Info Card -->
        <Card>
          <CardHeader>
            <div class="flex items-start justify-between">
              <div>
                <CardTitle class="text-xl">{{ workflow.name }}</CardTitle>
                <p class="mt-1 text-sm text-zinc-500 dark:text-zinc-400">
                  {{ workflow.description || 'No description' }}
                </p>
              </div>
              <Badge :variant="statusVariantMap[workflow.status] || 'secondary'">
                {{ workflow.status }}
              </Badge>
            </div>
          </CardHeader>
          <CardContent class="space-y-6">
            <!-- Metadata Grid -->
            <div class="grid grid-cols-2 gap-4 text-sm">
              <div class="flex items-center gap-2 rounded-md bg-zinc-100 p-3 dark:bg-zinc-900">
                <Hash class="h-4 w-4 text-zinc-500" />
                <div>
                  <div class="text-xs text-zinc-500 dark:text-zinc-400">ID</div>
                  <div class="font-mono text-zinc-900 dark:text-zinc-50">{{ workflow.id }}</div>
                </div>
              </div>
              <div class="flex items-center gap-2 rounded-md bg-zinc-100 p-3 dark:bg-zinc-900">
                <GitBranch class="h-4 w-4 text-zinc-500" />
                <div>
                  <div class="text-xs text-zinc-500 dark:text-zinc-400">Version</div>
                  <div class="font-medium text-zinc-900 dark:text-zinc-50">v{{ workflow.version }}</div>
                </div>
              </div>
              <div class="flex items-center gap-2 rounded-md bg-zinc-100 p-3 dark:bg-zinc-900">
                <Calendar class="h-4 w-4 text-zinc-500" />
                <div>
                  <div class="text-xs text-zinc-500 dark:text-zinc-400">Created</div>
                  <div class="text-zinc-900 dark:text-zinc-50">{{ workflow.created_at }}</div>
                </div>
              </div>
              <div class="flex items-center gap-2 rounded-md bg-zinc-100 p-3 dark:bg-zinc-900">
                <Clock class="h-4 w-4 text-zinc-500" />
                <div>
                  <div class="text-xs text-zinc-500 dark:text-zinc-400">Updated</div>
                  <div class="text-zinc-900 dark:text-zinc-50">{{ workflow.updated_at }}</div>
                </div>
              </div>
            </div>

            <Separator />

            <!-- Actions -->
            <div class="flex flex-wrap gap-3">
              <Button as="router-link" :to="`/workflows/${workflow.id}/design`">
                <Pencil class="h-4 w-4" />
                Open Designer
              </Button>
              
              <Button 
                variant="outline"
                @click="handleTriggerRun"
                :disabled="triggering"
              >
                <Loader2 v-if="triggering" class="h-4 w-4 animate-spin" />
                <Play v-else class="h-4 w-4" />
                {{ workflow.status === 'draft' ? 'Trial Run' : 'Run Workflow' }}
              </Button>
            </div>

            <p v-if="triggerError" class="text-sm text-red-500">
              {{ triggerError }}
            </p>

            <!-- Workflow Definition Summary -->
            <div v-if="workflow.nodes && workflow.nodes.length > 0">
              <Separator class="my-4" />
              <h3 class="text-sm font-medium text-zinc-900 dark:text-zinc-50">Workflow Definition</h3>
              <p class="mt-1 text-sm text-zinc-500 dark:text-zinc-400">
                {{ workflow.nodes.length }} node(s), {{ workflow.edges?.length || 0 }} edge(s)
              </p>
            </div>
          </CardContent>
        </Card>

        <!-- Recent Executions Card -->
        <Card>
          <CardHeader>
            <CardTitle class="text-base">Recent Executions</CardTitle>
          </CardHeader>
          <CardContent>
            <!-- Empty State -->
            <div v-if="runs.length === 0" class="flex flex-col items-center py-8">
              <Inbox class="h-10 w-10 text-zinc-300 dark:text-zinc-700" />
              <p class="mt-2 text-sm text-zinc-500 dark:text-zinc-400">No executions yet</p>
              <p class="text-xs text-zinc-400 dark:text-zinc-500">
                Click the button above to {{ workflow.status === 'draft' ? 'trial run' : 'run' }} this workflow
              </p>
            </div>
            
            <!-- Runs List -->
            <div v-else class="space-y-2">
              <button 
                v-for="run in runs.slice(0, 5)" 
                :key="run.id"
                @click="viewRun(run.id)"
                class="flex w-full items-center justify-between rounded-md border border-zinc-200 p-3 text-left transition-colors hover:bg-zinc-50 dark:border-zinc-800 dark:hover:bg-zinc-900"
              >
                <div class="flex items-center gap-3">
                  <Badge :variant="runStatusVariantMap[run.status] || 'secondary'" class="font-normal">
                    {{ run.status }}
                  </Badge>
                  <span class="text-sm font-medium text-zinc-900 dark:text-zinc-50">
                    Run #{{ run.id }}
                  </span>
                  <Badge v-if="run.trigger_type === 'trial'" variant="warning" class="font-normal">
                    Trial
                  </Badge>
                </div>
                <span class="text-xs text-zinc-500 dark:text-zinc-400">
                  {{ formatTime(run.created_at) }}
                </span>
              </button>
              
              <p v-if="runs.length > 5" class="pt-2 text-center text-xs text-zinc-500 dark:text-zinc-400">
                {{ runs.length - 5 }} more execution(s)...
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </main>
  </div>
</template>
