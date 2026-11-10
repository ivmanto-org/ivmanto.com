<script setup lang="ts">
const { pillars } = useCatalogue()
const year = new Date().getFullYear()

function openCookieSettings() {
  window.dispatchEvent(new Event('ivm:cookie-settings'))
}
</script>

<template>
  <footer class="site-footer">
    <div class="container">
      <div class="cols">
        <div>
          <NuxtLink to="/" class="brand" aria-label="ivmanto home"><SiteBrandLogo variant="wordmark" /></NuxtLink>
          <p class="small" style="margin-top: var(--space-4)">
            AI governance, risk and compliance services for SMEs. Delivered by Nikolay Tonev.
          </p>
        </div>
        <div>
          <h4>Services</h4>
          <ul>
            <li v-for="p in pillars" :key="p.id"><NuxtLink :to="`/services/${p.id}`">{{ p.name }}</NuxtLink></li>
            <li><NuxtLink to="/packages">All packages</NuxtLink></li>
          </ul>
        </div>
        <div>
          <h4>Company</h4>
          <ul>
            <li><NuxtLink to="/about">About</NuxtLink></li>
            <li><NuxtLink to="/blog">Insights</NuxtLink></li>
            <li><NuxtLink to="/booking">Book a call</NuxtLink></li>
            <li><NuxtLink to="/#contact">Contact</NuxtLink></li>
          </ul>
        </div>
        <div>
          <h4>Legal</h4>
          <ul>
            <li><NuxtLink to="/imprint">Imprint</NuxtLink></li>
            <li><NuxtLink to="/privacy-policy">Privacy</NuxtLink></li>
            <li><NuxtLink to="/terms">Terms</NuxtLink></li>
            <li><a href="#" @click.prevent="openCookieSettings">Cookie settings</a></li>
            <li><NuxtLink to="/ai-use">AI use on this site</NuxtLink></li>
          </ul>
        </div>
      </div>
      <div class="legal-line">
        <span>© {{ year }} ivmanto</span>
        <span>Information on this site is not legal advice.</span>
      </div>
    </div>
  </footer>
</template>
