<script setup lang="ts">
import { getGaSessionInfo } from '~/services/analytics'
import { toYYYYMMDD } from '~/composables/useDate'

useSeoMeta({
  title: 'Contact us | ivmanto.com',
  description:
    'Get in touch with IVMANTO to discuss your data and AI challenges. Book a free consultation or send us a message to start your data transformation journey.',
  robots: 'noindex, nofollow',
})

type TimeSlot = {
  id: string
  start: string
  end: string
}

const route = useRoute()

const selectedDate = ref(new Date())
// Visitors can book up to ~3 months ahead. The AfB scheduler (see
// ~/.hermes/profiles/ivmo/scripts/afb_scheduler.py) populates events with
// WEEKS_AHEAD=12 starting 4 weeks out, so the 90-day window always has
// bookable slots once the Monday cron has run for the week.
const bookingWindowDays = 90
const availableSlots = ref<TimeSlot[]>([])
// Calendar bounds. `today` is captured at mount so the bounds do not drift
// while the page is open; matches the prior day-stepper's behaviour where
// the prev-day button was disabled once the selected day was today.
const today = new Date()
const maxBookingDate = new Date(today.getFullYear(), today.getMonth(), today.getDate() + bookingWindowDays)
const selectedSlot = ref<TimeSlot | null>(null)
const isLoading = ref(false)
const error = ref<string | null>(null)

const bookingDetails = ref({
  name: '',
  email: '',
  notes: '',
})

const isBookingConfirmed = ref(false)

const formattedDate = computed(() => {
  return selectedDate.value.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
})

// timezoneIana is the visitor's IANA timezone name (e.g. "Europe/Berlin").
// Sent in the POST body so the backend can localise the confirmation
// email and .ics attachment to the visitor's zone.
const timezoneIana = computed(() => {
  if (typeof Intl === 'undefined' || typeof Intl.DateTimeFormat === 'undefined') return ''
  return Intl.DateTimeFormat().resolvedOptions().timeZone
})

// timezone is the human-readable form of the visitor's IANA name
// (underscores replaced with spaces) for the on-screen "Timezone:"
// label. Falls back to "Unknown" when Intl is unavailable.
const timezone = computed(() => {
  const iana = timezoneIana.value
  return iana ? iana.replace(/_/g, ' ') : 'Unknown'
})

async function fetchAvailability(date: Date) {
  isLoading.value = true
  error.value = null
  selectedSlot.value = null
  availableSlots.value = []
  try {
    const dateStr = toYYYYMMDD(date)
    const response = await fetch(`/api/booking/availability?date=${dateStr}`)
    if (!response.ok) throw new Error(`Server responded with status ${response.status}`)
    const data = await response.json().catch(() => null)
    let slots = data || []

    const now = new Date()
    const isFetchingForToday =
      date.getFullYear() === now.getFullYear() &&
      date.getMonth() === now.getMonth() &&
      date.getDate() === now.getDate()

    if (isFetchingForToday) {
      slots = slots.filter((slot: TimeSlot) => new Date(slot.start).getTime() > now.getTime())
    }

    availableSlots.value = slots
  } catch (e: any) {
    console.error('Failed to fetch availability:', e)
    error.value = 'Could not load available time slots. Please try again later.'
    availableSlots.value = []
  } finally {
    isLoading.value = false
  }
}

function selectSlot(slot: TimeSlot) {
  selectedSlot.value = slot
}

async function handleBookingSubmit() {
  if (!selectedSlot.value) return
  isLoading.value = true
  error.value = null
  try {
    const measurementId = 'G-W1TJ3KMZ6V'
    const { clientId, sessionId } = getGaSessionInfo(measurementId)
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
    if (!response.ok) {
      const errorData = await response.text()
      throw new Error(errorData || 'Booking failed. Please try again.')
    }
    isBookingConfirmed.value = true
  } catch (e: any) {
    error.value = e.message || 'An unexpected error occurred during booking.'
  } finally {
    isLoading.value = false
  }
}

function resetBookingProcess() {
  isBookingConfirmed.value = false
  selectedSlot.value = null
  bookingDetails.value = { name: '', email: '', notes: '' }
  fetchAvailability(selectedDate.value)
}

function formatTime(dateString: string) {
  return new Date(dateString).toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit', hour12: true })
}

