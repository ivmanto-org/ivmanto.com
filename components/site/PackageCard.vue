<script setup lang="ts">
import type { Package } from '~/composables/useCatalogue'

// badge: 'site' marks the one overall recommended start, 'pillar' marks each pillar's start (one amber badge per page).
const props = withDefaults(defineProps<{ pkg: Package; badge?: 'site' | 'pillar' | 'none' }>(), { badge: 'site' })
const { pillar } = useCatalogue()
const pillarName = computed(() => pillar(props.pkg.pillar)?.name ?? '')
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
      <span class="package__price">{{ formatPrice(pkg) }}<small>fixed, excl. VAT</small></span>
    </div>
    <div class="package__meta">
      <span>{{ pkg.duration_weeks }} weeks</span>
      <span>{{ pkg.deliverable }}</span>
    </div>
    <ul class="ticks">
      <li v-for="item in pkg.includes" :key="item">{{ item }}</li>
    </ul>
    <NuxtLink :to="`/packages/${pkg.id}`" class="btn btn--secondary btn--block">View package</NuxtLink>
  </article>
</template>
