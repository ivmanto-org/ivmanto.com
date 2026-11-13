<script setup lang="ts">
// Pillar line icons (DS v1.4 assets/icons/pillar-*.svg), inline so they follow the theme.
// 44px in the pillar index, 14px inside the pillar tag on a ticket; size and stroke width
// come from the stylesheet. Every shape has pathLength 1 for the draw animation.
import type { PillarId } from '~/composables/useCatalogue'

defineProps<{ pillar: PillarId }>()
</script>

<template>
  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
    <template v-if="pillar === 'governance'">
      <circle pathLength="1" cx="12" cy="5.5" r="2.3" />
      <path pathLength="1" d="M12 7.8v4.2M6 16.2V14a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v2.2" />
      <circle pathLength="1" cx="6" cy="18.5" r="2.3" />
      <circle pathLength="1" cx="18" cy="18.5" r="2.3" />
    </template>
    <template v-else-if="pillar === 'risk'">
      <circle pathLength="1" cx="12" cy="12" r="9" />
      <circle pathLength="1" cx="12" cy="12" r="4.5" />
      <path pathLength="1" d="M12 12l6.4-6.4" />
      <path pathLength="1" d="M12 3v2M21 12h-2M12 21v-2M3 12h2" />
      <circle cx="16.6" cy="15.4" r="1" fill="currentColor" stroke="none" />
    </template>
    <template v-else>
      <path pathLength="1" d="M3.5 6l1.8 1.8L8.5 4.5" />
      <path pathLength="1" d="M12 6.2h8.5" />
      <path pathLength="1" d="M3.5 12l1.8 1.8L8.5 10.5" />
      <path pathLength="1" d="M12 12.2h8.5" />
      <path pathLength="1" d="M3.8 17.2h4.4v3.6H3.8z" />
      <path pathLength="1" d="M12 19h5.5" />
    </template>
  </svg>
</template>
