<script setup lang="ts">
// About (DS v1.5, pages/about.md; approved by Nick, 2026-11-13).
// The two "Why governance" paragraphs are Nick's own wording: any change needs his yes first.
import { CREDLY_URL, MIT_CREDENTIAL_URL } from '~/utils/credentials'
import { jsonLd } from '~/utils/insights'

const DESCRIPTION =
  'Nikolay Tonev helps small and mid-sized companies govern data and AI: clear owners, known risks, and proof for the EU AI Act and ISO/IEC 42001.'

useSeoMeta({
  title: 'About Nikolay Tonev | ivmanto',
  description: DESCRIPTION,
  ogTitle: 'About Nikolay Tonev | ivmanto',
  ogDescription: DESCRIPTION,
  twitterTitle: 'About Nikolay Tonev | ivmanto',
  twitterDescription: DESCRIPTION,
})

useHead({
  script: [
    {
      id: 'person-schema',
      type: 'application/ld+json',
      innerHTML: jsonLd({
        '@context': 'https://schema.org',
        '@type': 'Person',
        '@id': 'https://ivmanto.com/about#person',
        name: 'Nikolay Tonev',
        url: 'https://ivmanto.com/about',
        image: 'https://ivmanto.com/images/nick-portrait.jpg',
        jobTitle: 'AI & Data Governance Consultant',
        worksFor: { '@id': 'https://ivmanto.com/#organization' },
        sameAs: ['https://linkedin.com/in/nikolaytonev', 'https://github.com/dasiyes', CREDLY_URL],
        knowsAbout: [
          'AI governance',
          'Data governance',
          'Data management',
          'DAMA-DMBOK',
          'EU AI Act',
          'ISO/IEC 42001',
          'AI risk management',
          'Google Cloud',
          'Data architecture',
        ],
      }),
    },
  ],
})

const principles = [
  { title: 'One contact', text: 'You work with me directly, from the first call to the handover.' },
  { title: 'Fixed scope, fixed price', text: 'Both are agreed before the work starts.' },
  { title: 'Your tools', text: 'Registers, policies and reports stay in the tools you already use.' },
  { title: 'Not legal advice', text: 'Where the law needs interpreting, I work alongside your lawyer.' },
]

// A Google Cloud certificate always names its years (Nick, 2026-11-11): a bare badge reads as current.
const googleCertificates = [
  { name: 'Google Cloud Professional Cloud Architect', meta: 'Google Cloud · certified 2019, 2022 and 2024' },
  { name: 'Google Cloud Professional Data Engineer', meta: 'Google Cloud · certified 2023' },
  { name: 'Google Cloud Professional Cloud Developer', meta: 'Google Cloud · certified 2023' },
]
</script>

