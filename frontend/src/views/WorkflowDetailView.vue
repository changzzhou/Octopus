<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getWorkflowDetail, type WorkflowDetail } from '../api/workflows'

const route = useRoute()
const workflow = ref<WorkflowDetail | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    const id = Number(route.params.id)
    const data = await getWorkflowDetail(id)
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

          <div class="flex gap-4">
            <router-link
              :to="`/workflows/${workflow.id}/design`"
              class="flex-1 px-6 py-3 bg-purple-600 hover:bg-purple-700 text-white font-medium rounded-lg text-center transition-colors"
            >
              🎨 Open Designer
            </router-link>
          </div>

          <div v-if="workflow.nodes && workflow.nodes.length > 0" class="mt-8">
            <h2 class="text-lg font-semibold text-slate-900 dark:text-white mb-4">
              Workflow Definition
            </h2>
            <div class="text-sm text-slate-600 dark:text-slate-400">
              <p>{{ workflow.nodes.length }} node(s), {{ workflow.edges?.length || 0 }} edge(s)</p>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
