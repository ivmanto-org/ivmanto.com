<script setup lang="ts">
// Order flow, B2B order request (docs/design-system/pages/order.html). Whether this
// becomes Stripe checkout or stays "order request + invoice" is open with Nick (#112).
const route = useRoute()
const { pkg, pillar } = useCatalogue()

const current = pkg(String(route.params.id))
if (!current) throw createError({ statusCode: 404, statusMessage: 'Package not found', fatal: true })
const plName = pillar(current.pillar)?.name ?? LINE_NAMES[current.line]

useSeoMeta({ title: `Order: ${current.name}`, robots: 'noindex' })
</script>

<template>
  <div class="container">
    <div class="order">
      <div>
        <ol class="stepper" aria-label="Order steps">
          <li class="done"><span>Package</span></li>
          <li aria-current="step"><span>Your details</span></li>
          <li><span>Confirmation</span></li>
        </ol>
        <h1 class="h2">Order request</h1>
        <p class="muted">
          No payment now. Nick confirms scope and start date by email, then sends the invoice. [Final order terms
          pending Nick.]
        </p>
        <AdvisorLeadForm source="order" :package-id="current.id" :prefill="`Order request: ${current.name}`" submit-label="Send order request" />
      </div>
      <aside class="card summary" aria-label="Summary">
        <h3>{{ current.name }}</h3>
        <dl>
          <dt>Service</dt><dd>{{ plName }}</dd>
          <dt>Duration</dt><dd>{{ formatDuration(current) }}</dd>
          <dt>Deliverable</dt><dd>{{ current.deliverable }}</dd>
          <dt class="total">Price ({{ priceTerms(current) }})</dt><dd class="total">{{ formatPrice(current) }}</dd>
        </dl>
        <NuxtLink :to="`/packages/${current.id}`" class="btn btn--ghost" style="margin-top: var(--space-4)">← Back to package</NuxtLink>
      </aside>
    </div>
  </div>
</template>
