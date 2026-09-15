<script setup lang="ts">
import { computed } from 'vue'
import type { ScriptNodeConfig, HttpNodeConfig, HumanNodeConfig, NodeType, WorkflowNodeData } from '../../types/workflow'
import { NODE_TYPE_CONFIGS } from '../../types/workflow'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from '@/components/ui/select'
import { X, Trash2, FileCode, Globe, User, MousePointer } from '@lucide/vue'

const props = defineProps<{
  node: any | null
}>()

const emit = defineEmits<{
  (e: 'update', nodeId: string, data: WorkflowNodeData): void
  (e: 'delete', nodeId: string): void
  (e: 'close'): void
}>()

const nodeConfig = computed(() => {
  if (!props.node?.data) return null
  return NODE_TYPE_CONFIGS[props.node.data.type as NodeType]
})

const nodeIcons: Record<NodeType, any> = {
  script: FileCode,
  http: Globe,
  human: User,
}

function updateLabel(value: string) {
  if (props.node?.data) {
    emit('update', props.node.id, { ...props.node.data, label: value })
  }
}

function updateConfig(key: string, value: any) {
  if (props.node?.data) {
    emit('update', props.node.id, {
      ...props.node.data,
      config: { ...props.node.data.config, [key]: value }
    })
  }
}

function deleteNode() {
  if (props.node && confirm('Delete this node?')) {
    emit('delete', props.node.id)
  }
}

function handleHeadersInput(value: string) {
  try {
    const parsed = JSON.parse(value)
    updateConfig('headers', parsed)
  } catch {
    // Ignore invalid JSON while typing
  }
}
</script>

<template>
  <Card class="w-64">
    <!-- Node Selected State -->
    <template v-if="node && nodeConfig">
      <CardHeader 
        class="flex flex-row items-center justify-between px-4 py-3"
        :style="{ backgroundColor: nodeConfig.color }"
      >
        <div class="flex items-center gap-2 text-white">
          <component :is="nodeIcons[node.data?.type as NodeType]" class="h-4 w-4" />
          <CardTitle class="text-sm text-white">{{ nodeConfig.label }}</CardTitle>
        </div>
        <Button 
          variant="ghost" 
          size="icon" 
          class="h-6 w-6 text-white hover:bg-white/20 hover:text-white"
          @click="$emit('close')"
        >
          <X class="h-4 w-4" />
        </Button>
      </CardHeader>
      
      <CardContent v-if="node.data" class="space-y-4 p-4">
        <!-- Label Field -->
        <div class="space-y-1.5">
          <Label class="text-xs">Label</Label>
          <Input
            :model-value="node.data.label"
            @update:model-value="updateLabel($event as string)"
            placeholder="Node label"
          />
        </div>
        
        <!-- Script Node Config -->
        <template v-if="node.data.type === 'script'">
          <div class="space-y-1.5">
            <Label class="text-xs">Language</Label>
            <Select
              :model-value="(node.data.config as ScriptNodeConfig).language"
              @update:model-value="updateConfig('language', $event)"
            >
              <SelectTrigger>
                <SelectValue placeholder="Select language" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="javascript">JavaScript</SelectItem>
                <SelectItem value="python">Python</SelectItem>
                <SelectItem value="shell">Shell</SelectItem>
              </SelectContent>
            </Select>
          </div>
          
          <div class="space-y-1.5">
            <Label class="text-xs">Code</Label>
            <Textarea
              :model-value="(node.data.config as ScriptNodeConfig).code"
              @update:model-value="updateConfig('code', $event)"
              :rows="6"
              placeholder="Enter your code here..."
              class="font-mono text-xs"
            />
          </div>
          
          <div class="space-y-1.5">
            <Label class="text-xs">Timeout (ms)</Label>
            <Input
              type="number"
              :model-value="(node.data.config as ScriptNodeConfig).timeout"
              @update:model-value="updateConfig('timeout', Number($event))"
              :min="0"
            />
          </div>
        </template>
        
        <!-- HTTP Node Config -->
        <template v-else-if="node.data.type === 'http'">
          <div class="space-y-1.5">
            <Label class="text-xs">Method</Label>
            <Select
              :model-value="(node.data.config as HttpNodeConfig).method"
              @update:model-value="updateConfig('method', $event)"
            >
              <SelectTrigger>
                <SelectValue placeholder="Select method" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="GET">GET</SelectItem>
                <SelectItem value="POST">POST</SelectItem>
                <SelectItem value="PUT">PUT</SelectItem>
                <SelectItem value="DELETE">DELETE</SelectItem>
                <SelectItem value="PATCH">PATCH</SelectItem>
              </SelectContent>
            </Select>
          </div>
          
          <div class="space-y-1.5">
            <Label class="text-xs">URL</Label>
            <Input
              :model-value="(node.data.config as HttpNodeConfig).url"
              @update:model-value="updateConfig('url', $event)"
              placeholder="https://api.example.com/endpoint"
            />
          </div>
          
          <div class="space-y-1.5">
            <Label class="text-xs">Headers (JSON)</Label>
            <Textarea
              :model-value="JSON.stringify((node.data.config as HttpNodeConfig).headers || {}, null, 2)"
              @update:model-value="handleHeadersInput($event as string)"
              :rows="3"
              placeholder='{"Authorization": "Bearer ..."}'
              class="font-mono text-xs"
            />
          </div>
          
          <div class="space-y-1.5">
            <Label class="text-xs">Body</Label>
            <Textarea
              :model-value="(node.data.config as HttpNodeConfig).body"
              @update:model-value="updateConfig('body', $event)"
              :rows="3"
              placeholder="Request body (JSON, text, etc.)"
              class="font-mono text-xs"
            />
          </div>
          
          <div class="space-y-1.5">
            <Label class="text-xs">Timeout (ms)</Label>
            <Input
              type="number"
              :model-value="(node.data.config as HttpNodeConfig).timeout"
              @update:model-value="updateConfig('timeout', Number($event))"
              :min="0"
            />
          </div>
        </template>
        
        <!-- Human Node Config -->
        <template v-else-if="node.data.type === 'human'">
          <div class="space-y-1.5">
            <Label class="text-xs">Assignee</Label>
            <Input
              :model-value="(node.data.config as HumanNodeConfig).assignee"
              @update:model-value="updateConfig('assignee', $event)"
              placeholder="user@example.com"
            />
          </div>
          
          <div class="space-y-1.5">
            <Label class="text-xs">Instructions</Label>
            <Textarea
              :model-value="(node.data.config as HumanNodeConfig).instructions"
              @update:model-value="updateConfig('instructions', $event)"
              :rows="4"
              placeholder="Enter instructions for the human task..."
            />
          </div>
        </template>
        
        <Separator />
        
        <Button variant="destructive" class="w-full" @click="deleteNode">
          <Trash2 class="h-4 w-4" />
          Delete Node
        </Button>
      </CardContent>
    </template>
    
    <!-- Empty State -->
    <template v-else>
      <CardContent class="flex flex-col items-center justify-center py-12 text-center">
        <MousePointer class="h-8 w-8 text-zinc-300 dark:text-zinc-700" />
        <p class="mt-3 text-sm text-zinc-500 dark:text-zinc-400">Select a node to configure</p>
      </CardContent>
    </template>
  </Card>
</template>
