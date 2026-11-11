<script setup lang="ts">
const route = useRoute()
const menuOpen = ref(false)

const govOpen = ref(false)
const { pillars } = useCatalogue()

// Governance leads; Data & AI engineering and ivGIL support it (Nick, 2026-11-11).
const govLinks = [...pillars.map((p) => ({ to: `/services/${p.id}`, label: p.name })), { to: '/packages', label: 'All packages' }]
const nav = [
  { to: '/data-ai', label: 'Data & AI engineering' },
  { to: '/ivgil', label: 'ivGIL' },
  { to: '/blog', label: 'Insights' },
  { to: '/about', label: 'About' },
]
const govActive = computed(() => route.path.startsWith('/services') || route.path.startsWith('/packages'))

const isActive = (to: string) => route.path === to || route.path.startsWith(`${to}/`)

watch(() => route.fullPath, () => {
  menuOpen.value = false
  govOpen.value = false
})
</script>

<template>
  <header class="site-header">
    <div class="container">
      <NuxtLink to="/" class="brand" aria-label="ivmanto home">
        <SiteBrandLogo variant="lockup" />
      </NuxtLink>
      <nav class="nav" aria-label="Main">
        <div class="dropdown" @mouseleave="govOpen = false" @keydown.esc="govOpen = false">
          <button
            type="button"
            class="dropdown__toggle"
            :class="{ 'is-active': govActive }"
            :aria-expanded="govOpen"
            aria-controls="gov-menu"
            @click="govOpen = !govOpen"
          >
            Governance
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" aria-hidden="true"><path d="M6 9l6 6 6-6" /></svg>
          </button>
          <div v-show="govOpen" id="gov-menu" class="dropdown__panel">
            <NuxtLink v-for="g in govLinks" :key="g.to" :to="g.to" :aria-current="isActive(g.to) ? 'page' : undefined">{{ g.label }}</NuxtLink>
          </div>
        </div>
        <NuxtLink v-for="item in nav" :key="item.to" :to="item.to" :aria-current="isActive(item.to) ? 'page' : undefined">
          {{ item.label }}
        </NuxtLink>
      </nav>
      <div class="header-actions">
        <SiteThemeToggle />
        <NuxtLink to="/booking" class="btn btn--primary btn--sm">Book a call</NuxtLink>
        <button
          type="button"
          class="icon-btn menu-btn"
          :aria-expanded="menuOpen"
          aria-controls="mobile-menu"
          aria-label="Menu"
          @click="menuOpen = !menuOpen"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
            <path v-if="!menuOpen" d="M4 7h16M4 12h16M4 17h16" />
            <path v-else d="M6 6l12 12M18 6L6 18" />
          </svg>
        </button>
      </div>
    </div>
    <nav v-if="menuOpen" id="mobile-menu" class="mobile-menu" aria-label="Main">
      <div class="container">
        <p class="mobile-menu__group">Governance</p>
        <NuxtLink v-for="g in govLinks" :key="g.to" :to="g.to" class="mobile-menu__sub" :aria-current="isActive(g.to) ? 'page' : undefined">
          {{ g.label }}
        </NuxtLink>
        <NuxtLink v-for="item in nav" :key="item.to" :to="item.to" :aria-current="isActive(item.to) ? 'page' : undefined">
          {{ item.label }}
        </NuxtLink>
        <NuxtLink to="/booking" class="btn btn--primary btn--block">Book a call</NuxtLink>
      </div>
    </nav>
  </header>
</template>

<style scoped>
.mobile-menu { border-top: 1px solid var(--border); padding: var(--space-4) 0 var(--space-6); }
.mobile-menu .container { display: flex; flex-direction: column; gap: var(--space-2); height: auto; align-items: stretch; }
.mobile-menu a:not(.btn) { color: var(--text); text-decoration: none; padding: 10px 0; border-bottom: 1px solid var(--border); }
.mobile-menu a[aria-current='page'] { color: var(--accent); }
.mobile-menu .btn { margin-top: var(--space-3); }
@media (min-width: 960px) { .mobile-menu { display: none; } }
.mobile-menu__group { margin: var(--space-2) 0 0; font-size: var(--text-sm); color: var(--text-muted); }
.mobile-menu a.mobile-menu__sub { padding-left: var(--space-4); }
.dropdown { position: relative; display: inline-flex; }
.dropdown__toggle { display: inline-flex; align-items: center; gap: 6px; background: none; border: 0; border-bottom: 2px solid transparent; padding: 6px 0; font: inherit; font-size: var(--text-sm); font-weight: 500; color: var(--text-muted); cursor: pointer; }
.dropdown__toggle:hover { color: var(--text); }
.dropdown__toggle.is-active { color: var(--text); border-bottom-color: var(--accent); }
.dropdown__panel { position: absolute; top: calc(100% + 10px); left: -12px; min-width: 220px; display: flex; flex-direction: column; padding: var(--space-2); background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-md, 8px); box-shadow: 0 8px 24px rgb(0 0 0 / 0.12); z-index: 50; }
.dropdown__panel::before { content: ''; position: absolute; top: -12px; left: 0; right: 0; height: 12px; }
.dropdown__panel a { border-bottom: 0; padding: 8px 12px; border-radius: 6px; color: var(--text); text-decoration: none; white-space: nowrap; }
.dropdown__panel a:hover, .dropdown__panel a[aria-current='page'] { background: var(--surface-sunken); color: var(--accent); }
</style>
