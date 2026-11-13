<script setup lang="ts">
// One Data & AI engineering service. Offer tiers render structured content; background
// (expertise) pages render their long-form component. Old prices are not shown: projects
// are scoped on request (Nick, 2026-11-11).
import { getServiceById } from '~/data/services'

const route = useRoute()
const current = getServiceById(String(route.params.id))
if (!current || current.tier === 'entry') throw createError({ statusCode: 404, statusMessage: 'Service not found', fatal: true })

const next = getServiceById(current.nextStepId)
useSeoMeta({ title: current.menuTitle, description: current.seoDescription })
</script>

<template>
  <div>
    <section class="pillar-hero hexgrid">
      <div class="container">
        <nav class="crumbs" aria-label="Breadcrumb"><NuxtLink to="/data-ai">Data &amp; AI engineering</NuxtLink> / {{ current.menuTitle }}</nav>
        <div class="split">
          <div>
            <p class="overline">{{ current.menuTitle }}</p>
            <h1>{{ current.headline ?? current.menuTitle }}</h1>
            <p class="lead">{{ current.intro ?? current.summary }}</p>
            <div v-if="current.stack" class="frameworks">
              <span v-for="t in current.stack" :key="t" class="tag tag--framework">{{ t }}</span>
            </div>
          </div>
          <div v-if="current.problems" class="card">
            <h3>You likely need this if…</h3>
            <ul class="problems">
              <li v-for="p in current.problems" :key="p">{{ p }}</li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <section v-if="current.deliverables" class="section section--sunken">
      <div class="container">
        <div class="section-head"><p class="overline">What you get</p><h2>Deliverables</h2></div>
        <div class="grid grid-2">
          <article v-for="d in current.deliverables" :key="d.title" class="card">
            <h3>{{ d.title }}</h3>
            <p class="muted">{{ d.text }}</p>
          </article>
        </div>
      </div>
    </section>

    <section v-if="current.useCases" class="section">
      <div class="container" style="max-width: 760px">
        <h2>Typical projects</h2>
        <ul class="ticks"><li v-for="u in current.useCases" :key="u">{{ u }}</li></ul>
      </div>
    </section>

    <section v-if="current.detailsComponent" class="section">
      <div class="container prose" style="max-width: 760px">
        <component :is="current.detailsComponent" />
      </div>
    </section>

    <section v-if="current.faqs" class="section section--sunken">
      <div class="container" style="max-width: 760px">
        <h2>Questions</h2>
        <div class="faq">
          <details v-for="f in current.faqs" :key="f.question"><summary>{{ f.question }}</summary><p>{{ f.answer }}</p></details>
        </div>
      </div>
    </section>

    <section class="section">
      <div class="container row">
        <NuxtLink to="/packages/eng-architecture-review" class="btn btn--primary">Start with the architecture review</NuxtLink>
        <NuxtLink to="/booking" class="btn btn--secondary">Book a call</NuxtLink>
        <NuxtLink v-if="next && next.tier !== 'entry'" :to="`/data-ai/${next.id}`" class="btn btn--ghost">Next: {{ next.menuTitle }}</NuxtLink>
      </div>
    </section>
  </div>
</template>
