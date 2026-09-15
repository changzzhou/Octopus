<script setup lang="ts">
import { cn } from '@/lib/utils'
import { DialogTitle } from 'radix-vue'
</script>

<template>
  <DialogTitle
    :class="cn('text-lg font-semibold leading-none tracking-tight', $attrs.class ?? '')"
  >
    <slot />
  </DialogTitle>
</template>
