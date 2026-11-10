<script setup lang="ts">
// Consent gate for analytics (GTM loads only after "Accept"). Accept and Reject have
// equal weight (TTDSG/GDPR). The footer's "Cookie settings" reopens this banner via
// the 'ivm:cookie-settings' window event.
import { initGtm } from '~/services/analytics'

const COOKIE_CONSENT_KEY = 'cookie_consent'

const showBanner = ref(false)

function readConsent(): string | null {
  try {
    return localStorage.getItem(COOKIE_CONSENT_KEY)
  } catch {
    return null
  }
}

function writeConsent(value: 'accepted' | 'declined') {
  try {
    localStorage.setItem(COOKIE_CONSENT_KEY, value)
  } catch {
    // Storage blocked: the choice applies to this page view only.
  }
}

const reopen = () => {
  showBanner.value = true
}

onMounted(() => {
  if (!readConsent()) showBanner.value = true
  window.addEventListener('ivm:cookie-settings', reopen)
})

onBeforeUnmount(() => window.removeEventListener('ivm:cookie-settings', reopen))

const acceptCookies = () => {
  writeConsent('accepted')
  showBanner.value = false
  initGtm()
}

const declineCookies = () => {
  writeConsent('declined')
  showBanner.value = false
}
</script>

<template>
  <div v-if="showBanner" class="cookie" role="dialog" aria-labelledby="cookie-title">
    <h4 id="cookie-title">Cookies</h4>
    <p>
      We use analytics cookies only if you agree. You can change this any time under
      “Cookie settings” in the footer.
      <NuxtLink to="/privacy-policy">Details</NuxtLink>
    </p>
    <div class="row">
      <button type="button" class="btn btn--secondary btn--sm" @click="declineCookies">Reject</button>
      <button type="button" class="btn btn--primary btn--sm" @click="acceptCookies">Accept</button>
    </div>
  </div>
</template>
