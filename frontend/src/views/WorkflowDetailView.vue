<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getWorkflow, type Workflow } from '../api/workflows'

const route = useRoute()
const workflow = ref<Workflow | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    const id = Number(route.params.id)
    const data = await getWorkflow(id)
    workflow.value = data.workflow
  } catch (e) {
    error.value = 'Failed to load workflow'
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
          <h1 class="text-3xl font-bold text-slate-900 dark:text-white mb-4">
            {{ workflow.name }}
          </h1>
          
          <p class="text-slate-600 dark:text-slate-400 mb-6">
            {{ workflow.description || 'No description' }}
          </p>

          <div class="grid grid-cols-2 gap-4 text-sm">
            <div class="bg-slate-50 dark:bg-slate-700 rounded-lg p-4">
              <div class="text-slate-500 dark:text-slate-400">ID</div>
              <div class="font-mono text-slate-900 dark:text-white">{{ workflow.id }}</div>
            </div>
            <div class="bg-slate-50 dark:bg-slate-700 rounded-lg p-4">
              <div class="text-slate-500 dark:text-slate-400">Status</div>
              <div class="font-medium text-slate-900 dark:text-white capitalize">
                {{ workflow.status }}
              </div>
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

          <div class="mt-8 p-6 bg-purple-50 dark:bg-purple-900/20 rounded-xl border-2 border-dashed border-purple-200 dark:border-purple-800">
            <div class="text-center">
              <div class="text-4xl mb-2">🎨</div>
              <h3 class="text-lg font-medium text-purple-700 dark:text-purple-300 mb-2">
                DAG Canvas Coming Soon
              </h3>
              <p class="text-sm text-purple-600 dark:text-purple-400">
                Vue Flow integration will be added in FE-2
              </p>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
