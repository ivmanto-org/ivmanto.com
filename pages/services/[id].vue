<script setup lang="ts">
// Pillar template (docs/design-system/pages/pillar.html). Framework clause references
// are checked by Ivgil before publishing; the checklist copy is a placeholder.
const route = useRoute()
const { pillars, pillar, packagesFor } = useCatalogue()

const current = pillar(String(route.params.id))
if (!current) throw createError({ statusCode: 404, statusMessage: 'Service not found', fatal: true })

const packages = packagesFor(current.id)
const others = pillars.filter((p) => p.id !== current.id)
const frameworks = [...new Set(packages.flatMap((p) => p.frameworks))]

useSeoMeta({ title: current.name, description: current.summary })
</script>

<template>
  <div>
    <section class="pillar-hero hexgrid">
      <div class="container">
        <nav class="crumbs" aria-label="Breadcrumb"><NuxtLink to="/services">Services</NuxtLink> / {{ current.name }}</nav>
        <div class="split">
          <div>
            <SitePillarIcon :pillar="current.id" />
            <p class="overline">{{ current.name }}</p>
            <h1>{{ current.summary }}</h1>
            <div class="frameworks">
              <span v-for="f in frameworks" :key="f" class="tag tag--framework">{{ f }}</span>
            </div>
          </div>
          <div class="card">
            <h3>You likely need this if…</h3>
            <ul class="problems">
              <li>[Situation 1 for {{ current.name }}]</li>
              <li>[Situation 2]</li>
              <li>[Situation 3]</li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <section class="section section--sunken">
      <div class="container">
        <div class="section-head">
          <p class="overline">Packages</p>
          <h2>{{ current.name }} packages</h2>
          <p class="lead">Fixed scope and fixed price. Start with the one closest to your deadline.</p>
        </div>
        <div class="grid grid-3">
          <SitePackageCard v-for="pk in packages" :key="pk.id" :pkg="pk" />
        </div>
      </div>
    </section>

    <section class="section">
      <div class="container">
        <h2>The other pillars</h2>
        <div class="grid grid-2">
          <article v-for="p in others" :key="p.id" class="card card--link pillar">
            <SitePillarIcon :pillar="p.id" />
            <h3><NuxtLink :to="`/services/${p.id}`" class="card__link" style="color: inherit; text-decoration: none">{{ p.name }}</NuxtLink></h3>
            <p class="muted">{{ p.summary }}</p>
          </article>
        </div>
      </div>
    </section>
  </div>
</template>
