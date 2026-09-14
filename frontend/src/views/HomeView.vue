<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { checkHealth } from '../api/workflows'

const health = ref<{ status: string; version: string } | null>(null)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    health.value = await checkHealth()
  } catch (e) {
    error.value = 'Backend not reachable'
  }
})
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
    <div class="container mx-auto px-4 py-16">
      <div class="text-center">
        <h1 class="text-5xl font-bold text-white mb-4">
          🐙 Octopus
        </h1>
        <p class="text-xl text-purple-200 mb-8">
          通用工作流平台 / Universal Workflow Platform
        </p>

        <div class="max-w-md mx-auto bg-white/10 backdrop-blur-lg rounded-2xl p-8 shadow-xl">
          <h2 class="text-2xl font-semibold text-white mb-4">System Status</h2>
          
          <div v-if="health" class="space-y-2">
            <div class="flex justify-between text-purple-200">
              <span>Backend Status:</span>
              <span class="text-green-400 font-medium">{{ health.status }}</span>
            </div>
            <div class="flex justify-between text-purple-200">
              <span>API Version:</span>
              <span class="text-white font-medium">{{ health.version }}</span>
            </div>
          </div>
          
          <div v-else-if="error" class="text-red-400">
            {{ error }}
          </div>
          
          <div v-else class="text-purple-300">
            Checking backend status...
          </div>
        </div>

        <nav class="mt-12 space-x-4">
          <router-link 
            to="/workflows" 
            class="inline-block px-6 py-3 bg-purple-600 hover:bg-purple-700 text-white font-medium rounded-lg transition-colors"
          >
            View Workflows
          </router-link>
        </nav>
      </div>
    </div>
  </div>
</template>
