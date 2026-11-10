<script setup lang="ts">
// DRAFT: Brio drafts the final text; Ivgil and Nick review (EU AI Act Art. 50, GDPR).
useSeoMeta({ title: 'AI use on this site', description: 'How the AI service advisor on ivmanto.com works and what happens to your message.' })
</script>

<template>
  <section class="section">
    <div class="container" style="max-width: 760px">
      <p class="overline">Draft, pending review</p>
      <h1>AI use on this site</h1>
      <p class="lead">The service advisor on the home page is an AI assistant. It suggests which of our services fits your situation. It does not give legal advice.</p>
      <h2>What happens to your message</h2>
      <ul class="ticks">
        <li>Processor: [Mistral AI, EU]. Purpose: suggesting a service.</li>
        <li>Retention: [N days], then deleted. [No training on your data.]</li>
        <li>Please don’t include personal or confidential data.</li>
        <li>If you leave your contact, the conversation is shared with Nick only with your consent.</li>
      </ul>
      <h2>Why the AI label</h2>
      <p>EU AI Act Art. 50: you should always know when you are talking to an AI system. [Text checked by Ivgil.]</p>
      <p><NuxtLink to="/privacy-policy">Privacy notice</NuxtLink></p>
    </div>
  </section>
</template>
