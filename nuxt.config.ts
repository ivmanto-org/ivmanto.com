// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: true,

  modules: ['@nuxtjs/tailwindcss', '@nuxtjs/sitemap', '@nuxt/image'],

  site: {
    url: 'https://ivmanto.com',
  },

  sitemap: {
    sources: ['/api/__sitemap__/blog'],
    exclude: ['/login', '/booking-demo', '/booking', '/booking/cancel', '/order/**', '/imprint'],
  },

  // Order matters: Tailwind (with tokens), then the design-system component and page
  // layers, legacy styles live inside main.css (they need Tailwind's @layer).
  css: ['~/assets/css/main.css', '~/assets/css/ds/components.css', '~/assets/css/ds/pages.css'],

  runtimeConfig: {
    apiBaseUrl: process.env.NUXT_API_BASE_URL || 'https://ivmanto.com',
    public: {
      geminiApiKey: '',
    },
  },

  nitro: {
    devProxy: {
      '/api': { target: 'http://localhost:8080', changeOrigin: true },
    },
    prerender: {
      // Crawl links to discover blog pages from the sitemap source
      crawlLinks: true,
    },
  },

  features: {
    inlineStyles: false
  },

  routeRules: {
    // Global Security Headers
    '/**': {
      headers: {
        'Strict-Transport-Security': 'max-age=31536000; includeSubDomains; preload',
      },
    },
    // Pre-render static marketing pages
    '/': { prerender: true },
    '/about': { prerender: true },
    '/services': { prerender: true },
    '/services/**': { prerender: true },
    '/packages': { prerender: true },
    '/packages/**': { prerender: true },
    '/order/**': { prerender: true },
    '/ai-use': { prerender: true },
    '/imprint': { prerender: true },
    '/booking-demo': { prerender: true },
    '/login': { prerender: true },
    '/privacy-policy': { prerender: true },
    '/terms': { prerender: true },
    // Client-only for dynamic pages
    '/booking': { prerender: true },
    '/booking/cancel': { prerender: true },
    '/assessment': { prerender: true },
    // Blog: pre-rendered at generate time for SEO
    '/blog': { prerender: true },
    '/blog/**': { prerender: true },
  },

  app: {
    head: {
      htmlAttrs: { lang: 'en' },
      charset: 'utf-8',
      viewport: 'width=device-width, initial-scale=1',
      link: [
        { rel: 'icon', href: '/favicon.ico', sizes: '48x48' },
        { rel: 'icon', href: '/favicon.svg', type: 'image/svg+xml' },
        { rel: 'apple-touch-icon', href: '/apple-touch-icon.png' },
        { rel: 'manifest', href: '/site.webmanifest' },
      ],
      meta: [
        { name: 'theme-color', content: '#ffffff', media: '(prefers-color-scheme: light)' },
        { name: 'theme-color', content: '#0a1220', media: '(prefers-color-scheme: dark)' },
      ],
      script: [
        {
          // Set the theme before first paint (no flash). Light by default; the OS
          // preference applies on first visit; the toggle stores 'ivm-theme'.
          key: 'ivm-theme',
          tagPosition: 'head',
          innerHTML:
            "(function(){try{var t=localStorage.getItem('ivm-theme');if(!t){t=matchMedia('(prefers-color-scheme: dark)').matches?'dark':'light'}document.documentElement.setAttribute('data-theme',t)}catch(e){document.documentElement.setAttribute('data-theme','light')}})()",
        },
      ],
    },
  },

  typescript: {
    strict: true,
  },

  hooks: {
    async 'prerender:routes'(ctx)
    {
      const backendUrl =
        process.env.NUXT_API_BASE_URL || 'https://ivmanto.com'
      try
      {
        const articles = await (await fetch(`${backendUrl}/api/articles`)).json()
        for (const article of articles)
        {
          ctx.routes.add(`/blog/${article.slug}`)
        }
        console.log(`[prerender] Discovered ${articles.length} blog routes`)
      } catch (e)
      {
        console.warn('[prerender] Could not fetch blog routes:', e)
      }
    },
  },

  compatibilityDate: '2025-02-28',
})
