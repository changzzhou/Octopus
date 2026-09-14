<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()
const workflowId = route.params.id as string
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
      <div class="max-w-4xl mx-auto">
        <div class="bg-white dark:bg-slate-800 rounded-xl shadow-lg p-8">
          <div class="flex items-center justify-between mb-6">
            <h1 class="text-3xl font-bold text-slate-900 dark:text-white">
              Workflow Design
            </h1>
            <span class="px-3 py-1 bg-purple-100 dark:bg-purple-900 text-purple-700 dark:text-purple-300 rounded-full text-sm font-medium">
              ID: {{ workflowId }}
            </span>
          </div>

          <div class="p-8 bg-purple-50 dark:bg-purple-900/20 rounded-xl border-2 border-dashed border-purple-200 dark:border-purple-800">
            <div class="text-center">
              <div class="text-6xl mb-4">🎨</div>
              <h3 class="text-xl font-medium text-purple-700 dark:text-purple-300 mb-2">
                Design-Time Canvas
              </h3>
              <p class="text-purple-600 dark:text-purple-400 mb-4">
                Vue Flow DAG editor will be integrated here in FE-2
              </p>
              <div class="text-sm text-purple-500 dark:text-purple-500">
                Placeholder for workflow design interface
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
