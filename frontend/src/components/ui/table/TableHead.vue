<script setup lang="ts">
import { cn } from '@/lib/utils'
</script>

<template>
  <th
    :class="cn(
      'h-10 px-4 text-left align-middle font-medium text-zinc-500 [&:has([role=checkbox])]:pr-0 dark:text-zinc-400',
      $attrs.class ?? ''
    )"
  >
    <slot />
  </th>
</template>
