<script setup lang="ts">
const route = useRoute()

const siteTitle = 'ivmanto | AI governance, risk and compliance'
const siteDescription =
  'Fixed-scope services that get your company ready for the EU AI Act and ISO/IEC 42001: clear owners, documented risks, audit-ready evidence.'
const siteUrl = 'https://ivmanto.com'
const ogImage = `${siteUrl}/og-image.png`

useHead({
  titleTemplate: (titleChunk) => {
    if (!titleChunk) return siteTitle;
    return titleChunk.toLowerCase().includes('ivmanto') ? titleChunk : `${titleChunk} | ivmanto.com`;
  },
  link: [
    {
      rel: 'canonical',
      href: computed(() => {
        let path = route.path === '/' ? '' : route.path
        if (path.endsWith('/') && path.length > 1) path = path.slice(0, -1)
        return `${siteUrl}${path}`
      }),
    },
  ],
})

useSeoMeta({
  title: siteTitle,
  description: siteDescription,
  author: 'Nikolay Tonev',
  ogType: 'website',
  ogSiteName: 'ivmanto.com',
  ogLocale: 'en_US',
  ogTitle: siteTitle,
  ogDescription: siteDescription,
  ogUrl: computed(() => {
    let path = route.path === '/' ? '' : route.path
    if (path.endsWith('/') && path.length > 1) path = path.slice(0, -1)
    return `${siteUrl}${path}`
  }),
  ogImage: ogImage,
  twitterCard: 'summary_large_image',
  twitterTitle: siteTitle,
  twitterDescription: siteDescription,
  twitterImage: ogImage,
})
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
