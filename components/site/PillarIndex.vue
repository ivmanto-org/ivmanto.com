<script setup lang="ts">
// Home pillar index that opens in place (DS v1.4 components/cards.md). Flat markup: head,
// panel, head, panel, head, panel. One panel is open at a time; the pillar page link and
// the pillar's packages are inside the panel.
import type { PillarId } from '~/composables/useCatalogue'

const { pillars, packagesFor } = useCatalogue()
const openId = ref<PillarId | null>(null)

function toggle(id: PillarId) {
  openId.value = openId.value === id ? null : id
}

const ARROW = 'M5 12h14M13 6l6 6-6 6'
</script>

<template>
  <div class="pillar-index">
    <template v-for="(p, i) in pillars" :key="p.id">
      <div class="pillar-index__head" :class="{ 'is-open': openId === p.id }" :style="{ '--i': i }">
        <span class="pillar-index__no" aria-hidden="true">{{ String(i + 1).padStart(2, '0') }}</span>
        <span class="pillar-index__icon"><SitePillarMark :pillar="p.id" /></span>
        <h3 :id="`pillar-h-${p.id}`">
          <button class="pillar-index__btn" type="button" :aria-expanded="openId === p.id" :aria-controls="`pillar-p-${p.id}`" @click="toggle(p.id)">{{ p.name }}</button>
        </h3>
        <p class="muted">{{ p.tagline }}</p>
        <span class="pillar-index__more" aria-hidden="true">
          <span class="t-show">Packages and details</span><span class="t-hide">Close</span>
          <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M3 6l5 5 5-5" /></svg>
        </span>
      </div>
      <div :id="`pillar-p-${p.id}`" class="pillar-index__panel" :class="{ 'is-open': openId === p.id }" role="region" :aria-labelledby="`pillar-h-${p.id}`">
        <div class="pillar-index__clip">
          <div class="pillar-index__card" :style="{ '--i': i }">
            <div class="pillar-index__side">
              <p class="pillar-index__sum">{{ p.summary }}</p>
              <NuxtLink :to="`/services/${p.id}`" class="btn btn--ghost">More about {{ p.name.toLowerCase() }} <SiteChevIcon /></NuxtLink>
            </div>
            <div>
              <p class="pillar-index__label">{{ packagesFor(p.id).length === 1 ? 'Package' : 'Packages' }}</p>
              <ul class="pillar-index__list">
                <li v-for="pk in packagesFor(p.id)" :key="pk.id">
                  <NuxtLink :to="`/packages/${pk.id}`">{{ pk.name }}</NuxtLink>
                  <small>{{ formatDuration(pk) }} · {{ pk.deliverable }}</small>
                  <span class="go" aria-hidden="true">
                    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path :d="ARROW" /></svg>
                  </span>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
