<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  getRun, 
  getRunSteps, 
  pollRunStatus,
  type Run, 
  type Step, 
  type RunStatus,
  type StepStatus,
} from '../api/workflows'
import type { Node } from '../types/workflow'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { 
  Workflow, 
  ArrowLeft, 
  Loader2,
  AlertTriangle,
  CheckCircle2,
  XCircle,
  Clock,
  Ban,
  Circle,
  ChevronRight,
  RefreshCw,
  GitBranch,
  Activity,
  Zap,
} from '@lucide/vue'

const route = useRoute()
const router = useRouter()

const workflowId = computed(() => Number(route.params.id))
const runId = computed(() => Number(route.params.executionId))

const run = ref<Run | null>(null)
const steps = ref<Step[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const selectedStepId = ref<string | null>(null)

let stopPolling: (() => void) | null = null

interface StatusConfig {
  label: string
  variant: 'secondary' | 'info' | 'success' | 'destructive' | 'warning'
  icon: any
}

const runStatusConfig: Record<RunStatus, StatusConfig> = {
  pending: { label: 'Pending', variant: 'secondary', icon: Clock },
  running: { label: 'Running', variant: 'info', icon: Loader2 },
  succeeded: { label: 'Succeeded', variant: 'success', icon: CheckCircle2 },
  failed: { label: 'Failed', variant: 'destructive', icon: XCircle },
  cancelled: { label: 'Cancelled', variant: 'warning', icon: Ban },
}

const stepStatusConfig: Record<StepStatus, StatusConfig & { borderClass: string }> = {
  pending: { label: 'Pending', variant: 'secondary', icon: Circle, borderClass: 'border-zinc-300 dark:border-zinc-600' },
  running: { label: 'Running', variant: 'info', icon: Loader2, borderClass: 'border-blue-400 dark:border-blue-500' },
  succeeded: { label: 'Succeeded', variant: 'success', icon: CheckCircle2, borderClass: 'border-emerald-400 dark:border-emerald-500' },
  failed: { label: 'Failed', variant: 'destructive', icon: XCircle, borderClass: 'border-red-400 dark:border-red-500' },
  skipped: { label: 'Skipped', variant: 'secondary', icon: Circle, borderClass: 'border-zinc-300 dark:border-zinc-600' },
}

const orderedSteps = computed(() => {
  if (!run.value?.definition_snapshot?.nodes) return steps.value
  
  const nodeOrder = new Map<string, number>()
  const nodes = run.value.definition_snapshot.nodes
  const edges = run.value.definition_snapshot.edges || []
  
  const visited = new Set<string>()
  const order: string[] = []
  
  function visit(nodeId: string) {
    if (visited.has(nodeId)) return
    visited.add(nodeId)
    order.push(nodeId)
    
    const outEdges = edges.filter(e => e.source === nodeId)
    for (const edge of outEdges) {
      visit(edge.target)
    }
  }
  
  const entryId = run.value.definition_snapshot.entry_node_id
  if (entryId) {
    visit(entryId)
  } else if (nodes.length > 0) {
    visit(nodes[0].id)
  }
  
  for (const node of nodes) {
    if (!visited.has(node.id)) {
      visit(node.id)
    }
  }
  
  order.forEach((id, idx) => nodeOrder.set(id, idx))
  
  return [...steps.value].sort((a, b) => {
    const orderA = nodeOrder.get(a.node_id) ?? 999
    const orderB = nodeOrder.get(b.node_id) ?? 999
    return orderA - orderB
  })
})

const selectedStep = computed(() => {
  if (!selectedStepId.value) return null
  return steps.value.find(s => s.node_id === selectedStepId.value) || null
})

const selectedNode = computed((): Node | null => {
  if (!selectedStepId.value || !run.value?.definition_snapshot?.nodes) return null
  return run.value.definition_snapshot.nodes.find(n => n.id === selectedStepId.value) || null
})

function getStepConfig(status: StepStatus) {
  return stepStatusConfig[status] || stepStatusConfig.pending
}

function getRunConfig(status: RunStatus) {
  return runStatusConfig[status] || runStatusConfig.pending
}

function formatTime(dateStr?: string): string {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

function formatDuration(startStr?: string, endStr?: string): string {
  if (!startStr) return '-'
  const start = new Date(startStr)
  const end = endStr ? new Date(endStr) : new Date()
  const diffMs = end.getTime() - start.getTime()
  
  if (diffMs < 1000) return `${diffMs}ms`
  if (diffMs < 60000) return `${(diffMs / 1000).toFixed(1)}s`
  const mins = Math.floor(diffMs / 60000)
  const secs = Math.floor((diffMs % 60000) / 1000)
  return `${mins}m ${secs}s`
}

function selectStep(nodeId: string) {
  selectedStepId.value = selectedStepId.value === nodeId ? null : nodeId
}

async function loadInitialData() {
  loading.value = true
  error.value = null
  
  try {
    const [runRes, stepsRes] = await Promise.all([
      getRun(runId.value),
      getRunSteps(runId.value),
    ])
    
    run.value = runRes.run
    steps.value = stepsRes.steps
    
    const isTerminal = ['succeeded', 'failed', 'cancelled'].includes(run.value.status)
    if (!isTerminal) {
      startPolling()
    }
  } catch (e) {
    error.value = 'Failed to load execution data'
    console.error(e)
  } finally {
    loading.value = false
  }
}

function startPolling() {
  if (stopPolling) stopPolling()
  
  stopPolling = pollRunStatus(runId.value, (updatedRun, updatedSteps) => {
    run.value = updatedRun
    steps.value = updatedSteps
  })
}

function goBack() {
  router.push(`/workflows/${workflowId.value}`)
}

watch(runId, () => {
  if (stopPolling) {
    stopPolling()
    stopPolling = null
  }
  loadInitialData()
})

onMounted(loadInitialData)

onUnmounted(() => {
  if (stopPolling) {
    stopPolling()
    stopPolling = null
  }
})
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
          <Button variant="ghost" size="sm" @click="goBack">
            <ArrowLeft class="h-4 w-4" />
            Back to Workflow
          </Button>
        </nav>
      </div>
    </header>

    <main class="container mx-auto px-4 py-6">
      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-20">
        <Loader2 class="h-8 w-8 animate-spin text-zinc-400" />
        <span class="ml-3 text-zinc-500">Loading execution...</span>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="flex flex-col items-center justify-center py-20">
        <AlertTriangle class="h-12 w-12 text-red-500" />
        <p class="mt-3 text-red-600">{{ error }}</p>
        <Button variant="outline" class="mt-4" @click="loadInitialData">
          <RefreshCw class="h-4 w-4" />
          Retry
        </Button>
      </div>

      <!-- Main Content -->
      <div v-else-if="run" class="space-y-6">
        <!-- Run Status Header -->
        <Card>
          <CardContent class="p-6">
            <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
              <div>
                <div class="flex items-center gap-3">
                  <h1 class="text-xl font-semibold text-zinc-900 dark:text-zinc-50">
                    Execution #{{ run.id }}
                  </h1>
                  <Badge :variant="getRunConfig(run.status).variant" class="gap-1">
                    <component 
                      :is="getRunConfig(run.status).icon" 
                      :class="['h-3 w-3', run.status === 'running' && 'animate-spin']" 
                    />
                    {{ getRunConfig(run.status).label }}
                  </Badge>
                  <Badge v-if="run.trigger_type === 'trial'" variant="warning">
                    Trial
                  </Badge>
                </div>
                <div class="mt-2 flex flex-wrap gap-4 text-sm text-zinc-500 dark:text-zinc-400">
                  <span class="flex items-center gap-1">
                    <Activity class="h-3.5 w-3.5" />
                    Workflow ID: {{ run.workflow_id }}
                  </span>
                  <span class="flex items-center gap-1">
                    <GitBranch class="h-3.5 w-3.5" />
                    Version: v{{ run.workflow_version }}
                  </span>
                  <span class="flex items-center gap-1">
                    <Zap class="h-3.5 w-3.5" />
                    Trigger: {{ run.trigger_type }}
                  </span>
                </div>
              </div>
              
              <div class="flex flex-col items-end gap-1 text-sm text-zinc-500 dark:text-zinc-400">
                <div>Started: {{ formatTime(run.started_at) }}</div>
                <div v-if="run.finished_at">Finished: {{ formatTime(run.finished_at) }}</div>
                <div class="font-medium text-zinc-900 dark:text-zinc-50">
                  Duration: {{ formatDuration(run.started_at, run.finished_at) }}
                </div>
              </div>
            </div>

            <!-- Error Message -->
            <div 
              v-if="run.error_message" 
              class="mt-4 flex items-start gap-2 rounded-md border border-red-200 bg-red-50 p-3 dark:border-red-900 dark:bg-red-950/50"
            >
              <XCircle class="mt-0.5 h-4 w-4 flex-shrink-0 text-red-500" />
              <div>
                <div class="text-sm font-medium text-red-700 dark:text-red-400">Execution Failed</div>
                <div class="mt-1 text-sm text-red-600 dark:text-red-500">{{ run.error_message }}</div>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Step Strip (Horizontal Timeline) -->
        <Card>
          <CardHeader class="pb-3">
            <CardTitle class="text-base">Step Progress</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="overflow-x-auto pb-2">
              <div class="flex min-w-max items-center gap-0">
                <template v-for="(step, index) in orderedSteps" :key="step.id">
                  <!-- Step Node -->
                  <button
                    @click="selectStep(step.node_id)"
                    :class="[
                      'relative flex min-w-[100px] flex-col items-center rounded-md border-2 p-3 transition-all',
                      selectedStepId === step.node_id 
                        ? 'ring-2 ring-zinc-900 ring-offset-2 dark:ring-zinc-50 dark:ring-offset-zinc-950' 
                        : 'hover:bg-zinc-50 dark:hover:bg-zinc-900',
                      getStepConfig(step.status).borderClass,
                    ]"
                  >
                    <!-- Status Icon -->
                    <div 
                      :class="[
                        'flex h-10 w-10 items-center justify-center rounded-full bg-white dark:bg-zinc-900',
                        step.status === 'running' && 'animate-pulse',
                      ]"
                    >
                      <component 
                        :is="getStepConfig(step.status).icon" 
                        :class="[
                          'h-5 w-5',
                          step.status === 'succeeded' && 'text-emerald-500',
                          step.status === 'failed' && 'text-red-500',
                          step.status === 'running' && 'text-blue-500 animate-spin',
                          step.status === 'pending' && 'text-zinc-400',
                          step.status === 'skipped' && 'text-zinc-400',
                        ]"
                      />
                    </div>
                    
                    <!-- Step Info -->
                    <div class="mt-2 text-center">
                      <div 
                        class="max-w-[80px] truncate text-xs font-medium text-zinc-900 dark:text-zinc-50"
                        :title="run?.definition_snapshot?.nodes?.find(n => n.id === step.node_id)?.name || step.node_id"
                      >
                        {{ run?.definition_snapshot?.nodes?.find(n => n.id === step.node_id)?.name || step.node_id }}
                      </div>
                      <div class="mt-0.5 text-[10px] text-zinc-500 dark:text-zinc-400">
                        {{ step.node_type }}
                      </div>
                      <div v-if="step.started_at" class="mt-0.5 text-[10px] text-zinc-400 dark:text-zinc-500">
                        {{ formatDuration(step.started_at, step.finished_at) }}
                      </div>
                    </div>

                    <!-- Error indicator -->
                    <div 
                      v-if="step.status === 'failed'" 
                      class="absolute -right-1 -top-1 flex h-4 w-4 items-center justify-center rounded-full bg-red-500"
                    >
                      <span class="text-[10px] font-bold text-white">!</span>
                    </div>
                  </button>

                  <!-- Connector Arrow -->
                  <div 
                    v-if="index < orderedSteps.length - 1"
                    class="flex items-center px-1"
                  >
                    <div 
                      :class="[
                        'h-0.5 w-6',
                        step.status === 'succeeded' ? 'bg-emerald-400' : 
                        step.status === 'failed' ? 'bg-red-400' :
                        step.status === 'running' ? 'bg-blue-400' : 
                        'bg-zinc-300 dark:bg-zinc-700'
                      ]"
                    ></div>
                    <ChevronRight 
                      :class="[
                        'h-4 w-4',
                        step.status === 'succeeded' ? 'text-emerald-400' : 
                        step.status === 'failed' ? 'text-red-400' :
                        step.status === 'running' ? 'text-blue-400' : 
                        'text-zinc-300 dark:text-zinc-700'
                      ]"
                    />
                  </div>
                </template>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Status Board (Two Column Layout) -->
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <!-- Steps List -->
          <Card>
            <CardHeader class="pb-3">
              <CardTitle class="text-base">Step Details</CardTitle>
            </CardHeader>
            <CardContent class="space-y-2">
              <button 
                v-for="step in orderedSteps" 
                :key="step.id"
                @click="selectStep(step.node_id)"
                :class="[
                  'flex w-full items-center justify-between rounded-md border-2 p-3 text-left transition-all',
                  selectedStepId === step.node_id 
                    ? 'border-zinc-900 bg-zinc-50 dark:border-zinc-50 dark:bg-zinc-900' 
                    : 'border-transparent hover:border-zinc-200 dark:hover:border-zinc-800',
                ]"
              >
                <div class="flex items-center gap-3">
                  <div 
                    :class="[
                      'flex h-8 w-8 items-center justify-center rounded-full',
                      step.status === 'running' && 'animate-pulse bg-blue-100 dark:bg-blue-900/30',
                      step.status !== 'running' && 'bg-zinc-100 dark:bg-zinc-800',
                    ]"
                  >
                    <component 
                      :is="getStepConfig(step.status).icon" 
                      :class="[
                        'h-4 w-4',
                        step.status === 'succeeded' && 'text-emerald-500',
                        step.status === 'failed' && 'text-red-500',
                        step.status === 'running' && 'text-blue-500 animate-spin',
                        step.status === 'pending' && 'text-zinc-400',
                        step.status === 'skipped' && 'text-zinc-400',
                      ]"
                    />
                  </div>
                  <div>
                    <div class="text-sm font-medium text-zinc-900 dark:text-zinc-50">
                      {{ run?.definition_snapshot?.nodes?.find(n => n.id === step.node_id)?.name || step.node_id }}
                    </div>
                    <div class="text-xs text-zinc-500 dark:text-zinc-400">
                      {{ step.node_type }} · {{ getStepConfig(step.status).label }}
                    </div>
                  </div>
                </div>
                
                <div class="text-right text-xs">
                  <div v-if="step.started_at" class="text-zinc-500 dark:text-zinc-400">
                    {{ formatTime(step.started_at) }}
                  </div>
                  <div v-if="step.started_at" class="text-zinc-400 dark:text-zinc-500">
                    {{ formatDuration(step.started_at, step.finished_at) }}
                  </div>
                </div>
              </button>
            </CardContent>
          </Card>

          <!-- Selected Step Details / Summary -->
          <Card>
            <CardHeader class="pb-3">
              <CardTitle class="text-base">
                {{ selectedStep ? 'Step Information' : 'Execution Summary' }}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <!-- Selected Step Info -->
              <div v-if="selectedStep && selectedNode" class="space-y-3">
                <div class="grid grid-cols-2 gap-3">
                  <div class="rounded-md bg-zinc-100 p-2.5 dark:bg-zinc-900">
                    <div class="text-[10px] uppercase tracking-wide text-zinc-500 dark:text-zinc-400">Node Name</div>
                    <div class="mt-0.5 text-sm font-medium text-zinc-900 dark:text-zinc-50">{{ selectedNode.name }}</div>
                  </div>
                  <div class="rounded-md bg-zinc-100 p-2.5 dark:bg-zinc-900">
                    <div class="text-[10px] uppercase tracking-wide text-zinc-500 dark:text-zinc-400">Node Type</div>
                    <div class="mt-0.5 text-sm font-medium text-zinc-900 dark:text-zinc-50">{{ selectedStep.node_type }}</div>
                  </div>
                  <div class="rounded-md bg-zinc-100 p-2.5 dark:bg-zinc-900">
                    <div class="text-[10px] uppercase tracking-wide text-zinc-500 dark:text-zinc-400">Status</div>
                    <Badge :variant="getStepConfig(selectedStep.status).variant" class="mt-0.5 gap-1">
                      <component :is="getStepConfig(selectedStep.status).icon" class="h-3 w-3" />
                      {{ getStepConfig(selectedStep.status).label }}
                    </Badge>
                  </div>
                  <div class="rounded-md bg-zinc-100 p-2.5 dark:bg-zinc-900">
                    <div class="text-[10px] uppercase tracking-wide text-zinc-500 dark:text-zinc-400">Duration</div>
                    <div class="mt-0.5 text-sm font-medium text-zinc-900 dark:text-zinc-50">
                      {{ formatDuration(selectedStep.started_at, selectedStep.finished_at) }}
                    </div>
                  </div>
                </div>

                <div class="rounded-md bg-zinc-100 p-2.5 dark:bg-zinc-900">
                  <div class="text-[10px] uppercase tracking-wide text-zinc-500 dark:text-zinc-400">Node ID</div>
                  <div class="mt-0.5 font-mono text-xs text-zinc-900 dark:text-zinc-50">{{ selectedStep.node_id }}</div>
                </div>

                <div v-if="selectedStep.started_at" class="rounded-md bg-zinc-100 p-2.5 dark:bg-zinc-900">
                  <div class="text-[10px] uppercase tracking-wide text-zinc-500 dark:text-zinc-400">Timeline</div>
                  <div class="mt-1 space-y-0.5 text-xs text-zinc-900 dark:text-zinc-50">
                    <div>Started: {{ selectedStep.started_at }}</div>
                    <div v-if="selectedStep.finished_at">Finished: {{ selectedStep.finished_at }}</div>
                  </div>
                </div>

                <div v-if="selectedStep.error_message" class="rounded-md border border-red-200 bg-red-50 p-2.5 dark:border-red-900 dark:bg-red-950/50">
                  <div class="text-[10px] uppercase tracking-wide text-red-500 dark:text-red-400">Error Message</div>
                  <div class="mt-1 text-xs text-red-700 dark:text-red-300">{{ selectedStep.error_message }}</div>
                </div>

                <div v-if="selectedStep.input_data" class="rounded-md bg-zinc-100 p-2.5 dark:bg-zinc-900">
                  <div class="text-[10px] uppercase tracking-wide text-zinc-500 dark:text-zinc-400">Input Data</div>
                  <pre class="mt-1 max-h-24 overflow-auto font-mono text-xs text-zinc-700 dark:text-zinc-300">{{ selectedStep.input_data }}</pre>
                </div>

                <div v-if="selectedStep.output_data" class="rounded-md bg-zinc-100 p-2.5 dark:bg-zinc-900">
                  <div class="text-[10px] uppercase tracking-wide text-zinc-500 dark:text-zinc-400">Output Data</div>
                  <pre class="mt-1 max-h-24 overflow-auto font-mono text-xs text-zinc-700 dark:text-zinc-300">{{ selectedStep.output_data }}</pre>
                </div>
              </div>

              <!-- Summary View (when no step selected) -->
              <div v-else class="space-y-4">
                <div class="grid grid-cols-2 gap-3">
                  <div class="rounded-md bg-zinc-100 p-4 text-center dark:bg-zinc-900">
                    <div class="text-2xl font-bold text-zinc-900 dark:text-zinc-50">{{ steps.length }}</div>
                    <div class="text-xs text-zinc-500 dark:text-zinc-400">Total Steps</div>
                  </div>
                  <div class="rounded-md bg-emerald-50 p-4 text-center dark:bg-emerald-950/30">
                    <div class="text-2xl font-bold text-emerald-600 dark:text-emerald-400">
                      {{ steps.filter(s => s.status === 'succeeded').length }}
                    </div>
                    <div class="text-xs text-zinc-500 dark:text-zinc-400">Completed</div>
                  </div>
                  <div class="rounded-md bg-blue-50 p-4 text-center dark:bg-blue-950/30">
                    <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">
                      {{ steps.filter(s => s.status === 'running').length }}
                    </div>
                    <div class="text-xs text-zinc-500 dark:text-zinc-400">Running</div>
                  </div>
                  <div class="rounded-md bg-red-50 p-4 text-center dark:bg-red-950/30">
                    <div class="text-2xl font-bold text-red-600 dark:text-red-400">
                      {{ steps.filter(s => s.status === 'failed').length }}
                    </div>
                    <div class="text-xs text-zinc-500 dark:text-zinc-400">Failed</div>
                  </div>
                </div>

                <div class="rounded-md bg-zinc-100 p-4 dark:bg-zinc-900">
                  <div class="mb-2 text-xs text-zinc-500 dark:text-zinc-400">Execution Progress</div>
                  <div class="h-2 w-full overflow-hidden rounded-full bg-zinc-200 dark:bg-zinc-800">
                    <div 
                      class="h-full transition-all duration-500"
                      :class="[
                        run.status === 'failed' ? 'bg-red-500' : 
                        run.status === 'succeeded' ? 'bg-emerald-500' : 'bg-blue-500'
                      ]"
                      :style="{ 
                        width: `${steps.length > 0 ? (steps.filter(s => ['succeeded', 'failed'].includes(s.status)).length / steps.length) * 100 : 0}%` 
                      }"
                    ></div>
                  </div>
                  <div class="mt-2 text-right text-xs text-zinc-500 dark:text-zinc-400">
                    {{ steps.filter(s => ['succeeded', 'failed'].includes(s.status)).length }} / {{ steps.length }} steps completed
                  </div>
                </div>

                <p class="py-4 text-center text-sm text-zinc-500 dark:text-zinc-400">
                  Click a step on the left to view details
                </p>
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- Polling Indicator -->
        <div 
          v-if="run && !['succeeded', 'failed', 'cancelled'].includes(run.status)"
          class="fixed bottom-4 right-4 flex items-center gap-2 rounded-full border border-zinc-200 bg-white px-4 py-2 text-sm shadow-lg dark:border-zinc-800 dark:bg-zinc-950"
        >
          <Loader2 class="h-4 w-4 animate-spin text-blue-500" />
          <span class="text-zinc-600 dark:text-zinc-400">Live updates...</span>
        </div>
      </div>
    </main>
  </div>
</template>
