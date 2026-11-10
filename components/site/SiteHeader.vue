<script setup lang="ts">
const route = useRoute()
const menuOpen = ref(false)

const nav = [
  { to: '/services', label: 'Services' },
  { to: '/packages', label: 'Packages' },
  { to: '/blog', label: 'Insights' },
  { to: '/about', label: 'About' },
]

const isActive = (to: string) => route.path === to || route.path.startsWith(`${to}/`)

watch(() => route.fullPath, () => {
  menuOpen.value = false
})
</script>

<template>
  <header class="site-header">
    <div class="container">
      <NuxtLink to="/" class="brand" aria-label="ivmanto home">
        <SiteBrandLogo variant="lockup" />
      </NuxtLink>
      <nav class="nav" aria-label="Main">
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
</style>
