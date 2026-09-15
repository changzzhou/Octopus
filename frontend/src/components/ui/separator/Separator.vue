<script setup lang="ts">
import { cn } from '@/lib/utils'

defineProps<{
  orientation?: 'horizontal' | 'vertical'
}>()
</script>

<template>
  <div
    :class="cn(
      'shrink-0 bg-zinc-200 dark:bg-zinc-800',
      orientation === 'vertical' ? 'h-full w-[1px]' : 'h-[1px] w-full',
      $attrs.class ?? ''
    )"
  />
</template>
