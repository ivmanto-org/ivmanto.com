<script setup lang="ts">
import type { Package } from '~/composables/useCatalogue'

const props = defineProps<{ pkg: Package; showRecommended?: boolean }>()
const { pillar } = useCatalogue()
const pillarName = computed(() => pillar(props.pkg.pillar)?.name ?? '')
const recommended = computed(() => props.showRecommended !== false && props.pkg.recommended)
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