watch(selectedDate, (newDate) => { fetchAvailability(newDate) })

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
  <section class="py-16 md:py-24 bg-surface-sunken min-h-[70vh]">
  <div class="max-w-2xl mx-auto bg-surface p-8 rounded-2xl shadow-lg">
    <div v-if="!isBookingConfirmed">
      <h1 class="text-2xl font-bold text-dark-slate mb-6 text-center">Book a Consultation</h1>

      <!-- Date Selector — month-grid calendar view -->
      <div class="mb-2">
        <BookingCalendarMonthGrid
          v-model="selectedDate"
          :min-date="today"
          :max-date="maxBookingDate"
          :disabled="isLoading"
        />
      </div>
      <p class="text-xs text-muted text-center mb-6">Timezone: {{ timezone }}</p>

      <!-- Error Message -->
      <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4" role="alert">
        <span class="block sm:inline">{{ error }}</span>
      </div>

      <!-- Loading Spinner -->
      <div v-if="isLoading" class="grid grid-cols-2 sm:grid-cols-3 gap-4">
        <div v-for="n in 6" :key="n" class="p-3 h-12 bg-border rounded-lg animate-pulse"></div>
      </div>

      <!-- Time Slots -->
      <div v-else-if="!selectedSlot" class="min-h-[100px]">
        <h4 class="text-sm font-medium text-muted mb-3 text-center">
          Available slots for {{ formattedDate }}
        </h4>
        <div v-if="availableSlots.length > 0" class="grid grid-cols-2 sm:grid-cols-3 gap-4">
          <button
            v-for="slot in availableSlots"
            :key="slot.start"
            @click="selectSlot(slot)"
            class="p-3 border rounded-lg text-center text-primary hover:bg-primary hover:text-white transition-colors"
          >
            {{ formatTime(slot.start) }}
          </button>
        </div>
        <div v-else-if="!error" class="flex items-center justify-center bg-surface-sunken border border-border text-muted px-4 py-3 rounded-lg" role="status">
          <svg class="w-5 h-5 mr-3 text-muted" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"></path>
          </svg>
          <span class="block sm:inline">There are no available slots for this day. Please try another date.</span>
        </div>
      </div>

      <!-- Booking Form -->
      <form v-else @submit.prevent="handleBookingSubmit">
        <h4 class="text-xl font-semibold mb-4">Confirming for: {{ formatTime(selectedSlot.start) }}</h4>
        <div class="space-y-4">
          <div>
            <label for="name" class="block text-muted font-medium mb-1">Full Name</label>
            <input type="text" id="name" v-model="bookingDetails.name" required class="form-field" />
          </div>
          <div>
            <label for="email" class="block text-muted font-medium mb-1">Email</label>
            <input type="email" id="email" v-model="bookingDetails.email" required class="form-field" />
          </div>
          <div>
            <label for="notes" class="block text-muted font-medium mb-1">Project Notes (Optional)</label>
            <textarea id="notes" v-model="bookingDetails.notes" rows="3" class="form-field"></textarea>
          </div>
        </div>
        <div class="flex justify-end gap-4 mt-6">
          <button type="button" @click="selectedSlot = null" class="px-6 py-2 bg-border text-text rounded-lg hover:bg-border">Back</button>
          <button type="submit" :disabled="isLoading" class="px-6 py-2 bg-primary text-white rounded-lg hover:bg-primary-dark disabled:bg-gray-400 disabled:cursor-not-allowed">
            <span v-if="isLoading">Booking...</span>
            <span v-else>Confirm Booking</span>
          </button>
        </div>
      </form>
    </div>

    <!-- Confirmation Message -->
    <div v-else class="text-center py-12">
      <h2 class="text-2xl font-bold text-green-600 mb-4">Booking Confirmed!</h2>
      <p class="text-muted">
        Thank you, {{ bookingDetails.name }}. A confirmation email has been sent to {{ bookingDetails.email }}.
      </p>
      <div class="flex justify-center gap-4 mt-8">
        <NuxtLink to="/" class="px-6 py-2 bg-border text-text rounded-lg hover:bg-border">Home</NuxtLink>
        <button @click="resetBookingProcess" class="px-6 py-2 bg-primary text-white rounded-lg hover:bg-primary-dark">Book Another</button>
      </div>
    </div>
  </div>
  </section>
</template>
