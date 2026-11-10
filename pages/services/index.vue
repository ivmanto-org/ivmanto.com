<script setup lang="ts">
const { pillars, packagesFor } = useCatalogue()
useSeoMeta({ title: 'Services', description: 'AI governance, risk management and compliance services for SMEs.' })
</script>

<template>
  <section class="section hexgrid">
    <div class="container">
      <div class="section-head">
        <p class="overline">Services</p>
        <h1>Three pillars, one accountable partner</h1>
        <p class="lead">Most SMEs need a little of each. Start where the pressure is.</p>
      </div>
      <div class="grid grid-3">
        <article v-for="p in pillars" :key="p.id" class="card card--link pillar">
          <SitePillarIcon :pillar="p.id" />
          <h3><NuxtLink :to="`/services/${p.id}`" class="card__link" style="color: inherit; text-decoration: none">{{ p.name }}</NuxtLink></h3>
          <p class="muted">{{ p.summary }}</p>
          <ul>
            <li v-for="pk in packagesFor(p.id)" :key="pk.id">{{ pk.name }}</li>
          </ul>
        </article>
      </div>
    </div>
  </section>
</template>
