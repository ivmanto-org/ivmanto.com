<script setup lang="ts">
import type { Package } from '~/composables/useCatalogue'

// badge: 'site' marks the one overall recommended start, 'pillar' marks each pillar's start (one amber badge per page).
// compact: no ticks list on the face (home page only); "What's included" opens it in place
// (DS v1.4 cards.md). The parent owns `open` because the tickets of one row open together.
const props = withDefaults(defineProps<{ pkg: Package; badge?: 'site' | 'pillar' | 'none'; compact?: boolean; open?: boolean }>(), {
  badge: 'site',
  compact: false,
  open: false,
})
defineEmits<{ toggle: [] }>()
const { pillar } = useCatalogue()
const pillarEntry = computed(() => pillar(props.pkg.pillar))
const pillarName = computed(() => pillarEntry.value?.name ?? LINE_NAMES[props.pkg.line])
const moreId = computed(() => `ticket-more-${props.pkg.id}`)
const recommended = computed(() => {
  if (props.badge === 'site') return Boolean(props.pkg.recommended)
  if (props.badge === 'pillar') return Boolean(props.pkg.pillar_recommended)
  return false
})
</script>

<template>
  <!-- Ticket (DS v1.3 cards.md): scope above the perforation, price on the stub. The button is the only link. -->
  <article class="ticket" :class="{ 'ticket--recommended': recommended, 'is-open': compact && open }">
    <div class="ticket__body">
      <div class="package__top">
        <span class="tag tag--pillar"><SitePillarMark v-if="pillarEntry" :pillar="pillarEntry.id" />{{ pillarName }}</span>
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
      <template v-if="compact">
        <button class="ticket__toggle" type="button" :aria-expanded="open" :aria-controls="moreId" @click="$emit('toggle')">
          <span>What’s included</span><span class="ticket__plus" aria-hidden="true"></span>
        </button>
        <div :id="moreId" class="ticket__more" :class="{ 'is-open': open }">
          <div class="ticket__clip">
            <ul class="ticks">
              <li v-for="item in pkg.includes" :key="item">{{ item }}</li>
            </ul>
            <p v-if="pkg.not_included.length" class="ticket__not"><b>Not included</b>{{ pkg.not_included.join('; ') }}</p>
          </div>
        </div>
      </template>
    </div>
    <div class="ticket__stub">
      <div class="ticket__price">{{ formatPrice(pkg) }}<small>{{ priceTerms(pkg) }}</small></div>
      <NuxtLink :to="`/packages/${pkg.id}`" class="btn" :class="recommended ? 'btn--primary' : 'btn--secondary'" :aria-label="`View package: ${pkg.name}`">View <SiteChevIcon /></NuxtLink>
    </div>
  </article>
</template>
