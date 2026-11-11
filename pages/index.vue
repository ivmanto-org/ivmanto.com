<script setup lang="ts">
// Home: the funnel (docs/design-system/pages/README.md, section order and why).
// Copy in [brackets] and '€ [price]' are placeholders for Nick.
const { pillars, packagesFor, featured } = useCatalogue()
const { fetchArticles, sortedArticles } = useArticles()

await useAsyncData('home-articles', () => fetchArticles())
const insights = computed(() => sortedArticles.value.slice(0, 3))

const frameworks = ['EU AI Act', 'ISO/IEC 42001', 'ISO/IEC 27001', 'NIST AI RMF', 'GDPR']

const steps = [
  { title: 'Describe', text: 'Ask the advisor above or book a 30-minute call.' },
  { title: 'Choose', text: 'Pick a package. Scope and price are fixed before we start.' },
  { title: 'Deliver', text: 'Workshops with your team, documents built with you, a clear owner for every item.' },
  { title: 'Hand over', text: 'You keep every artefact in your own tools. No lock-in.' },
]

const faqs = [
  {
    q: 'Is this legal advice?',
    a: 'No. We provide governance, risk and compliance consulting. Where legal interpretation is needed, we work alongside your lawyer.',
  },
  { q: 'We’re small. Does the EU AI Act even apply to us?', a: '[Short, accurate answer, checked by Ivgil.]' },
  { q: 'What happens after I order a package?', a: '[Order flow: confirmation, kickoff call scheduled, invoice terms.]' },
  { q: 'What does the AI advisor do with my message?', a: '[Processor, retention, no training on your data; link to privacy notice.]' },
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
    <section class="hero hexgrid">
      <div class="container hero-grid">
        <div class="hero-text">
          <p class="overline">AI governance · risk · compliance</p>
          <h1>Use AI with confidence.<br /><em>Govern it with clarity.</em></h1>
          <p class="lead">
            Fixed-scope services that get your company ready for the EU AI Act and ISO/IEC 42001: clear owners,
            documented risks, audit-ready evidence.
          </p>
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

    <!-- 4 · Packages -->
    <section class="section">
      <div class="container">
        <div class="section-head">
          <p class="overline">Packages</p>
          <h2>Fixed scope. Fixed price. Order directly.</h2>
          <p class="lead">No open-ended retainers. Each package states what you get, how long it takes and what it costs.</p>
        </div>
        <div class="grid grid-3">
          <SitePackageCard v-for="pk in featured()" :key="pk.id" :pkg="pk" />
        </div>
        <p style="margin-top: var(--space-8)"><NuxtLink to="/packages" class="btn btn--ghost">All packages →</NuxtLink></p>
      </div>
    </section>

    <!-- 5 · How it works -->
    <section class="section section--sunken">
      <div class="container">
        <div class="section-head">
          <p class="overline">How it works</p>
          <h2>From question to evidence in four steps</h2>
        </div>
        <ol class="flow">
          <li v-for="s in steps" :key="s.title">
            <h4>{{ s.title }}</h4>
            <p class="muted small">{{ s.text }}</p>
          </li>
        </ol>
      </div>
    </section>

    <!-- 6 · About Nick -->
    <section class="section">
      <div class="container about">
        <img class="portrait" src="/images/nick-portrait.jpg" alt="Nikolay Tonev" width="240" height="307" loading="lazy" style="object-fit: cover; width: 100%; height: auto" />
        <div>
          <p class="overline">Who you work with</p>
          <h2>One accountable expert, not an account team</h2>
          <p class="lead">
            I’m Nikolay Tonev. [Two to three sentences on background: cloud and data architecture, years in enterprise
            IT, why AI governance.] You work with me directly, from the first call to the handover.
          </p>
          <SiteCredentialBadge />
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
            <p class="muted small">{{ a.summary }}</p>
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
    <section id="contact" class="section section--sunken hexgrid closing">
      <div class="container">
        <h2>Not sure where to start?</h2>
        <p class="lead">A 30-minute call is enough to find the first step. No preparation needed.</p>
        <div class="row" style="justify-content: center">
          <NuxtLink to="/booking" class="btn btn--primary">Book a call</NuxtLink>
          <a href="#main" class="btn btn--secondary">Ask the advisor</a>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
/* Blog summaries can be long; keep insight cards even. */
.post .muted { display: -webkit-box; -webkit-line-clamp: 4; -webkit-box-orient: vertical; overflow: hidden; }
</style>
