<script setup lang="ts">
// Package template (docs/design-system/pages/package.html): everything needed to decide.
const route = useRoute()
const { pkg, pillar } = useCatalogue()

const current = pkg(String(route.params.id))
if (!current) throw createError({ statusCode: 404, statusMessage: 'Package not found', fatal: true })
const pl = pillar(current.pillar)!

// Even week-by-week split of the deliverables until Nick provides real timelines.
const timeline = current.includes.map((item, i) => {
  const per = Math.max(1, Math.round(current.duration_weeks / current.includes.length))
  const start = i * per + 1
  const end = i === current.includes.length - 1 ? current.duration_weeks : Math.min(current.duration_weeks, start + per - 1)
  return { weeks: start === end ? `Week ${start}` : `Weeks ${start} to ${end}`, item }
})

useSeoMeta({ title: current.name, description: current.outcome })
</script>

<template>
  <div class="container">
    <nav class="crumbs" aria-label="Breadcrumb">
      <NuxtLink to="/packages">Packages</NuxtLink> / <NuxtLink :to="`/services/${pl.id}`">{{ pl.name }}</NuxtLink> / {{ current.name }}
    </nav>
    <div class="pkg">
      <div>
        <div class="row" style="margin-bottom: var(--space-4)">
          <span class="tag tag--pillar">{{ pl.name }}</span>
          <span v-for="f in current.frameworks" :key="f" class="tag tag--framework">{{ f }}</span>
          <span v-if="current.recommended" class="badge-recommended">Recommended start</span>
        </div>
        <h1>{{ current.name }}</h1>
        <p class="lead">{{ current.outcome }}</p>

        <dl class="facts">
          <div class="fact"><dt>Duration</dt><dd>{{ current.duration_weeks }} weeks</dd></div>
          <div class="fact"><dt>You get</dt><dd>{{ current.deliverable }}</dd></div>
          <div class="fact"><dt>Your effort</dt><dd>[~N hours]</dd></div>
          <div class="fact"><dt>Format</dt><dd>Remote or on site, agreed with you</dd></div>
        </dl>

        <h2>What you get</h2>
        <ul class="ticks">
          <li v-for="item in current.includes" :key="item">{{ item }}</li>
        </ul>

        <h2 style="margin-top: var(--space-12)">Timeline</h2>
        <ol class="timeline">
          <li v-for="t in timeline" :key="t.item"><span class="wk">{{ t.weeks }}</span><span>{{ t.item }}</span></li>
        </ol>

        <h2 style="margin-top: var(--space-12)">Not included</h2>
        <ul class="ticks notin">
          <li v-for="item in current.not_included" :key="item">{{ item }}</li>
        </ul>
      </div>

      <aside class="buybox card" aria-label="Order">
        <span class="package__price">{{ formatPrice(current) }}<small>fixed, excl. VAT</small></span>
        <ul class="ticks" style="margin: var(--space-6) 0">
          <li>Scope and price fixed before we start</li>
          <li>Delivered by Nick, no hand-offs</li>
          <li>You keep every artefact</li>
        </ul>
        <div class="stack">
          <NuxtLink :to="`/order/${current.id}`" class="btn btn--primary btn--block">Order package</NuxtLink>
          <NuxtLink to="/booking" class="btn btn--secondary btn--block">Book a call first</NuxtLink>
        </div>
        <div style="margin-top: var(--space-6)"><SiteCredentialBadge /></div>
      </aside>
    </div>
  </div>
</template>
