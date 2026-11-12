<script setup lang="ts">
// Home: the funnel (docs/design-system/pages/README.md, section order and why).
// Copy in [brackets] and '€ [price]' are placeholders for Nick.
const { pillars, featured } = useCatalogue()
const { fetchArticles, sortedArticles } = useArticles()

await useAsyncData('home-articles', () => fetchArticles())
const insights = computed(() => sortedArticles.value.slice(0, 3))

const frameworks = ['EU AI Act', 'ISO/IEC 42001', 'ISO/IEC 27001', 'NIST AI RMF', 'GDPR']

// Copy: DS v1.1 pages/home-copy.md.
const steps = [
  { icon: 'describe', title: 'Describe', text: 'Ask the advisor or book a call.' },
  { icon: 'choose', title: 'Choose', text: 'Fixed scope and price, up front.' },
  { icon: 'deliver', title: 'Deliver', text: 'Workshops and documents, built with you.' },
  { icon: 'handover', title: 'Hand over', text: 'Everything stays in your tools.' },
] as const

const faqs = [
  {
    q: 'Is this legal advice?',
    a: 'No. We provide governance, risk and compliance consulting. Where legal interpretation is needed, we work alongside your lawyer.',
  },
  { q: 'We’re small. Does the EU AI Act even apply to us?', a: '[Short, accurate answer, checked by Ivgil.]' },
  { q: 'What happens after I order a package?', a: 'You send an order request. Nick replies within one working day to confirm the scope and start date, schedules the kickoff call and sends the invoice.' },
  { q: 'What does the AI advisor do with my message?', a: 'It sends your message to our AI provider, Mistral AI (EU), only to generate the answer. We don’t store the conversation. Details are on the AI use page and in the privacy policy.' },
]

useSeoMeta({
  title: 'ivmanto | AI governance, risk and compliance',
  description:
    'Fixed-scope services that get your company ready for the EU AI Act and ISO/IEC 42001: clear owners, documented risks, audit-ready evidence.',
})

function formatDate(iso: string): string {
  const d = new Date(iso)
  return Number.isNaN(d.getTime()) ? '' : d.toLocaleDateString('en-GB', { year: 'numeric', month: 'short' })
}
</script>

