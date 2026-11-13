<script setup lang="ts">
useSeoMeta({ title: 'Contact', description: 'Send Nick a message. You get a reply within one working day.' })
</script>

<template>
  <section class="section">
    <div class="container" style="max-width: 760px">
      <p class="overline">Contact</p>
      <h1>Send a message</h1>
      <p class="lead">Tell me briefly what you need. I reply within one working day.</p>
      <AdvisorLeadForm source="contact" submit-label="Send message" />
    </div>
  </section>
</template>
