<script setup lang="ts">
import { Handle, Position } from '@vue-flow/core'
import type { WorkflowNodeData, HumanNodeConfig } from '../../types/workflow'
import { User, Mail, Check, X } from '@lucide/vue'

defineProps<{
  id: string
  data: WorkflowNodeData
  selected?: boolean
}>()
</script>

<template>
  <div
    class="workflow-node human-node"
    :class="{ selected }"
  >
    <Handle
      type="target"
      :position="Position.Left"
      class="handle-input"
    />
    
    <div class="node-header">
      <User class="h-4 w-4" />
      <span class="node-title">{{ data.label || 'Human Task' }}</span>
    </div>
    
    <div class="node-content">
      <div v-if="(data.config as HumanNodeConfig).assignee" class="node-assignee">
        <Mail class="h-3 w-3" />
        {{ (data.config as HumanNodeConfig).assignee }}
      </div>
      <div class="node-instructions" :title="(data.config as HumanNodeConfig).instructions">
        {{ (data.config as HumanNodeConfig).instructions?.slice(0, 30) }}{{ ((data.config as HumanNodeConfig).instructions?.length || 0) > 30 ? '...' : '' }}
      </div>
    </div>
    
    <Handle
      id="success"
      type="source"
      :position="Position.Right"
      class="handle-success"
      :style="{ top: '35%' }"
    />
    <Handle
      id="failure"
      type="source"
      :position="Position.Right"
      class="handle-failure"
      :style="{ top: '65%' }"
    />
    
    <div class="handle-labels">
      <span class="handle-label success" style="top: 35%">
        <Check class="h-2.5 w-2.5" />
      </span>
      <span class="handle-label failure" style="top: 65%">
        <X class="h-2.5 w-2.5" />
      </span>
    </div>
  </div>
</template>

<style scoped>
.workflow-node {
  background: white;
  border: 2px solid #f59e0b;
  border-radius: 6px;
  padding: 0;
  min-width: 150px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  position: relative;
}

.workflow-node.selected {
  border-color: #18181b;
  box-shadow: 0 0 0 2px rgba(24, 24, 27, 0.2);
}

.node-header {
  background: #f59e0b;
  color: white;
  padding: 6px 10px;
  border-radius: 4px 4px 0 0;
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 500;
  font-size: 12px;
}

.node-content {
  padding: 8px 10px;
}

.node-assignee {
  font-size: 10px;
  color: #a16207;
  background: #fefce8;
  padding: 2px 6px;
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 4px;
}

.node-instructions {
  font-size: 10px;
  color: #71717a;
  line-height: 1.4;
}

.handle-input {
  width: 10px !important;
  height: 10px !important;
  background: #71717a !important;
  border: 2px solid white !important;
}

.handle-success {
  width: 10px !important;
  height: 10px !important;
  background: #10b981 !important;
  border: 2px solid white !important;
}

.handle-failure {
  width: 10px !important;
  height: 10px !important;
  background: #ef4444 !important;
  border: 2px solid white !important;
}

.handle-labels {
  position: absolute;
  right: 14px;
  top: 0;
  bottom: 0;
  pointer-events: none;
}

.handle-label {
  position: absolute;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.handle-label.success {
  color: #10b981;
}

.handle-label.failure {
  color: #ef4444;
}
</style>
