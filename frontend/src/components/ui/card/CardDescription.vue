<script setup lang="ts">
import { cn } from '@/lib/utils'
</script>

<template>
  <p :class="cn('text-sm text-zinc-500 dark:text-zinc-400', $attrs.class ?? '')">
    <slot />
  </p>
</template>
