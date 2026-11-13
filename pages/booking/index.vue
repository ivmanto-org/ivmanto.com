<script setup lang="ts">
// Book a call (DS v1.5, pages/booking.md). Layout, copy and formats only:
// the API calls and what they send are the ones the page had before.
import { getGaSessionInfo } from '~/services/analytics'
import { toYYYYMMDD } from '~/composables/useDate'
import { CALL_MINUTES, formatLongDate, formatTime24 } from '~/utils/bookingFormat'

useSeoMeta({
  title: 'Book a call | ivmanto.com',
  description: 'Book a 30-minute call with Nikolay Tonev to find the first step for AI governance in your company.',
  robots: 'noindex, nofollow',
})

type TimeSlot = { id: string; start: string; end: string }
type Step = 1 | 2 | 3

const GA_MEASUREMENT_ID = 'G-W1TJ3KMZ6V'
const COOKIE_CONSENT_KEY = 'cookie_consent'
const HTTP_CONFLICT = 409
const SKELETON_SLOTS = 3
// Visitors can book up to ~3 months ahead. The AfB scheduler fills the calendar
// 12 weeks ahead, so the 90-day window always has bookable slots.
const BOOKING_WINDOW_DAYS = 90
const STEPS = ['Pick a time', 'Your details', 'Confirmed'] as const

const route = useRoute()

// `today` is captured at setup so the bounds do not drift while the page is open.
const today = new Date()
const maxBookingDate = new Date(today.getFullYear(), today.getMonth(), today.getDate() + BOOKING_WINDOW_DAYS)

const step = ref<Step>(1)
const selectedDate = ref(new Date())
const availableSlots = ref<TimeSlot[]>([])
const selectedSlot = ref<TimeSlot | null>(null)
const isLoadingSlots = ref(false)
const isSubmitting = ref(false)
const slotsError = ref(false)
const slotTaken = ref(false)
const submitError = ref(false)
const fieldErrors = ref<{ name?: string; email?: string }>({})

const bookingDetails = ref({ name: '', email: '', notes: '' })

const nameInput = ref<HTMLInputElement | null>(null)
const doneHeading = ref<HTMLElement | null>(null)

// The visitor's IANA time zone (e.g. "Europe/Berlin"). Sent with the booking so the
// backend can localise the confirmation email and the calendar file.
const timezoneIana = computed(() => {
  if (typeof Intl === 'undefined' || typeof Intl.DateTimeFormat === 'undefined') return ''
  return Intl.DateTimeFormat().resolvedOptions().timeZone
})

const pickedLine = computed(() =>
  selectedSlot.value ? `${formatLongDate(selectedDate.value, true)}, ${formatTime24(selectedSlot.value.start)}` : '',
)

async function fetchAvailability(date: Date) {
  isLoadingSlots.value = true
  slotsError.value = false
  selectedSlot.value = null
  availableSlots.value = []
  try {
    const response = await fetch(`/api/booking/availability?date=${toYYYYMMDD(date)}`)
    if (!response.ok) throw new Error(`Server responded with status ${response.status}`)
    const data = await response.json().catch(() => null)
    const slots: TimeSlot[] = Array.isArray(data) ? data : []
    const now = Date.now()
    // Only today can hold times that have already passed.
    availableSlots.value = slots.filter((slot) => new Date(slot.start).getTime() > now)
  } catch (e) {
    console.error('Failed to fetch availability:', e)
    slotsError.value = true
  } finally {
    isLoadingSlots.value = false
  }
}

function selectSlot(slot: TimeSlot) {
  selectedSlot.value = slot
  slotTaken.value = false
  submitError.value = false
  step.value = 2
  nextTick(() => nameInput.value?.focus())
}

function backToTimes() {
  step.value = 1
  selectedSlot.value = null
}

function validate(): boolean {
  const errors: { name?: string; email?: string } = {}
  if (!bookingDetails.value.name.trim()) errors.name = 'Please enter your name.'
  if (!/^\S+@\S+\.\S+$/.test(bookingDetails.value.email.trim())) errors.email = 'Please enter a valid email address.'
  fieldErrors.value = errors
  return Object.keys(errors).length === 0
}

