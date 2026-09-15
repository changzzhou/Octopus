<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { listWorkflows, createWorkflow, triggerRun, type WorkflowSummary, type WorkflowStatus } from '../api/workflows'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Card, CardContent } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { 
  Plus, 
  Play, 
  Pencil, 
  Workflow, 
  FileText,
  Loader2,
  AlertCircle,
} from '@lucide/vue'

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
const triggeringId = ref<number | null>(null)

const statusVariantMap: Record<WorkflowStatus, 'secondary' | 'success' | 'destructive'> = {
  draft: 'secondary',
  enabled: 'success',
  disabled: 'destructive',
}

function getStatusVariant(status: WorkflowStatus) {
  return statusVariantMap[status] || 'secondary'
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

async function handleTriggerRun(workflow: WorkflowSummary, event: Event) {
  event.stopPropagation()
  event.preventDefault()
  
  if (triggeringId.value) return
  
  triggeringId.value = workflow.id
  try {
    const result = await triggerRun(workflow.id)
    router.push(`/workflows/${workflow.id}/executions/${result.run_id}`)
  } catch (e) {
    console.error('Failed to trigger run:', e)
  } finally {
    triggeringId.value = null
  }
}

onMounted(loadWorkflows)
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
            class="text-sm text-zinc-600 transition-colors hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-50"
          >
            Workflows
          </router-link>
        </nav>
      </div>
    </header>

    <main class="container mx-auto px-4 py-8">
      <!-- Page Header -->
      <div class="mb-8 flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-semibold tracking-tight text-zinc-900 dark:text-zinc-50">
            Workflows
          </h1>
          <p class="mt-1 text-sm text-zinc-500 dark:text-zinc-400">
            Manage and run your automation workflows
          </p>
        </div>
        <Button @click="openCreateModal">
          <Plus class="h-4 w-4" />
          New Workflow
        </Button>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-16">
        <Loader2 class="h-6 w-6 animate-spin text-zinc-400" />
        <span class="ml-2 text-sm text-zinc-500">Loading workflows...</span>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="flex flex-col items-center justify-center py-16">
        <AlertCircle class="h-10 w-10 text-red-500" />
        <p class="mt-2 text-sm text-red-600">{{ error }}</p>
        <Button variant="outline" class="mt-4" @click="loadWorkflows">
          Try Again
        </Button>
      </div>

      <!-- Empty State -->
      <Card v-else-if="workflows.length === 0" class="mx-auto max-w-md">
        <CardContent class="flex flex-col items-center py-12">
          <FileText class="h-12 w-12 text-zinc-300 dark:text-zinc-700" />
          <h2 class="mt-4 font-medium text-zinc-900 dark:text-zinc-50">
            No workflows yet
          </h2>
          <p class="mt-1 text-center text-sm text-zinc-500 dark:text-zinc-400">
            Create your first workflow to get started with automation.
          </p>
          <Button class="mt-6" @click="openCreateModal">
            <Plus class="h-4 w-4" />
            Create Workflow
          </Button>
        </CardContent>
      </Card>

      <!-- Workflows Table -->
      <Card v-else>
        <Table>
          <TableHeader>
            <TableRow class="hover:bg-transparent">
              <TableHead>Name</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Updated</TableHead>
              <TableHead class="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow 
              v-for="workflow in workflows" 
              :key="workflow.id"
            >
              <TableCell>
                <router-link 
                  :to="`/workflows/${workflow.id}`"
                  class="font-medium text-zinc-900 hover:underline dark:text-zinc-50"
                >
                  {{ workflow.name }}
                </router-link>
                <p v-if="workflow.description" class="mt-0.5 text-sm text-zinc-500 dark:text-zinc-400">
                  {{ workflow.description }}
                </p>
              </TableCell>
              <TableCell>
                <Badge :variant="getStatusVariant(workflow.status)">
                  {{ workflow.status }}
                </Badge>
              </TableCell>
              <TableCell class="text-zinc-500 dark:text-zinc-400">
                {{ formatUpdatedTime(workflow.updated_at) }}
              </TableCell>
              <TableCell class="text-right">
                <div class="flex items-center justify-end gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    @click="handleTriggerRun(workflow, $event)"
                    :disabled="triggeringId === workflow.id"
                  >
                    <Loader2 v-if="triggeringId === workflow.id" class="h-3.5 w-3.5 animate-spin" />
                    <Play v-else class="h-3.5 w-3.5" />
                    {{ workflow.status === 'draft' ? 'Trial' : 'Run' }}
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    as="router-link"
                    :to="`/workflows/${workflow.id}/design`"
                  >
                    <Pencil class="h-3.5 w-3.5" />
                    Design
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </Card>

      <!-- Total Count -->
      <div v-if="total > 0" class="mt-4 text-center text-sm text-zinc-500 dark:text-zinc-400">
        {{ total }} workflow(s) total
      </div>
    </main>

    <!-- Create Workflow Dialog -->
    <Dialog :open="showCreateModal" @update:open="(val) => val ? null : closeCreateModal()">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create New Workflow</DialogTitle>
          <DialogDescription>
            Add a new workflow to automate your processes.
          </DialogDescription>
        </DialogHeader>
        
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="name">Name</Label>
            <Input 
              id="name"
              v-model="newWorkflowName"
              placeholder="Enter workflow name"
              @keyup.enter="handleCreate"
            />
          </div>
          
          <div class="space-y-2">
            <Label for="description">Description</Label>
            <Textarea 
              id="description"
              v-model="newWorkflowDescription"
              placeholder="Enter description (optional)"
              :rows="3"
            />
          </div>
          
          <p v-if="createError" class="text-sm text-red-500">
            {{ createError }}
          </p>
        </div>
        
        <DialogFooter>
          <Button variant="outline" @click="closeCreateModal">
            Cancel
          </Button>
          <Button @click="handleCreate" :disabled="creating">
            <Loader2 v-if="creating" class="h-4 w-4 animate-spin" />
            {{ creating ? 'Creating...' : 'Create' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