<template>
  <div>
    <!-- 1 · Hero + advisor -->
    <section class="hero hexgrid hero-motion">
      <SiteHeroMotion />
      <div class="container hero-grid">
        <div class="hero-text">
          <p class="overline">AI governance · risk · compliance</p>
          <h1><span class="line"><span>Use AI with confidence.</span></span> <span class="line"><em>Govern it with clarity.</em></span></h1>
          <p class="lead">Fixed-scope services for the EU AI Act and ISO/IEC 42001.</p>
        </div>
        <div class="hero-adv">
          <AdvisorChat />
        </div>
        <div class="hero-more">
          <div class="row">
            <NuxtLink to="/packages" class="btn btn--primary">See packages and prices</NuxtLink>
            <NuxtLink to="/booking" class="btn btn--secondary">Book a call</NuxtLink>
          </div>
          <div class="trust-line">
            <span>MIT Professional Education certificate</span>
            <span>Fixed price, fixed scope</span>
            <span>Based in Germany</span>
          </div>
        </div>
      </div>
    </section>

    <!-- 2 · Frameworks -->
    <section class="section--line" style="padding: var(--space-8) 0">
      <div class="container fw-strip">
        <span class="small muted">Frameworks we work with</span>
        <span v-for="f in frameworks" :key="f" class="tag tag--framework">{{ f }}</span>
      </div>
    </section>

    <!-- 3 · Pillars -->
    <section class="section section--sunken">
      <div class="container">
        <div class="section-head">
          <p class="overline">What we do</p>
          <h2>Three pillars, one accountable partner</h2>
        </div>
        <div class="grid grid-3">
          <article v-for="p in pillars" :key="p.id" class="card card--link pillar">
            <SitePillarIcon :pillar="p.id" />
            <h3><NuxtLink :to="`/services/${p.id}`" class="card__link" style="color: inherit; text-decoration: none">{{ p.name }}</NuxtLink></h3>
            <p class="muted" style="margin: 0">{{ p.tagline }}</p>
          </article>
        </div>
      </div>
    </section>

    <!-- 4 · Packages -->
    <section class="section">
      <div class="container">
        <div class="section-head">
          <p class="overline">Packages</p>
          <h2>Fixed scope. Fixed price. Order directly.</h2>
        </div>
        <div class="grid grid-3">
          <SitePackageCard v-for="pk in featured()" :key="pk.id" :pkg="pk" compact />
        </div>
        <p style="margin-top: var(--space-8)"><NuxtLink to="/packages" class="btn btn--ghost">All packages →</NuxtLink></p>
      </div>
    </section>

    <!-- 4b + 5 · One dark band (DS v1.2 components/backgrounds.md) -->
    <div class="band">
      <section class="section">
        <div class="container">
          <div class="section-head">
            <p class="overline">Built on real engineering</p>
            <h2>We build what we govern</h2>
          </div>
          <div class="grid grid-2">
            <article class="card card--link">
              <h3><NuxtLink to="/data-ai" class="card__link" style="color: inherit; text-decoration: none">Data &amp; AI engineering on Google Cloud</NuxtLink></h3>
              <p class="muted" style="margin: 0">Data platforms and AI agents, built in your own Google Cloud project.</p>
            </article>
            <article class="card card--link">
              <h3><NuxtLink to="/ivgil" class="card__link" style="color: inherit; text-decoration: none">ivGIL</NuxtLink></h3>
              <p class="muted" style="margin: 0">Our AI governance platform. Early access is open.</p>
            </article>
          </div>
        </div>
      </section>
      <section class="section">
        <div class="container">
          <div class="section-head">
            <p class="overline">How it works</p>
            <h2>From question to evidence in four steps</h2>
          </div>
          <ol class="flow">
            <li v-for="s in steps" :key="s.title">
              <span class="flow__icon"><SiteStepIcon :step="s.icon" /></span>
              <div>
                <h4>{{ s.title }}</h4>
                <p class="small muted">{{ s.text }}</p>
              </div>
            </li>
          </ol>
        </div>
      </section>
    </div>

    <!-- 6 · About Nick -->
    <section class="section">
      <div class="container about">
        <img class="portrait" src="/images/nick-portrait.jpg" alt="Nikolay Tonev" width="1024" height="1536" loading="lazy" style="object-fit: cover; width: 100%; height: auto" />
        <div>
          <p class="overline">Who you work with</p>
          <h2>One accountable expert, not an account team</h2>
          <p class="lead">I’m Nikolay Tonev. You work with me directly, from the first call to the handover.</p>
          <div class="credentials">
            <SiteCredentialBadge />
            <a class="gcp-badge" href="https://www.credly.com/users/nikolay-tonev.a60be769/badges" target="_blank" rel="noopener noreferrer">
              <img src="/certs/cert_ca.webp" alt="" width="64" height="64" loading="lazy" />
              <span><strong>Google Cloud Professional Cloud Architect</strong><small>Certified 2019, 2022 and 2024</small></span>
            </a>
          </div>
          <p style="margin-top: var(--space-6)"><NuxtLink to="/about" class="btn btn--ghost">More about me →</NuxtLink></p>
        </div>
      </div>
    </section>

    <!-- 7 · Insights -->
    <section v-if="insights.length" class="section section--sunken">
      <div class="container">
        <div class="section-head">
          <p class="overline">Insights</p>
          <h2>Notes on AI governance</h2>
        </div>
        <div class="grid grid-3">
          <article v-for="a in insights" :key="a.slug" class="card card--link post">
            <p class="meta">{{ formatDate(a.date) }}</p>
            <h3><NuxtLink :to="`/blog/${a.slug}`" class="card__link" style="color: inherit; text-decoration: none">{{ a.title }}</NuxtLink></h3>
          </article>
        </div>
        <p style="margin-top: var(--space-8)"><NuxtLink to="/blog" class="btn btn--ghost">All insights →</NuxtLink></p>
      </div>
    </section>

    <!-- 8 · FAQ -->
    <section class="section">
      <div class="container" style="max-width: 820px">
        <div class="section-head">
          <p class="overline">Questions</p>
          <h2>Before you ask</h2>
        </div>
        <div class="faq">
          <details v-for="f in faqs" :key="f.q">
            <summary>{{ f.q }}</summary>
            <p>{{ f.a }}</p>
          </details>
        </div>
      </div>
    </section>

    <!-- 9 · Closing CTA + contact -->
    <section id="contact" class="section section--tint hexgrid closing">
      <div class="container">
        <h2>Not sure where to start?</h2>
        <p class="lead">A 30-minute call finds the first step.</p>
        <div class="row" style="justify-content: center">
          <NuxtLink to="/booking" class="btn btn--primary">Book a call</NuxtLink>
          <a href="#main" class="btn btn--secondary">Ask the advisor</a>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.credentials { display: flex; flex-wrap: wrap; gap: var(--space-4); align-items: center; margin-top: var(--space-6); }
.gcp-badge { display: inline-flex; gap: var(--space-3); align-items: center; color: inherit; text-decoration: none; }
.gcp-badge img { width: 64px; height: 64px; object-fit: contain; }
.gcp-badge span { display: grid; font-size: var(--text-sm); }
.gcp-badge small { color: var(--text-muted); }
</style>
