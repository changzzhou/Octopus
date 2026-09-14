<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { listWorkflows, type WorkflowSummary } from '../api/workflows'

const workflows = ref<WorkflowSummary[]>([])
const total = ref(0)
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    const data = await listWorkflows()
    workflows.value = data.workflows
    total.value = data.total
  } catch (e) {
    error.value = 'Failed to load workflows'
  } finally {
    loading.value = false
  }
})
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
          @click="() => { /* TODO: FE-1 - Create workflow modal */ }"
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
        <p class="text-slate-500 dark:text-slate-400">
          Create your first workflow to get started.
        </p>
      </div>

      <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <router-link
          v-for="workflow in workflows"
          :key="workflow.id"
          :to="`/workflows/${workflow.id}`"
          class="block p-6 bg-white dark:bg-slate-800 rounded-xl shadow hover:shadow-lg transition-shadow"
        >
          <h3 class="text-lg font-semibold text-slate-900 dark:text-white mb-2">
            {{ workflow.name }}
          </h3>
          <p class="text-sm text-slate-500 dark:text-slate-400 mb-4">
            {{ workflow.description || 'No description' }}
          </p>
          <div class="flex justify-between text-xs text-slate-400">
            <span>ID: {{ workflow.id }}</span>
            <span>{{ new Date(workflow.created_at).toLocaleDateString() }}</span>
          </div>
        </router-link>
      </div>

      <div v-if="total > 0" class="mt-6 text-center text-slate-500">
        Total: {{ total }} workflow(s)
      </div>
    </main>
  </div>
</template>
