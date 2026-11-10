<script setup lang="ts">
const { pillars, packagesFor } = useCatalogue()
useSeoMeta({ title: 'Packages and prices', description: 'Fixed-scope, fixed-price AI governance, risk and compliance packages.' })
</script>

<template>
  <div>
    <section class="section hexgrid" style="padding-bottom: var(--space-8)">
      <div class="container">
        <div class="section-head">
          <p class="overline">Packages</p>
          <h1>Fixed scope. Fixed price. Order directly.</h1>
          <p class="lead">No open-ended retainers. Each package states what you get, how long it takes and what it costs.</p>
        </div>
      </div>
    </section>
    <section v-for="(p, i) in pillars" :key="p.id" class="section" :class="{ 'section--sunken': i % 2 === 0 }">
      <div class="container">
        <div class="section-head">
          <h2><NuxtLink :to="`/services/${p.id}`" style="color: inherit; text-decoration: none">{{ p.name }}</NuxtLink></h2>
          <p class="muted">{{ p.summary }}</p>
        </div>
        <div class="grid grid-3">
          <SitePackageCard v-for="pk in packagesFor(p.id)" :key="pk.id" :pkg="pk" />
        </div>
      </div>
    </section>
  </div>
</template>
