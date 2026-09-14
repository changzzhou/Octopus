<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()
const workflowId = route.params.id as string
const executionId = route.params.executionId as string
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
            :to="`/workflows/${workflowId}`"
            class="text-slate-600 dark:text-slate-300 hover:text-purple-600 dark:hover:text-purple-400"
          >
            ← Back to Workflow
          </router-link>
        </nav>
      </div>
    </header>

    <main class="container mx-auto px-4 py-8">
      <div class="max-w-4xl mx-auto">
        <div class="bg-white dark:bg-slate-800 rounded-xl shadow-lg p-8">
          <div class="flex items-center justify-between mb-6">
            <h1 class="text-3xl font-bold text-slate-900 dark:text-white">
              Workflow Execution
            </h1>
            <div class="flex gap-2">
              <span class="px-3 py-1 bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 rounded-full text-sm font-medium">
                Workflow: {{ workflowId }}
              </span>
              <span class="px-3 py-1 bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-300 rounded-full text-sm font-medium">
                Execution: {{ executionId }}
              </span>
            </div>
          </div>

          <div class="p-8 bg-green-50 dark:bg-green-900/20 rounded-xl border-2 border-dashed border-green-200 dark:border-green-800">
            <div class="text-center">
              <div class="text-6xl mb-4">▶️</div>
              <h3 class="text-xl font-medium text-green-700 dark:text-green-300 mb-2">
                Execution-Time View
              </h3>
              <p class="text-green-600 dark:text-green-400 mb-4">
                Real-time execution monitoring will be integrated here
              </p>
              <div class="text-sm text-green-500 dark:text-green-500">
                Placeholder for workflow execution status and logs
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
