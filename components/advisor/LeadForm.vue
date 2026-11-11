<script setup lang="ts">
// Lead capture for the advisor ("Leave contact") and the order request.
// POST /api/leads (backend #124). Consent is required and unticked by default.
import { trackEvent } from '~/services/analytics'

const props = withDefaults(
  defineProps<{ source: 'advisor' | 'order' | 'contact' | 'ivgil'; prefill?: string; packageId?: string; submitLabel?: string }>(),
  { prefill: '', packageId: undefined, submitLabel: 'Send' },
)

const form = reactive({ name: '', email: '', company: '', message: props.prefill, consent: false, insights: false, role: '', interest: 'early_access' as 'early_access' | 'design_partner' })
const status = ref<'idle' | 'sending' | 'sent' | 'error'>('idle')
const errors = ref<Record<string, string>>({})

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

function validate(): boolean {
  const e: Record<string, string> = {}
  if (!form.name.trim()) e.name = 'Please enter your name.'
  if (!EMAIL_RE.test(form.email.trim())) e.email = 'Please enter a valid work email.'
  if (!form.company.trim()) e.company = 'Please enter your company.'
  if (!form.consent) e.consent = 'Please confirm so we can answer your request.'
  errors.value = e
  return Object.keys(e).length === 0
}

async function submit() {
  if (!validate()) return
  status.value = 'sending'
  try {
    await $fetch('/api/leads', {
      method: 'POST',
      body: {
        source: props.source,
        package_id: props.packageId,
        name: form.name.trim(),
        email: form.email.trim(),
        company: form.company.trim(),
        message: form.message.trim(),
        consent: form.consent,
        insights_optin: form.insights,
        ...(props.source === 'ivgil' ? { role: form.role.trim(), interest: form.interest } : {}),
      },
      timeout: 15000,
    })
    status.value = 'sent'
    trackEvent('lead_submitted', { source: props.source, package_id: props.packageId })
  } catch {
    status.value = 'error'
  }
}
</script>

<template>
  <div v-if="status === 'sent'" class="notice" role="status">
    <p><strong>Thanks.</strong> Nick will reply within one working day.</p>
  </div>
  <form v-else novalidate @submit.prevent="submit">
    <div class="two">
      <div class="field" :class="{ 'field--error': errors.name }">
        <label :for="`${source}-name`">Name</label>
        <input :id="`${source}-name`" v-model="form.name" class="input" autocomplete="name" required />
        <span v-if="errors.name" class="hint">{{ errors.name }}</span>
      </div>
      <div class="field" :class="{ 'field--error': errors.email }">
        <label :for="`${source}-email`">Work email</label>
        <input :id="`${source}-email`" v-model="form.email" class="input" type="email" autocomplete="email" required />
        <span v-if="errors.email" class="hint">{{ errors.email }}</span>
      </div>
    </div>
    <div class="field" :class="{ 'field--error': errors.company }">
      <label :for="`${source}-company`">Company</label>
      <input :id="`${source}-company`" v-model="form.company" class="input" autocomplete="organization" required />
      <span v-if="errors.company" class="hint">{{ errors.company }}</span>
    </div>
    <template v-if="source === 'ivgil'">
      <div class="field">
        <label :for="`${source}-role`">Role (optional)</label>
        <input :id="`${source}-role`" v-model="form.role" class="input" autocomplete="organization-title" />
      </div>
      <fieldset class="field">
        <legend>I'm interested in</legend>
        <label class="consent"><input v-model="form.interest" type="radio" value="early_access" /> <span>Early access</span></label>
        <label class="consent"><input v-model="form.interest" type="radio" value="design_partner" /> <span>Becoming a design partner (free implementation, in return for testing and feedback)</span></label>
      </fieldset>
    </template>
    <div class="field">
      <label :for="`${source}-message`">{{ source === 'ivgil' ? 'Your question (optional)' : 'Message' }}</label>
      <textarea :id="`${source}-message`" v-model="form.message" rows="4" maxlength="2000" />
    </div>
    <label class="consent">
      <input v-model="form.consent" type="checkbox" required />
      <span>
        I agree that ivmanto may use my details<template v-if="source === 'advisor'"> and my advisor conversation</template>
        to answer this request. Withdraw any time. <NuxtLink to="/privacy-policy">Privacy notice</NuxtLink>
      </span>
    </label>
    <p v-if="errors.consent" class="small" style="color: var(--danger)">{{ errors.consent }}</p>
    <label class="consent">
      <input v-model="form.insights" type="checkbox" />
      <span>Also send me occasional insights by email (optional; you confirm by email first).</span>
    </label>
    <div v-if="status === 'error'" class="notice notice--error" role="alert">
      <p>That didn't go through. Please try again, or <NuxtLink to="/booking">book a call</NuxtLink>.</p>
    </div>
    <button type="submit" class="btn btn--primary" :disabled="status === 'sending'">
      {{ status === 'sending' ? 'Sending…' : submitLabel }}
    </button>
  </form>
</template>
