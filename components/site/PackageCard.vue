<script setup lang="ts">
import type { Package } from '~/composables/useCatalogue'

// badge: 'site' marks the one overall recommended start, 'pillar' marks each pillar's start (one amber badge per page).
// compact: no ticks list (home page only, DS v1.3 cards.md).
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
  <!-- Ticket (DS v1.3 cards.md): scope above the perforation, price on the stub. The button is the only link. -->
  <article class="ticket" :class="{ 'ticket--recommended': recommended }">
    <div class="ticket__body">
      <div class="package__top">
        <span class="tag tag--pillar">{{ pillarName }}</span>
        <span v-if="recommended" class="badge-recommended">Recommended start</span>
      </div>
      <h3>{{ pkg.name }}</h3>
      <ul v-if="!compact" class="ticks">
        <li v-for="item in pkg.includes" :key="item">{{ item }}</li>
      </ul>
      <dl class="ticket__spec">
        <div><dt>Duration</dt><dd>{{ formatDuration(pkg) }}</dd></div>
        <div v-if="pkg.deliverable"><dt>You get</dt><dd>{{ pkg.deliverable }}</dd></div>
      </dl>
    </div>
    <div class="ticket__stub">
      <div class="ticket__price">{{ formatPrice(pkg) }}<small>{{ priceTerms(pkg) }}</small></div>
      <NuxtLink :to="`/packages/${pkg.id}`" class="btn" :class="recommended ? 'btn--primary' : 'btn--secondary'" :aria-label="`View package: ${pkg.name}`">View <SiteChevIcon /></NuxtLink>
    </div>
  </article>
</template>
