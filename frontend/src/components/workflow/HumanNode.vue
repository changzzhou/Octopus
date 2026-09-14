<script setup lang="ts">
import { Handle, Position } from '@vue-flow/core'
import type { WorkflowNodeData, HumanNodeConfig } from '../../types/workflow'

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
      <span class="node-icon">👤</span>
      <span class="node-title">{{ data.label || 'Human Task' }}</span>
    </div>
    
    <div class="node-content">
      <div v-if="(data.config as HumanNodeConfig).assignee" class="node-assignee">
        <span class="assignee-icon">📧</span>
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
      <span class="handle-label success" style="top: 35%">✓</span>
      <span class="handle-label failure" style="top: 65%">✗</span>
    </div>
  </div>
</template>

<style scoped>
.workflow-node {
  background: white;
  border: 2px solid #f59e0b;
  border-radius: 8px;
  padding: 0;
  min-width: 160px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: relative;
}

.workflow-node.selected {
  border-color: #7c3aed;
  box-shadow: 0 0 0 2px rgba(124, 58, 237, 0.3);
}

.node-header {
  background: #f59e0b;
  color: white;
  padding: 8px 12px;
  border-radius: 6px 6px 0 0;
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 13px;
}

.node-icon {
  font-size: 16px;
}

.node-content {
  padding: 10px 12px;
}

.node-assignee {
  font-size: 11px;
  color: #92400e;
  background: #fef3c7;
  padding: 2px 6px;
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 6px;
}

.assignee-icon {
  font-size: 10px;
}

.node-instructions {
  font-size: 11px;
  color: #6b7280;
  line-height: 1.4;
}

.handle-input {
  width: 12px !important;
  height: 12px !important;
  background: #6b7280 !important;
  border: 2px solid white !important;
}

.handle-success {
  width: 12px !important;
  height: 12px !important;
  background: #10b981 !important;
  border: 2px solid white !important;
}

.handle-failure {
  width: 12px !important;
  height: 12px !important;
  background: #ef4444 !important;
  border: 2px solid white !important;
}

.handle-labels {
  position: absolute;
  right: 16px;
  top: 0;
  bottom: 0;
  pointer-events: none;
}

.handle-label {
  position: absolute;
  font-size: 10px;
  transform: translateY(-50%);
}

.handle-label.success {
  color: #10b981;
}

.handle-label.failure {
  color: #ef4444;
}
</style>
