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
        <li>Your message is sent to our AI provider, Mistral AI (EU), only to generate the answer.</li>
        <li>We don’t store the conversation. How long Mistral AI itself keeps data is decided by Mistral AI under its own terms; see our <NuxtLink to="/privacy-policy">privacy policy</NuxtLink>.</li>
        <li>Please don’t include personal or confidential information.</li>
        <li>If you send the contact form, Nick receives your details and a short summary of the chat to answer your request.</li>
      </ul>
      <h2>Why the AI label</h2>
      <p>EU AI Act Art. 50: you should always know when you are talking to an AI system. [Text checked by Ivgil.]</p>
      <h2>Texts on this site</h2>
      <p>The texts on this site are drafted with the help of AI tools. Nick reviews and approves every text before it is published.</p>
      <p><NuxtLink to="/privacy-policy">Privacy notice</NuxtLink></p>
    </div>
  </section>
</template>

<style scoped>
h2 { margin-top: var(--space-12); }
</style>
