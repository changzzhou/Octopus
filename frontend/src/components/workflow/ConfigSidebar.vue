<script setup lang="ts">
import { computed } from 'vue'
import type { WorkflowNode, ScriptNodeConfig, HttpNodeConfig, HumanNodeConfig, NodeType, WorkflowNodeData } from '../../types/workflow'
import { NODE_TYPE_CONFIGS } from '../../types/workflow'

const props = defineProps<{
  node: WorkflowNode | null
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
  <div class="config-sidebar" :class="{ open: !!node }">
    <div v-if="node && nodeConfig" class="sidebar-content">
      <div class="sidebar-header" :style="{ backgroundColor: nodeConfig.color }">
        <div class="header-title">
          <span class="header-icon">{{ nodeConfig.icon }}</span>
          <span>{{ nodeConfig.label }}</span>
        </div>
        <button class="close-btn" @click="$emit('close')">×</button>
      </div>
      
      <div v-if="node.data" class="sidebar-body">
        <div class="form-group">
          <label class="form-label">Label</label>
          <input
            type="text"
            class="form-input"
            :value="node.data.label"
            @input="updateLabel(($event.target as HTMLInputElement).value)"
            placeholder="Node label"
          />
        </div>
        
        <template v-if="node.data.type === 'script'">
          <div class="form-group">
            <label class="form-label">Language</label>
            <select
              class="form-select"
              :value="(node.data.config as ScriptNodeConfig).language"
              @change="updateConfig('language', ($event.target as HTMLSelectElement).value)"
            >
              <option value="javascript">JavaScript</option>
              <option value="python">Python</option>
              <option value="shell">Shell</option>
            </select>
          </div>
          
          <div class="form-group">
            <label class="form-label">Code</label>
            <textarea
              class="form-textarea"
              :value="(node.data.config as ScriptNodeConfig).code"
              @input="updateConfig('code', ($event.target as HTMLTextAreaElement).value)"
              rows="8"
              placeholder="Enter your code here..."
            ></textarea>
          </div>
          
          <div class="form-group">
            <label class="form-label">Timeout (ms)</label>
            <input
              type="number"
              class="form-input"
              :value="(node.data.config as ScriptNodeConfig).timeout"
              @input="updateConfig('timeout', Number(($event.target as HTMLInputElement).value))"
              min="0"
            />
          </div>
        </template>
        
        <template v-else-if="node.data.type === 'http'">
          <div class="form-group">
            <label class="form-label">Method</label>
            <select
              class="form-select"
              :value="(node.data.config as HttpNodeConfig).method"
              @change="updateConfig('method', ($event.target as HTMLSelectElement).value)"
            >
              <option value="GET">GET</option>
              <option value="POST">POST</option>
              <option value="PUT">PUT</option>
              <option value="DELETE">DELETE</option>
              <option value="PATCH">PATCH</option>
            </select>
          </div>
          
          <div class="form-group">
            <label class="form-label">URL</label>
            <input
              type="text"
              class="form-input"
              :value="(node.data.config as HttpNodeConfig).url"
              @input="updateConfig('url', ($event.target as HTMLInputElement).value)"
              placeholder="https://api.example.com/endpoint"
            />
          </div>
          
          <div class="form-group">
            <label class="form-label">Headers (JSON)</label>
            <textarea
              class="form-textarea"
              :value="JSON.stringify((node.data.config as HttpNodeConfig).headers || {}, null, 2)"
              @input="handleHeadersInput(($event.target as HTMLTextAreaElement).value)"
              rows="4"
              placeholder='{"Authorization": "Bearer ..."}'
            ></textarea>
          </div>
          
          <div class="form-group">
            <label class="form-label">Body</label>
            <textarea
              class="form-textarea"
              :value="(node.data.config as HttpNodeConfig).body"
              @input="updateConfig('body', ($event.target as HTMLTextAreaElement).value)"
              rows="4"
              placeholder="Request body (JSON, text, etc.)"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label class="form-label">Timeout (ms)</label>
            <input
              type="number"
              class="form-input"
              :value="(node.data.config as HttpNodeConfig).timeout"
              @input="updateConfig('timeout', Number(($event.target as HTMLInputElement).value))"
              min="0"
            />
          </div>
        </template>
        
        <template v-else-if="node.data.type === 'human'">
          <div class="form-group">
            <label class="form-label">Assignee</label>
            <input
              type="text"
              class="form-input"
              :value="(node.data.config as HumanNodeConfig).assignee"
              @input="updateConfig('assignee', ($event.target as HTMLInputElement).value)"
              placeholder="user@example.com"
            />
          </div>
          
          <div class="form-group">
            <label class="form-label">Instructions</label>
            <textarea
              class="form-textarea"
              :value="(node.data.config as HumanNodeConfig).instructions"
              @input="updateConfig('instructions', ($event.target as HTMLTextAreaElement).value)"
              rows="4"
              placeholder="Enter instructions for the human task..."
            ></textarea>
          </div>
        </template>
        
        <div class="form-actions">
          <button class="btn btn-danger" @click="deleteNode">
            Delete Node
          </button>
        </div>
      </div>
    </div>
    
    <div v-else class="sidebar-empty">
      <div class="empty-icon">👆</div>
      <p>Select a node to configure</p>
    </div>
  </div>
</template>

<style scoped>
.config-sidebar {
  width: 280px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  max-height: 100%;
}

.sidebar-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.sidebar-header {
  padding: 12px 16px;
  color: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.header-icon {
  font-size: 18px;
}

.close-btn {
  width: 24px;
  height: 24px;
  border: none;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 18px;
  line-height: 1;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.sidebar-body {
  padding: 16px;
  overflow-y: auto;
  flex: 1;
}

.form-group {
  margin-bottom: 16px;
}

.form-label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 6px;
}

.form-input,
.form-select,
.form-textarea {
  width: 100%;
  padding: 8px 10px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 13px;
  transition: border-color 0.2s;
}

.form-input:focus,
.form-select:focus,
.form-textarea:focus {
  outline: none;
  border-color: #7c3aed;
  box-shadow: 0 0 0 2px rgba(124, 58, 237, 0.1);
}

.form-textarea {
  resize: vertical;
  font-family: monospace;
  font-size: 12px;
}

.form-actions {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
}

.btn {
  width: 100%;
  padding: 10px 16px;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  transition: all 0.2s;
}

.btn-danger {
  background: #fef2f2;
  color: #dc2626;
}

.btn-danger:hover {
  background: #fee2e2;
}

.sidebar-empty {
  padding: 40px 20px;
  text-align: center;
  color: #9ca3af;
}

.empty-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.sidebar-empty p {
  margin: 0;
  font-size: 14px;
}
</style>
