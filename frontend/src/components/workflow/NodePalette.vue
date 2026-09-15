<script setup lang="ts">
import { NODE_TYPE_CONFIGS, type NodeType } from '../../types/workflow'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { FileCode, Globe, User, CircleCheck, CircleX } from '@lucide/vue'

const emit = defineEmits<{
  (e: 'dragstart', type: NodeType, event: DragEvent): void
}>()

const nodeTypes = Object.entries(NODE_TYPE_CONFIGS) as [NodeType, typeof NODE_TYPE_CONFIGS[NodeType]][]

const nodeIcons: Record<NodeType, any> = {
  script: FileCode,
  http: Globe,
  human: User,
}

function onDragStart(type: NodeType, event: DragEvent) {
  if (event.dataTransfer) {
    event.dataTransfer.setData('application/vueflow-nodetype', type)
    event.dataTransfer.effectAllowed = 'move'
  }
  emit('dragstart', type, event)
}
</script>

<template>
  <Card class="w-48">
    <CardHeader class="px-4 py-3">
      <CardTitle class="text-sm">Node Palette</CardTitle>
      <p class="text-xs text-zinc-500 dark:text-zinc-400">Drag nodes to canvas</p>
    </CardHeader>
    <CardContent class="space-y-2 px-4 pb-4">
      <div
        v-for="[type, config] in nodeTypes"
        :key="type"
        class="flex cursor-grab items-center gap-2 rounded-md border border-zinc-200 p-2 transition-all hover:-translate-y-0.5 hover:shadow-md active:cursor-grabbing dark:border-zinc-800"
        :style="{ borderLeftColor: config.color, borderLeftWidth: '3px' }"
        draggable="true"
        @dragstart="onDragStart(type, $event)"
      >
        <div 
          class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-md"
          :style="{ backgroundColor: config.color }"
        >
          <component :is="nodeIcons[type]" class="h-4 w-4 text-white" />
        </div>
        <div class="min-w-0 flex-1">
          <div class="text-xs font-medium text-zinc-900 dark:text-zinc-50">{{ config.label }}</div>
          <div class="font-mono text-[10px] text-zinc-500 dark:text-zinc-400">{{ type }}</div>
        </div>
      </div>
      
      <Separator class="my-3" />
      
      <div class="space-y-1">
        <p class="text-[10px] font-medium uppercase tracking-wide text-zinc-500 dark:text-zinc-400">Edge Types</p>
        <div class="flex items-center gap-2 text-xs text-zinc-600 dark:text-zinc-400">
          <CircleCheck class="h-3 w-3 text-emerald-500" />
          <span>Success</span>
        </div>
        <div class="flex items-center gap-2 text-xs text-zinc-600 dark:text-zinc-400">
          <CircleX class="h-3 w-3 text-red-500" />
          <span>Failure</span>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
