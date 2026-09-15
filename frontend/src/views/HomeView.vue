<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { checkHealth } from '../api/workflows'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Workflow, ArrowRight, Server, Activity } from '@lucide/vue'

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
  <div class="min-h-screen bg-zinc-50 dark:bg-zinc-950">
    <div class="container mx-auto px-4 py-16">
      <div class="mx-auto max-w-2xl text-center">
        <!-- Logo & Title -->
        <div class="mb-8 flex items-center justify-center gap-3">
          <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-zinc-900 dark:bg-zinc-50">
            <Workflow class="h-6 w-6 text-zinc-50 dark:text-zinc-900" />
          </div>
          <h1 class="text-3xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50">
            Octopus
          </h1>
        </div>
        
        <p class="text-lg text-zinc-600 dark:text-zinc-400">
          Universal Workflow Platform
        </p>
        <p class="mt-1 text-sm text-zinc-500 dark:text-zinc-500">
          通用工作流平台
        </p>

        <!-- System Status Card -->
        <Card class="mx-auto mt-12 max-w-sm">
          <CardHeader class="pb-3">
            <CardTitle class="flex items-center justify-center gap-2 text-base">
              <Server class="h-4 w-4" />
              System Status
            </CardTitle>
          </CardHeader>
          <CardContent class="space-y-3">
            <div v-if="health" class="space-y-2">
              <div class="flex items-center justify-between text-sm">
                <span class="text-zinc-500 dark:text-zinc-400">Backend</span>
                <Badge variant="success">{{ health.status }}</Badge>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span class="text-zinc-500 dark:text-zinc-400">API Version</span>
                <span class="font-mono text-zinc-900 dark:text-zinc-50">{{ health.version }}</span>
              </div>
            </div>
            
            <div v-else-if="error" class="flex items-center justify-center gap-2 text-sm text-red-500">
              <Activity class="h-4 w-4" />
              {{ error }}
            </div>
            
            <div v-else class="text-sm text-zinc-500">
              Checking backend status...
            </div>
          </CardContent>
        </Card>

        <!-- CTA -->
        <div class="mt-12">
          <Button as="router-link" to="/workflows" size="lg">
            View Workflows
            <ArrowRight class="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>
