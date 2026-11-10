<script setup lang="ts">
const props = defineProps<{ pillar: string; size?: number }>()
const letter = computed(() => ({ governance: 'G', risk: 'R', compliance: 'C' })[props.pillar] ?? '•')
</script>

<template>
  <svg class="card__icon" viewBox="0 0 44 44" :width="size ?? 44" :height="size ?? 44" aria-hidden="true">
    <rect x="3" y="3" width="38" height="38" rx="9" fill="var(--accent-subtle)" stroke="var(--accent)" stroke-width="1.5" />
    <text x="22" y="27.5" text-anchor="middle" font-family="Montserrat, Inter, sans-serif" font-size="15" font-weight="600" fill="var(--accent)">{{ letter }}</text>
  </svg>
</template>
