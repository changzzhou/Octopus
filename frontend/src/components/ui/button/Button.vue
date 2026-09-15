<script setup lang="ts">
import { type HTMLAttributes, computed } from 'vue'
import { cn } from '@/lib/utils'
import { buttonVariants, type ButtonVariants } from './variants'

interface Props extends /* @vue-ignore */ HTMLAttributes {
  variant?: ButtonVariants['variant']
  size?: ButtonVariants['size']
  as?: string
}

const props = withDefaults(defineProps<Props>(), {
  as: 'button',
})

const delegatedProps = computed(() => {
  const { as: _as, variant: _variant, size: _size, class: _class, ...rest } = props
  return rest
})
</script>

<template>
  <component
    :is="as"
    :class="cn(buttonVariants({ variant, size }), $attrs.class ?? '')"
    v-bind="delegatedProps"
  >
    <slot />
  </component>
</template>
