<template>
  <div class="container mx-auto px-4 py-16 text-center">
    <div v-if="isLoading" class="prose lg:prose-xl">
      <h1 class="text-2xl font-bold">Cancelling Your Booking...</h1>
      <p>Please wait while we process your request.</p>
    </div>

    <div v-else-if="error" class="prose lg:prose-xl text-red-600">
      <h1 class="text-2xl font-bold">Cancellation Failed</h1>
      <p>{{ error }}</p>
      <p>
        Please
        <NuxtLink to="/booking" class="text-muted hover:underline">contact us</NuxtLink>
        if you need further assistance.
      </p>
    </div>

    <div v-else class="prose lg:prose-xl">
      <h1 class="text-2xl font-bold text-green-600">Booking Cancelled Successfully</h1>
      <p>Your consultation has been cancelled. You will receive a confirmation email shortly.</p>
      <p>
        If you'd like to schedule a new time, feel free to visit our
        <NuxtLink to="/booking" class="text-muted hover:underline">booking page</NuxtLink>
        again.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
useSeoMeta({
  title: 'Cancel Booking | ivmanto.com',
  description: 'Cancel your scheduled consultation with IVMANTO. Contact us if you need further assistance modifying your appointment.',
  robots: 'noindex, nofollow',
})

const route = useRoute()
const isLoading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  const token = route.query.token

  if (!token) {
    error.value = 'Invalid or missing cancellation token in the URL.'
    isLoading.value = false
    return
  }

  try {
    const response = await fetch('/api/booking/cancel', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ token }),
    })

    if (!response.ok) {
      const errData = await response.json()
      throw new Error(errData.message || 'The server returned an error.')
    }
  } catch (e: any) {
    error.value = e.message || 'An unexpected error occurred.'
  } finally {
    isLoading.value = false
  }
})
</script>
