<script setup lang="ts">
import { cn } from '@/lib/utils'
import { DialogDescription } from 'radix-vue'
</script>

<template>
  <DialogDescription
    :class="cn('text-sm text-zinc-500 dark:text-zinc-400', $attrs.class ?? '')"
  >
    <slot />
  </DialogDescription>
</template>
