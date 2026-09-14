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

const runStatusConfig: Record<RunStatus, { label: string; color: string; bgColor: string; icon: string }> = {
  pending: { label: '等待中', color: 'text-slate-600', bgColor: 'bg-slate-100 dark:bg-slate-700', icon: '⏳' },
  running: { label: '运行中', color: 'text-blue-600', bgColor: 'bg-blue-100 dark:bg-blue-900', icon: '🔄' },
  succeeded: { label: '成功', color: 'text-green-600', bgColor: 'bg-green-100 dark:bg-green-900', icon: '✅' },
  failed: { label: '失败', color: 'text-red-600', bgColor: 'bg-red-100 dark:bg-red-900', icon: '❌' },
  cancelled: { label: '已取消', color: 'text-orange-600', bgColor: 'bg-orange-100 dark:bg-orange-900', icon: '🚫' },
}

const stepStatusConfig: Record<StepStatus, { label: string; color: string; bgColor: string; borderColor: string; icon: string }> = {
  pending: { label: '等待中', color: 'text-slate-500', bgColor: 'bg-slate-50 dark:bg-slate-800', borderColor: 'border-slate-300 dark:border-slate-600', icon: '○' },
  running: { label: '运行中', color: 'text-blue-600', bgColor: 'bg-blue-50 dark:bg-blue-900/30', borderColor: 'border-blue-400 dark:border-blue-500', icon: '◎' },
  succeeded: { label: '成功', color: 'text-green-600', bgColor: 'bg-green-50 dark:bg-green-900/30', borderColor: 'border-green-400 dark:border-green-500', icon: '●' },
  failed: { label: '失败', color: 'text-red-600', bgColor: 'bg-red-50 dark:bg-red-900/30', borderColor: 'border-red-400 dark:border-red-500', icon: '✕' },
  skipped: { label: '已跳过', color: 'text-slate-400', bgColor: 'bg-slate-50 dark:bg-slate-800', borderColor: 'border-slate-300 dark:border-slate-600', icon: '◌' },
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
    error.value = '加载执行数据失败'
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
  <div class="min-h-screen bg-slate-100 dark:bg-slate-900">
    <!-- Header -->
    <header class="bg-white dark:bg-slate-800 shadow">
      <div class="container mx-auto px-4 py-4 flex items-center justify-between">
        <router-link to="/" class="text-2xl font-bold text-purple-600 dark:text-purple-400">
          🐙 Octopus
        </router-link>
        <nav class="space-x-4">
          <button 
            @click="goBack"
            class="text-slate-600 dark:text-slate-300 hover:text-purple-600 dark:hover:text-purple-400"
          >
            ← 返回工作流
          </button>
        </nav>
      </div>
    </header>

    <main class="container mx-auto px-4 py-6">
      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-4 border-purple-500 border-t-transparent"></div>
        <span class="ml-4 text-slate-600 dark:text-slate-300">加载中...</span>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-20">
        <div class="text-6xl mb-4">⚠️</div>
        <p class="text-red-500 text-lg">{{ error }}</p>
        <button 
          @click="loadInitialData" 
          class="mt-4 px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700"
        >
          重试
        </button>
      </div>

      <!-- Main Content -->
      <div v-else-if="run" class="space-y-6">
        <!-- Run Status Header -->
        <div class="bg-white dark:bg-slate-800 rounded-xl shadow-lg p-6">
          <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
            <div>
              <div class="flex items-center gap-3 mb-2">
                <h1 class="text-2xl font-bold text-slate-900 dark:text-white">
                  执行 #{{ run.id }}
                </h1>
                <span 
                  :class="[
                    getRunConfig(run.status).bgColor,
                    getRunConfig(run.status).color,
                    'px-3 py-1 rounded-full text-sm font-medium flex items-center gap-1'
                  ]"
                >
                  <span>{{ getRunConfig(run.status).icon }}</span>
                  {{ getRunConfig(run.status).label }}
                </span>
                <span 
                  v-if="run.trigger_type === 'trial'"
                  class="px-2 py-1 bg-amber-100 dark:bg-amber-900 text-amber-700 dark:text-amber-300 rounded text-xs font-medium"
                >
                  试运行
                </span>
              </div>
              <div class="text-sm text-slate-500 dark:text-slate-400 flex flex-wrap gap-4">
                <span>工作流 ID: {{ run.workflow_id }}</span>
                <span>版本: v{{ run.workflow_version }}</span>
                <span>触发类型: {{ run.trigger_type }}</span>
              </div>
            </div>
            
            <div class="flex flex-col items-end gap-1 text-sm text-slate-600 dark:text-slate-400">
              <div>开始: {{ formatTime(run.started_at) }}</div>
              <div v-if="run.finished_at">结束: {{ formatTime(run.finished_at) }}</div>
              <div class="font-medium">
                耗时: {{ formatDuration(run.started_at, run.finished_at) }}
              </div>
            </div>
          </div>

          <!-- Error Message -->
          <div 
            v-if="run.error_message" 
            class="mt-4 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg"
          >
            <div class="flex items-start gap-2">
              <span class="text-red-500">❌</span>
              <div>
                <div class="font-medium text-red-700 dark:text-red-300">执行失败</div>
                <div class="text-sm text-red-600 dark:text-red-400 mt-1">{{ run.error_message }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Step Strip (Horizontal Timeline) -->
        <div class="bg-white dark:bg-slate-800 rounded-xl shadow-lg p-6">
          <h2 class="text-lg font-semibold text-slate-900 dark:text-white mb-4">步骤进度</h2>
          
          <div class="overflow-x-auto pb-2">
            <div class="flex items-center min-w-max gap-0">
              <template v-for="(step, index) in orderedSteps" :key="step.id">
                <!-- Step Node -->
                <button
                  @click="selectStep(step.node_id)"
                  :class="[
                    'relative flex flex-col items-center p-3 rounded-lg transition-all duration-200 min-w-[120px]',
                    selectedStepId === step.node_id 
                      ? 'ring-2 ring-purple-500 ring-offset-2 dark:ring-offset-slate-800' 
                      : 'hover:bg-slate-50 dark:hover:bg-slate-700/50',
                    getStepConfig(step.status).bgColor,
                    'border-2',
                    getStepConfig(step.status).borderColor,
                  ]"
                >
                  <!-- Status Icon -->
                  <div 
                    :class="[
                      'w-10 h-10 rounded-full flex items-center justify-center text-lg font-bold mb-2',
                      step.status === 'running' ? 'animate-pulse' : '',
                      getStepConfig(step.status).color,
                      step.status === 'running' ? 'bg-blue-200 dark:bg-blue-800' : 'bg-white dark:bg-slate-700',
                    ]"
                  >
                    {{ getStepConfig(step.status).icon }}
                  </div>
                  
                  <!-- Step Info -->
                  <div class="text-center">
                    <div 
                      :class="[
                        'text-sm font-medium truncate max-w-[100px]',
                        getStepConfig(step.status).color
                      ]"
                      :title="run?.definition_snapshot?.nodes?.find(n => n.id === step.node_id)?.name || step.node_id"
                    >
                      {{ run?.definition_snapshot?.nodes?.find(n => n.id === step.node_id)?.name || step.node_id }}
                    </div>
                    <div class="text-xs text-slate-500 dark:text-slate-400 mt-1">
                      {{ step.node_type }}
                    </div>
                    <div 
                      v-if="step.started_at" 
                      class="text-xs text-slate-400 dark:text-slate-500 mt-1"
                    >
                      {{ formatDuration(step.started_at, step.finished_at) }}
                    </div>
                  </div>

                  <!-- Error indicator -->
                  <div 
                    v-if="step.status === 'failed'" 
                    class="absolute -top-1 -right-1 w-4 h-4 bg-red-500 rounded-full flex items-center justify-center"
                  >
                    <span class="text-white text-xs">!</span>
                  </div>
                </button>

                <!-- Connector Arrow -->
                <div 
                  v-if="index < orderedSteps.length - 1"
                  class="flex items-center px-2"
                >
                  <div 
                    :class="[
                      'h-0.5 w-8',
                      step.status === 'succeeded' ? 'bg-green-400' : 
                      step.status === 'failed' ? 'bg-red-400' :
                      step.status === 'running' ? 'bg-blue-400 animate-pulse' : 
                      'bg-slate-300 dark:bg-slate-600'
                    ]"
                  ></div>
                  <div 
                    :class="[
                      'w-0 h-0 border-y-4 border-y-transparent border-l-8',
                      step.status === 'succeeded' ? 'border-l-green-400' : 
                      step.status === 'failed' ? 'border-l-red-400' :
                      step.status === 'running' ? 'border-l-blue-400' : 
                      'border-l-slate-300 dark:border-l-slate-600'
                    ]"
                  ></div>
                </div>
              </template>
            </div>
          </div>
        </div>

        <!-- Status Board (Two Column Layout) -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Steps List -->
          <div class="bg-white dark:bg-slate-800 rounded-xl shadow-lg p-6">
            <h2 class="text-lg font-semibold text-slate-900 dark:text-white mb-4">步骤详情</h2>
            
            <div class="space-y-3">
              <div 
                v-for="step in orderedSteps" 
                :key="step.id"
                @click="selectStep(step.node_id)"
                :class="[
                  'p-4 rounded-lg cursor-pointer transition-all duration-200 border-2',
                  selectedStepId === step.node_id 
                    ? 'border-purple-500 bg-purple-50 dark:bg-purple-900/20' 
                    : 'border-transparent hover:border-slate-200 dark:hover:border-slate-600',
                  getStepConfig(step.status).bgColor,
                ]"
              >
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-3">
                    <span 
                      :class="[
                        'w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold',
                        step.status === 'running' ? 'animate-pulse bg-blue-200 dark:bg-blue-800' : 'bg-white dark:bg-slate-700',
                        getStepConfig(step.status).color
                      ]"
                    >
                      {{ getStepConfig(step.status).icon }}
                    </span>
                    <div>
                      <div class="font-medium text-slate-900 dark:text-white">
                        {{ run?.definition_snapshot?.nodes?.find(n => n.id === step.node_id)?.name || step.node_id }}
                      </div>
                      <div class="text-xs text-slate-500 dark:text-slate-400">
                        {{ step.node_type }} · {{ getStepConfig(step.status).label }}
                      </div>
                    </div>
                  </div>
                  
                  <div class="text-right text-sm">
                    <div v-if="step.started_at" class="text-slate-500 dark:text-slate-400">
                      {{ formatTime(step.started_at) }}
                    </div>
                    <div v-if="step.started_at" class="text-xs text-slate-400 dark:text-slate-500">
                      {{ formatDuration(step.started_at, step.finished_at) }}
                    </div>
                  </div>
                </div>

                <!-- Step Error -->
                <div 
                  v-if="step.error_message" 
                  class="mt-3 p-2 bg-red-100 dark:bg-red-900/30 rounded text-sm text-red-600 dark:text-red-400"
                >
                  {{ step.error_message }}
                </div>
              </div>
            </div>
          </div>

          <!-- Selected Step Details / Summary -->
          <div class="bg-white dark:bg-slate-800 rounded-xl shadow-lg p-6">
            <h2 class="text-lg font-semibold text-slate-900 dark:text-white mb-4">
              {{ selectedStep ? '步骤信息' : '执行摘要' }}
            </h2>
            
            <!-- Selected Step Info -->
            <div v-if="selectedStep && selectedNode" class="space-y-4">
              <div class="grid grid-cols-2 gap-4">
                <div class="bg-slate-50 dark:bg-slate-700 rounded-lg p-3">
                  <div class="text-xs text-slate-500 dark:text-slate-400">节点名称</div>
                  <div class="font-medium text-slate-900 dark:text-white">{{ selectedNode.name }}</div>
                </div>
                <div class="bg-slate-50 dark:bg-slate-700 rounded-lg p-3">
                  <div class="text-xs text-slate-500 dark:text-slate-400">节点类型</div>
                  <div class="font-medium text-slate-900 dark:text-white">{{ selectedStep.node_type }}</div>
                </div>
                <div class="bg-slate-50 dark:bg-slate-700 rounded-lg p-3">
                  <div class="text-xs text-slate-500 dark:text-slate-400">状态</div>
                  <div :class="['font-medium', getStepConfig(selectedStep.status).color]">
                    {{ getStepConfig(selectedStep.status).icon }} {{ getStepConfig(selectedStep.status).label }}
                  </div>
                </div>
                <div class="bg-slate-50 dark:bg-slate-700 rounded-lg p-3">
                  <div class="text-xs text-slate-500 dark:text-slate-400">耗时</div>
                  <div class="font-medium text-slate-900 dark:text-white">
                    {{ formatDuration(selectedStep.started_at, selectedStep.finished_at) }}
                  </div>
                </div>
              </div>

              <div class="bg-slate-50 dark:bg-slate-700 rounded-lg p-3">
                <div class="text-xs text-slate-500 dark:text-slate-400 mb-1">节点 ID</div>
                <div class="font-mono text-sm text-slate-900 dark:text-white">{{ selectedStep.node_id }}</div>
              </div>

              <div v-if="selectedStep.started_at" class="bg-slate-50 dark:bg-slate-700 rounded-lg p-3">
                <div class="text-xs text-slate-500 dark:text-slate-400 mb-1">时间线</div>
                <div class="text-sm text-slate-900 dark:text-white space-y-1">
                  <div>开始: {{ selectedStep.started_at }}</div>
                  <div v-if="selectedStep.finished_at">结束: {{ selectedStep.finished_at }}</div>
                </div>
              </div>

              <div v-if="selectedStep.error_message" class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-3">
                <div class="text-xs text-red-500 dark:text-red-400 mb-1">错误信息</div>
                <div class="text-sm text-red-700 dark:text-red-300">{{ selectedStep.error_message }}</div>
              </div>

              <div v-if="selectedStep.input_data" class="bg-slate-50 dark:bg-slate-700 rounded-lg p-3">
                <div class="text-xs text-slate-500 dark:text-slate-400 mb-1">输入数据</div>
                <pre class="text-xs text-slate-700 dark:text-slate-300 overflow-auto max-h-32">{{ selectedStep.input_data }}</pre>
              </div>

              <div v-if="selectedStep.output_data" class="bg-slate-50 dark:bg-slate-700 rounded-lg p-3">
                <div class="text-xs text-slate-500 dark:text-slate-400 mb-1">输出数据</div>
                <pre class="text-xs text-slate-700 dark:text-slate-300 overflow-auto max-h-32">{{ selectedStep.output_data }}</pre>
              </div>
            </div>

            <!-- Summary View (when no step selected) -->
            <div v-else class="space-y-4">
              <div class="grid grid-cols-2 gap-4">
                <div class="bg-slate-50 dark:bg-slate-700 rounded-lg p-4 text-center">
                  <div class="text-3xl font-bold text-slate-900 dark:text-white">{{ steps.length }}</div>
                  <div class="text-sm text-slate-500 dark:text-slate-400">总步骤数</div>
                </div>
                <div class="bg-green-50 dark:bg-green-900/20 rounded-lg p-4 text-center">
                  <div class="text-3xl font-bold text-green-600 dark:text-green-400">
                    {{ steps.filter(s => s.status === 'succeeded').length }}
                  </div>
                  <div class="text-sm text-slate-500 dark:text-slate-400">已完成</div>
                </div>
                <div class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4 text-center">
                  <div class="text-3xl font-bold text-blue-600 dark:text-blue-400">
                    {{ steps.filter(s => s.status === 'running').length }}
                  </div>
                  <div class="text-sm text-slate-500 dark:text-slate-400">运行中</div>
                </div>
                <div class="bg-red-50 dark:bg-red-900/20 rounded-lg p-4 text-center">
                  <div class="text-3xl font-bold text-red-600 dark:text-red-400">
                    {{ steps.filter(s => s.status === 'failed').length }}
                  </div>
                  <div class="text-sm text-slate-500 dark:text-slate-400">失败</div>
                </div>
              </div>

              <div class="bg-slate-50 dark:bg-slate-700 rounded-lg p-4">
                <div class="text-sm text-slate-500 dark:text-slate-400 mb-2">执行进度</div>
                <div class="w-full bg-slate-200 dark:bg-slate-600 rounded-full h-3 overflow-hidden">
                  <div 
                    class="h-full transition-all duration-500 rounded-full"
                    :class="[
                      run.status === 'failed' ? 'bg-red-500' : 
                      run.status === 'succeeded' ? 'bg-green-500' : 'bg-blue-500'
                    ]"
                    :style="{ 
                      width: `${steps.length > 0 ? (steps.filter(s => ['succeeded', 'failed'].includes(s.status)).length / steps.length) * 100 : 0}%` 
                    }"
                  ></div>
                </div>
                <div class="text-xs text-slate-500 dark:text-slate-400 mt-2 text-right">
                  {{ steps.filter(s => ['succeeded', 'failed'].includes(s.status)).length }} / {{ steps.length }} 步骤已完成
                </div>
              </div>

              <div class="text-center text-sm text-slate-500 dark:text-slate-400 py-4">
                点击左侧步骤查看详情
              </div>
            </div>
          </div>
        </div>

        <!-- Polling Indicator -->
        <div 
          v-if="run && !['succeeded', 'failed', 'cancelled'].includes(run.status)"
          class="fixed bottom-4 right-4 bg-blue-600 text-white px-4 py-2 rounded-full shadow-lg flex items-center gap-2 text-sm"
        >
          <div class="animate-spin w-4 h-4 border-2 border-white border-t-transparent rounded-full"></div>
          <span>实时更新中...</span>
        </div>
      </div>
    </main>
  </div>
</template>