// Analytics ids go along only while the visitor's consent stands. Cookies from an
// earlier "Accept" can outlive a later "Reject", so the cookies alone are not enough.
function analyticsIds() {
  let hasConsent = false
  try {
    hasConsent = localStorage.getItem(COOKIE_CONSENT_KEY) === 'accepted'
  } catch {
    hasConsent = false
  }
  return hasConsent ? getGaSessionInfo(GA_MEASUREMENT_ID) : { clientId: null, sessionId: null }
}

async function handleBookingSubmit() {
  if (!selectedSlot.value || isSubmitting.value || !validate()) return
  isSubmitting.value = true
  submitError.value = false
  try {
    const { clientId, sessionId } = analyticsIds()
    const response = await fetch('/api/booking/book', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        ga_client_id: clientId,
        ga_session_id: sessionId,
        eventId: selectedSlot.value.id,
        visitorTimezone: timezoneIana.value,
        ...bookingDetails.value,
      }),
    })
    if (response.status === HTTP_CONFLICT) {
      // Someone else took the time: back to step 1 with a fresh list for the same day.
      step.value = 1
      await fetchAvailability(selectedDate.value)
      slotTaken.value = true
      return
    }
    if (!response.ok) throw new Error(`Server responded with status ${response.status}`)
    step.value = 3
    nextTick(() => doneHeading.value?.focus())
  } catch (e) {
    console.error('Booking failed:', e)
    submitError.value = true
  } finally {
    isSubmitting.value = false
  }
}

watch(selectedDate, (newDate) => {
  slotTaken.value = false
  fetchAvailability(newDate)
})

onMounted(() => {
  const { topic, summary } = route.query
  if (topic && summary && typeof topic === 'string' && typeof summary === 'string') {
    bookingDetails.value.notes = `I'm interested in exploring the topic: "${topic}".\n\nSummary: ${summary}`
    window.history.replaceState({}, '', route.path)
  }
  fetchAvailability(selectedDate.value)
})
</script>

