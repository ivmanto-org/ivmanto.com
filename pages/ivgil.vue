<script setup lang="ts">
// ivGIL landing page (Nick, 2026-11-11: named at launch; early access by default, design
// partner as the second option). Product rules: no "compliant", "certified", "blockchain",
// no absolute claims. Regulatory wording is DRAFT until Ivgil's content gate passes.
useSeoMeta({
  title: 'ivGIL: the AI governance platform we are building',
  description: 'ivGIL is an AI governance platform for EU companies, in development: AI inventory, EU AI Act classification, runtime control of AI agents and verifiable evidence. Request early access.',
})

const features = [
  { title: 'AI inventory', text: 'One register of the AI systems and agents your company uses, each with an owner.' },
  { title: 'EU AI Act classification', text: 'A structured way to work out which risk class each system falls into and which obligations follow.' },
  { title: 'Runtime control of AI agents', text: 'Rules for what your AI agents may do, checked while they run, not only on paper.' },
  { title: 'Verifiable evidence', text: 'Records of decisions and controls that an auditor can verify without having to trust ivGIL.' },
]
</script>

<template>
  <div>
    <section class="pillar-hero hexgrid">
      <div class="container">
        <p class="overline">ivGIL · in development</p>
        <h1>The AI governance platform we are building</h1>
        <p class="lead" style="max-width: 720px">
          ivGIL turns the governance work we do with clients into software: know your AI systems, classify them,
          control what your AI agents do, and keep evidence that holds up. It is in development now, and we are
          looking for early users and design partners.
        </p>
        <div class="row" style="margin-top: var(--space-6)">
          <a href="#early-access" class="btn btn--primary">Request early access</a>
          <NuxtLink to="/packages" class="btn btn--secondary">Need results now? Start with a package</NuxtLink>
        </div>
      </div>
    </section>

    <section class="section section--sunken">
      <div class="container">
        <div class="section-head"><p class="overline">What it does</p><h2>Four things, in one place</h2></div>
        <div class="grid grid-2">
          <article v-for="f in features" :key="f.title" class="card">
            <h3>{{ f.title }}</h3>
            <p class="muted">{{ f.text }}</p>
          </article>
        </div>
      </div>
    </section>

    <section class="section">
      <div class="container">
        <div class="section-head"><p class="overline">Two ways in</p><h2>Early access or design partner</h2></div>
        <div class="grid grid-2">
          <article class="card">
            <h3>Early access</h3>
            <p class="muted">Be among the first to use ivGIL when the first version is ready. We tell you when it opens and what it can do.</p>
          </article>
          <article class="card">
            <h3>Design partner</h3>
            <p class="muted">We implement ivGIL in your company free of charge. In return, you test it in daily work and give us regular feedback. Places are limited.</p>
          </article>
        </div>
      </div>
    </section>

    <section id="early-access" class="section section--sunken">
      <div class="container" style="max-width: 760px">
        <h2>Request early access</h2>
        <p class="lead">Tell us who you are and which option interests you. Nick replies within one working day.</p>
        <AdvisorLeadForm source="ivgil" submit-label="Send request" />
      </div>
    </section>
  </div>
</template>
