<script setup lang="ts">
import type { Package } from '~/composables/useCatalogue'

// badge: 'site' marks the one overall recommended start, 'pillar' marks each pillar's start (one amber badge per page).
// compact: no ticks list (home page only, DS v1.1 cards.md).
const props = withDefaults(defineProps<{ pkg: Package; badge?: 'site' | 'pillar' | 'none'; compact?: boolean }>(), {
  badge: 'site',
  compact: false,
})
const { pillar } = useCatalogue()
const pillarName = computed(() => pillar(props.pkg.pillar)?.name ?? LINE_NAMES[props.pkg.line])
const recommended = computed(() => {
  if (props.badge === 'site') return Boolean(props.pkg.recommended)
  if (props.badge === 'pillar') return Boolean(props.pkg.pillar_recommended)
  return false
})
</script>

<template>
  <article class="card package" :class="{ 'package--recommended': recommended }">
    <div class="package__top">
      <span class="tag tag--pillar">{{ pillarName }}</span>
      <span v-if="recommended" class="badge-recommended">Recommended start</span>
    </div>
    <h3>{{ pkg.name }}</h3>
    <div>
      <span class="package__price">{{ formatPrice(pkg) }}<small>{{ priceTerms(pkg) }}</small></span>
    </div>
    <div class="package__meta">
      <span>{{ formatDuration(pkg) }}</span>
      <span>{{ pkg.deliverable }}</span>
    </div>
    <ul v-if="!compact" class="ticks">
      <li v-for="item in pkg.includes" :key="item">{{ item }}</li>
    </ul>
    <NuxtLink :to="`/packages/${pkg.id}`" class="btn btn--secondary btn--block">View package</NuxtLink>
  </article>
</template>
