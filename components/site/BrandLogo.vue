<script setup lang="ts">
// Inline logo so it follows the theme (ink = --text, tick/chevron = --brand).
// Source: docs/design-system/assets/logo; regenerate there, never hand-edit paths.
import lockup from '~/assets/brand/lockup.svg?raw'
import wordmark from '~/assets/brand/wordmark.svg?raw'

const props = withDefaults(defineProps<{ variant?: 'lockup' | 'wordmark' }>(), { variant: 'lockup' })
const svg = computed(() => (props.variant === 'wordmark' ? wordmark : lockup))
</script>

<template>
  <span class="brand-logo" role="img" aria-label="ivmanto" v-html="svg" />
</template>
