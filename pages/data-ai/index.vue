<script setup lang="ts">
// Data & AI engineering on Google Cloud (IA proposal 2026-11-11, Nick approved).
// Content comes from data/services.ts (the old site's GCP services); Brio rewrites it in the new voice.
import { services } from '~/data/services'

const { packagesInLine } = useCatalogue()
const entry = packagesInLine('engineering')
const offers = services.filter((s) => s.tier === 'core' || s.tier === 'ongoing')
const background = services.filter((s) => s.tier === 'expertise')

useSeoMeta({
  title: 'Data & AI engineering on Google Cloud',
  description: 'Data platforms, AI agents and private knowledge assistants on Google Cloud, designed and built in your own Google Cloud project.',
})
</script>

<template>
  <div>
    <section class="pillar-hero hexgrid">
      <div class="container">
        <p class="overline">Data &amp; AI engineering</p>
        <h1>We build what we govern: data and AI systems on Google Cloud</h1>
        <p class="lead" style="max-width: 720px">
          Data platforms, AI agents and private knowledge assistants, designed and built in your own Google Cloud
          project. Start with a fixed-price architecture review; everything after that is a project scoped with you.
        </p>
        <div class="row" style="margin-top: var(--space-6)">
          <NuxtLink v-if="entry[0]" :to="`/packages/${entry[0].id}`" class="btn btn--primary">Start with the architecture review</NuxtLink>
          <NuxtLink to="/booking" class="btn btn--secondary">Book a call</NuxtLink>
        </div>
      </div>
    </section>

    <section v-if="entry.length" class="section section--sunken">
      <div class="container">
        <div class="section-head">
          <p class="overline">Fixed price</p>
          <h2>Entry package</h2>
        </div>
        <div class="grid grid-3">
          <SitePackageCard v-for="pk in entry" :key="pk.id" :pkg="pk" badge="none" />
        </div>
      </div>
    </section>

    <section class="section">
      <div class="container">
        <div class="section-head">
          <p class="overline">Projects on request</p>
          <h2>What we build</h2>
          <p class="lead">Scope and price are agreed with you after the review or a first call.</p>
        </div>
        <div class="grid grid-2">
          <article v-for="s in offers" :key="s.id" class="card card--link pillar">
            <h3><NuxtLink :to="`/data-ai/${s.id}`" class="card__link" style="color: inherit; text-decoration: none">{{ s.menuTitle }}</NuxtLink></h3>
            <p class="muted">{{ s.summary }}</p>
          </article>
        </div>
      </div>
    </section>

    <section class="section section--sunken">
      <div class="container">
        <div class="section-head">
          <p class="overline">Background</p>
          <h2>How we think about data and AI</h2>
        </div>
        <div class="grid grid-3">
          <article v-for="s in background" :key="s.id" class="card card--link">
            <h3><NuxtLink :to="`/data-ai/${s.id}`" class="card__link" style="color: inherit; text-decoration: none">{{ s.menuTitle }}</NuxtLink></h3>
            <p class="muted small">{{ s.summary }}</p>
          </article>
        </div>
      </div>
    </section>

    <section class="section">
      <div class="container">
        <div class="card" style="max-width: 760px">
          <h3>Governance comes with it</h3>
          <p class="muted">Every system we build ships with an owner, a risk record and documentation your auditors can read. If you need the governance side first, start with a package.</p>
          <NuxtLink to="/packages" class="btn btn--ghost">See governance packages</NuxtLink>
        </div>
      </div>
    </section>
  </div>
</template>
