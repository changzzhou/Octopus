<script setup lang="ts">
import { cn } from '@/lib/utils'
import {
  SelectItem,
  SelectItemIndicator,
  SelectItemText,
} from 'radix-vue'
import { Check } from '@lucide/vue'

defineProps<{
  value: string
  disabled?: boolean
}>()
</script>

<template>
  <SelectItem
    :value="value"
    :disabled="disabled"
    :class="cn(
      'relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-2 pr-8 text-sm outline-none focus:bg-zinc-100 focus:text-zinc-900 data-[disabled]:pointer-events-none data-[disabled]:opacity-50 dark:focus:bg-zinc-800 dark:focus:text-zinc-50',
      $attrs.class ?? ''
    )"
  >
    <span class="absolute right-2 flex h-3.5 w-3.5 items-center justify-center">
      <SelectItemIndicator>
        <Check class="h-4 w-4" />
      </SelectItemIndicator>
    </span>
    <SelectItemText>
      <slot />
    </SelectItemText>
  </SelectItem>
</template>
