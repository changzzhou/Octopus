<script setup lang="ts">
import { cn } from '@/lib/utils'
</script>

<template>
  <tr
    :class="cn(
      'border-b border-zinc-200 transition-colors hover:bg-zinc-100/50 data-[state=selected]:bg-zinc-100 dark:border-zinc-800 dark:hover:bg-zinc-800/50 dark:data-[state=selected]:bg-zinc-800',
      $attrs.class ?? ''
    )"
  >
    <slot />
  </tr>
</template>
