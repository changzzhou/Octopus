<script setup lang="ts">
import { SelectValue } from 'radix-vue'

defineProps<{
  placeholder?: string
}>()
</script>

<template>
  <SelectValue :placeholder="placeholder" />
</template>
