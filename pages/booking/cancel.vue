<script setup lang="ts">
// Cancel page (DS v1.5, pages/booking.md): one centred panel, three states.
useSeoMeta({
  title: 'Cancel your call | ivmanto.com',
  description: 'Cancel a call you booked with ivmanto.',
  robots: 'noindex, nofollow',
})

const route = useRoute()
const isLoading = ref(true)
const hasFailed = ref(false)

onMounted(async () => {
  const token = route.query.token
  if (!token || typeof token !== 'string') {
    hasFailed.value = true
    isLoading.value = false
    return
  }
  try {
    const response = await fetch('/api/booking/cancel', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ token }),
    })
    if (!response.ok) throw new Error(`Server responded with status ${response.status}`)
  } catch (e) {
    console.error('Cancelling failed:', e)
    hasFailed.value = true
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <section class="section" style="padding-top: var(--space-12)">
    <div class="container">
      <div class="booking__panel cancel-panel">
        <p v-if="isLoading" class="muted" role="status" style="margin: 0">Cancelling your call…</p>
        <template v-else-if="hasFailed">
          <h1 class="h3">This cancel link no longer works</h1>
          <p class="muted">The call may already be cancelled. If not, reply to the confirmation email and I’ll cancel it for you.</p>
          <NuxtLink to="/booking" class="btn btn--secondary">Book a call</NuxtLink>
        </template>
        <template v-else>
          <h1 class="h3">Your call is cancelled</h1>
          <p class="muted">The time is free again. You can book a new one whenever it suits you.</p>
          <NuxtLink to="/booking" class="btn btn--primary">Book a new time</NuxtLink>
        </template>
      </div>
    </div>
  </section>
</template>

<style scoped>
.cancel-panel { max-width: 560px; margin: 0 auto; }
</style>