<template>
  <section class="section" style="padding-top: var(--space-12)">
    <div class="container booking">
      <div>
        <p class="overline">Book a call</p>
        <h1>30 minutes to find the first step</h1>
        <p class="lead">Tell me where you stand with AI. You leave the call knowing which service fits, or that you need none.</p>
        <ul class="facts">
          <li>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M12 7v5l3 2" /></svg>
            <div><strong>30 minutes</strong><small>[Free and without obligation: Nick to confirm]</small></div>
          </li>
          <li>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="3" y="6" width="12" height="12" rx="2" /><path d="M15 10.5l6-3.5v10l-6-3.5z" /></svg>
            <div><strong>Video call</strong><small>The link comes with the confirmation email.</small></div>
          </li>
          <li>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="8" r="4" /><path d="M4 20.5a8 8 0 0 1 16 0" /></svg>
            <div><strong>With Nikolay Tonev</strong><small>[In English or German: Nick to confirm]</small></div>
          </li>
        </ul>
        <p class="small muted" style="margin-top: var(--space-6)">Prefer to write? <NuxtLink to="/contact">Send a message</NuxtLink></p>
      </div>

      <div class="booking__panel">
        <ol class="stepper" aria-label="Booking steps">
          <li v-for="(label, i) in STEPS" :key="label" :class="{ done: i + 1 < step }" :aria-current="i + 1 === step ? 'step' : undefined">
            <span>{{ label }}</span>
          </li>
        </ol>

        <!-- The date, time zone and times depend on the visitor's clock, so the panel renders
             in the browser only (the pre-rendered page would carry the build date). -->
        <ClientOnly>
          <!-- Step 1: pick a time -->
          <div v-if="step === 1">
            <div v-if="slotTaken" class="notice notice--error" role="alert" style="margin-bottom: var(--space-4)">
              <p><strong>This time has just been taken.</strong> Please pick another one.</p>
            </div>
            <BookingCalendarMonthGrid v-model="selectedDate" :min-date="today" :max-date="maxBookingDate" />
            <p class="cal__tz">Times are shown in your time zone: {{ timezoneIana || 'unknown' }}</p>

            <h3 class="slots__title">Free times on {{ formatLongDate(selectedDate) }}</h3>
            <div v-if="isLoadingSlots" class="slots">
              <span v-for="n in SKELETON_SLOTS" :key="n" class="slot slot--skeleton" aria-hidden="true"></span>
              <span class="sr-only" role="status">Loading the times…</span>
            </div>
            <div v-else-if="slotsError" class="notice notice--error" role="alert">
              <p><strong>That didn’t work.</strong> Please try again, or <NuxtLink to="/contact">send a message</NuxtLink> instead.</p>
            </div>
            <div v-else-if="availableSlots.length > 0" class="slots">
              <button v-for="slot in availableSlots" :key="slot.id" type="button" class="slot" @click="selectSlot(slot)">
                {{ formatTime24(slot.start) }}
              </button>
            </div>
            <div v-else class="notice" role="status">
              <p style="margin: 0">No free times on this day. Please pick another date.</p>
            </div>
          </div>

          <!-- Step 2: your details -->
          <form v-else-if="step === 2" novalidate @submit.prevent="handleBookingSubmit">
            <div class="picked">
              <span><strong>{{ pickedLine }}</strong> · {{ CALL_MINUTES }} minutes</span>
              <button type="button" @click="backToTimes">Change</button>
            </div>
            <div class="field" :class="{ 'field--error': fieldErrors.name }">
              <label for="bk-name">Name</label>
              <input id="bk-name" ref="nameInput" v-model="bookingDetails.name" class="input" autocomplete="name" required :aria-invalid="fieldErrors.name ? 'true' : undefined" :aria-describedby="fieldErrors.name ? 'bk-name-err' : undefined" />
              <span v-if="fieldErrors.name" id="bk-name-err" class="hint">{{ fieldErrors.name }}</span>
            </div>
            <div class="field" :class="{ 'field--error': fieldErrors.email }">
              <label for="bk-mail">Work email</label>
              <input id="bk-mail" v-model="bookingDetails.email" class="input" type="email" autocomplete="email" required :aria-invalid="fieldErrors.email ? 'true' : undefined" aria-describedby="bk-mail-hint" />
              <span id="bk-mail-hint" class="hint">{{ fieldErrors.email || 'The confirmation and the calendar file go here.' }}</span>
            </div>
            <div class="field">
              <label for="bk-notes">What would you like to talk about? <span class="muted" style="font-weight: 400">(optional)</span></label>
              <textarea id="bk-notes" v-model="bookingDetails.notes" rows="3"></textarea>
            </div>
            <p class="form-note">
              I use these details only to hold and prepare this call (Art. 6(1)(b) GDPR).
              <NuxtLink to="/privacy-policy">Privacy notice</NuxtLink>
            </p>
            <div v-if="submitError" class="notice notice--error" role="alert" style="margin-bottom: var(--space-4)">
              <p><strong>That didn’t work.</strong> Please try again, or <NuxtLink to="/contact">send a message</NuxtLink> instead.</p>
            </div>
            <button class="btn btn--primary btn--block" type="submit" :disabled="isSubmitting">
              {{ isSubmitting ? 'Booking…' : 'Book the call' }}
            </button>
          </form>

          <!-- Step 3: confirmed -->
          <div v-else>
            <span class="done-mark">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M5 12.5l4.5 4.5L19 7.5" /></svg>
            </span>
            <h2 ref="doneHeading" class="h3" tabindex="-1">Your call is booked</h2>
            <p><strong>{{ pickedLine }}</strong> · {{ CALL_MINUTES }} minutes · video call</p>
            <p class="muted">A confirmation with a calendar file is on its way to your inbox. It also has a link to cancel.</p>
            <div class="row" style="margin-top: var(--space-6)">
              <NuxtLink to="/" class="btn btn--secondary">Back to home</NuxtLink>
              <NuxtLink to="/blog" class="btn btn--ghost">Read the insights <SiteChevIcon /></NuxtLink>
            </div>
          </div>

          <template #fallback>
            <p class="muted small" role="status">Loading the calendar…</p>
          </template>
        </ClientOnly>
      </div>
    </div>
  </section>
</template>

<style scoped>
h2[tabindex='-1']:focus { outline: none; }
</style>
