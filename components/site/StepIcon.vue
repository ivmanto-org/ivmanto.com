<script setup lang="ts">
// How-it-works step icons (DS v1.1 assets/icons/step-*.svg), inline so they follow the theme.
export type Step = 'describe' | 'choose' | 'deliver' | 'handover'
defineProps<{ step: Step }>()

const PATHS: Record<Step, string[]> = {
  describe: ['M4 5.5h16v11H11l-4.5 3.5v-3.5H4z', 'M8 9.5h8M8 12.5h5'],
  choose: ['M12 3l8 4.5v9L12 21l-8-4.5v-9z', 'M8.5 12.2l2.6 2.6 4.6-5.2'],
  deliver: ['M6 3.5h8l4 4v13H6z', 'M14 3.5v4h4M9 12h6M9 15.5h6'],
  handover: ['M4 14v5.5h16V14', 'M12 4v10.5M8 10.5l4 4 4-4'],
}
</script>

<template>
  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
    <path v-for="d in PATHS[step]" :key="d" :d="d" />
  </svg>
</template>
