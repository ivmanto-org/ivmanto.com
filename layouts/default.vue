<script setup lang="ts">
// Global Schema (JSON-LD structured data)
useHead({
  script: [
    {
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': ['Organization', 'ProfessionalService'],
        '@id': 'https://ivmanto.com/#organization',
        name: 'ivmanto',
        alternateName: 'ivmanto - AI governance, risk and compliance services',
        url: 'https://ivmanto.com',
        logo: 'https://ivmanto.com/icon-512.png',
        image: 'https://ivmanto.com/og-image.png',
        description:
          'AI governance, risk and compliance services for SMEs: EU AI Act, ISO/IEC 42001, NIST AI RMF. Delivered by Nikolay Tonev.',
        address: {
          '@type': 'PostalAddress',
          addressCountry: 'DE',
        },
        founder: {
          '@id': 'https://ivmanto.com/about#person',
        },
        sameAs: [
          'https://linkedin.com/in/nikolaytonev',
          'https://github.com/dasiyes',
          'https://g.dev/ivmanto-nikolaytonev',
        ],
      }),
    },
    {
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'WebSite',
        '@id': 'https://ivmanto.com/#website',
        name: 'IVMANTO',
        url: 'https://ivmanto.com',
        publisher: {
          '@id': 'https://ivmanto.com/#organization',
        },
      }),
    },
  ],
})
</script>

<template>
  <div class="flex flex-col min-h-screen">
    <a href="#main" class="sr-only focus:not-sr-only">Skip to content</a>
    <SiteHeader />
    <main id="main" class="flex-1">
      <slot />
    </main>
    <SiteFooter />
    <CookieBanner />
  </div>
</template>