<template>
  <div>
    <!-- 1 · Who -->
    <section class="page-head hexgrid" style="padding-bottom: var(--space-16)">
      <div class="container about-hero">
        <img class="portrait" src="/images/nick-portrait.jpg" alt="Nikolay Tonev" width="1024" height="1536" />
        <div>
          <p class="overline">About</p>
          <h1>AI and data governance from someone who builds the systems</h1>
          <p class="lead">I’m Nikolay Tonev. I help small and mid-sized companies govern their data and their AI: who decides, what can go wrong, and what the law asks for.</p>
          <div class="row" style="margin-top: var(--space-6)">
            <NuxtLink to="/booking" class="btn btn--primary">Book a call</NuxtLink>
            <NuxtLink to="/packages" class="btn btn--ghost">See the packages <SiteChevIcon /></NuxtLink>
          </div>
          <ul class="profile-links">
            <li><a href="https://linkedin.com/in/nikolaytonev" target="_blank" rel="noopener">LinkedIn</a></li>
            <li><a href="https://github.com/dasiyes" target="_blank" rel="noopener">GitHub</a></li>
            <li><a :href="CREDLY_URL" target="_blank" rel="noopener">Credly</a></li>
          </ul>
        </div>
      </div>
    </section>

    <!-- 2 · Why governance (Nick's wording, verbatim) -->
    <section class="section section--sunken">
      <div class="container two-col">
        <div>
          <p class="overline">Why governance</p>
          <h2>Rules only work if they fit how the systems are built</h2>
        </div>
        <div>
          <p>I’ve spent more than 20 years in enterprise IT, starting as a business analyst at Wrigley (Mars), and the last ten years designing data platforms and data products in the cloud. At MediaMarktSaturn I built the customer data platform that 100+ marketers relied on, and at Eviden I architected governed data products and generative AI agents on Google Cloud for enterprise clients. Today I bring that experience to SMEs as an independent AI and data governance consultant.</p>
          <p>At Eviden (Atos Group), I ran a Google Agentspace proof of concept and watched AI agents spread through the enterprise much faster than anyone had planned for. Getting an agent working was the easy part. The harder questions were who controls it, what data it touches, and whether it complies with the rules. That’s why I moved into AI governance: to help companies scale AI quickly without losing control of it.</p>
        </div>
      </div>
    </section>

    <!-- 3 · How I work -->
    <section class="section">
      <div class="container">
        <div class="section-head">
          <p class="overline">How I work</p>
          <h2>Four things you can count on</h2>
        </div>
        <ul class="principles">
          <li v-for="p in principles" :key="p.title">
            <h3>{{ p.title }}</h3>
            <p>{{ p.text }}</p>
          </li>
        </ul>
      </div>
    </section>

    <!-- 4 · Hands-on background -->
    <section class="section section--sunken">
      <div class="container two-col">
        <div>
          <p class="overline">Hands-on background</p>
          <h2>I still build data platforms and AI agents</h2>
        </div>
        <div>
          <p>Data platforms and AI agents, built in your own Google Cloud project. That work keeps the governance advice practical.</p>
          <p><NuxtLink to="/data-ai" class="btn btn--ghost">Data &amp; AI engineering <SiteChevIcon /></NuxtLink></p>
        </div>
      </div>
    </section>

    <!-- 5 · Credentials -->
    <section class="section">
      <div class="container">
        <div class="section-head">
          <p class="overline">Credentials</p>
          <h2>Training and certificates</h2>
        </div>
        <!-- The certificate image is served from this site: a hotlink to Accredible would send
             the visitor's IP address to a third party before any consent. -->
        <div class="credential credential--certificate" style="margin-bottom: var(--space-8)">
          <a class="credential__image" :href="MIT_CREDENTIAL_URL" target="_blank" rel="noopener">
            <img src="/images/mit-pe-certificate.png" width="600" height="424" loading="lazy" alt="Certificate from MIT Professional Education: Nikolay Tonev has successfully completed the online program Applied Agentic AI for Organizational Transformation, September 1 to October 27, 2026." />
          </a>
          <p>
            <strong style="color: var(--text)">MIT Professional Education certificate</strong><br />
            Applied Agentic AI for Organizational Transformation · 2026<br />
            <a :href="MIT_CREDENTIAL_URL" target="_blank" rel="noopener">Verified on Accredible · View credential</a>
          </p>
        </div>
        <ul class="cred-list">
          <li v-for="c in googleCertificates" :key="c.name">
            <div><strong>{{ c.name }}</strong><small>{{ c.meta }}</small></div>
            <a :href="CREDLY_URL" target="_blank" rel="noopener">View on Credly</a>
          </li>
          <li>
            <div><strong>Data Architecture Workshop: Practical Exercises in Architecture Design</strong><small>DATAVERSITY · 2021 · 24 CEUs</small></div>
            <a href="/pdfs/dataversity-data-architecture-workshop-2021.pdf" target="_blank" rel="noopener">Certificate (PDF)</a>
          </li>
        </ul>
        <p class="small muted" style="margin-top: var(--space-4)">
          Google Cloud certificates are listed with the years they were earned.
          <a :href="CREDLY_URL" target="_blank" rel="noopener">All badges on Credly</a>
        </p>
      </div>
    </section>

    <!-- 6 · Closing (same as home) -->
    <section class="section section--tint hexgrid closing">
      <div class="container">
        <h2>Not sure where to start?</h2>
        <p class="lead">A 30-minute call finds the first step.</p>
        <div class="row" style="margin-top: var(--space-6)">
          <NuxtLink to="/booking" class="btn btn--primary">Book a call</NuxtLink>
          <NuxtLink to="/#main" class="btn btn--secondary">Ask the advisor</NuxtLink>
        </div>
      </div>
    </section>
  </div>
</template>
