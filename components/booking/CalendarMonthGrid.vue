<script setup lang="ts">
import { buildMonthGrid, isSameDay, clampDate } from '~/composables/useDate'

/**
 * Month-grid date picker used by the booking page. Renders a 6×7 grid of
 * days for the currently displayed month (with out-of-month padding days as
 * inert cells) and lets visitors jump between months and select a day.
 *
 * Selection is two-way bound via v-model:modelValue. The parent's
 * `selectedDate` is the source of truth — this component only holds an
 * internal `displayMonth` ref for the month currently being viewed, which
 * is initialised from `modelValue` and does NOT change `selectedDate` when
 * the visitor navigates months. Only a day click updates the bound date.
 *
 * Bounds: cells outside [minDate, maxDate] are rendered as disabled. The
 * prev-month button is disabled when `displayMonth` is the current month;
 * the next-month button is disabled when `displayMonth` is the month
 * containing `maxDate`.
 */

const props = withDefaults(
  defineProps<{
    modelValue: Date
    minDate?: Date
    maxDate?: Date
    disabled?: boolean
  }>(),
  {
    minDate: undefined,
    maxDate: undefined,
    disabled: false,
  },
)

const emit = defineEmits<{
  'update:modelValue': [Date]
}>()

// Effective bounds: today and today+30 by default. The parent normally
// passes these explicitly; the defaults are here so the component can be
// rendered standalone (Storybook / preview) without surprises.
const today = new Date()
const minDate = computed(() => props.minDate ?? new Date(today.getFullYear(), today.getMonth(), today.getDate()))
const maxDate = computed(() => {
  if (props.maxDate) return props.maxDate
  const d = new Date(today.getFullYear(), today.getMonth(), today.getDate())
  d.setDate(d.getDate() + 30)
  return d
})

// Internal month being shown. Initialised from the bound value. Prev/next
// arrows only mutate this ref, never modelValue.
const displayMonth = ref(new Date(props.modelValue.getFullYear(), props.modelValue.getMonth(), 1))

// Today (midnight) — used for the "today" ring on the cell. Captured once at
// component setup so the ring stays put as the wall clock advances during a
// long session.
const todayMidnight = new Date(today.getFullYear(), today.getMonth(), today.getDate())

// If the bound value jumps to a different month (e.g. parent reset), realign.
watch(
  () => props.modelValue,
  (v) => {
    if (v.getFullYear() !== displayMonth.value.getFullYear() || v.getMonth() !== displayMonth.value.getMonth()) {
      displayMonth.value = new Date(v.getFullYear(), v.getMonth(), 1)
    }
  },
)

const monthLabel = computed(() =>
  displayMonth.value.toLocaleDateString('en-US', { month: 'long', year: 'numeric' }),
)

const grid = computed(() => buildMonthGrid(displayMonth.value))

const isPrevMonthDisabled = computed(() => {
  const firstOfDisplay = new Date(displayMonth.value.getFullYear(), displayMonth.value.getMonth(), 1)
  const firstOfMinMonth = new Date(minDate.value.getFullYear(), minDate.value.getMonth(), 1)
  return firstOfDisplay.getTime() <= firstOfMinMonth.getTime()
})

const isNextMonthDisabled = computed(() => {
  const d = displayMonth.value
  const m = maxDate.value
  return d.getFullYear() > m.getFullYear()
    || (d.getFullYear() === m.getFullYear() && d.getMonth() >= m.getMonth())
})

const dayLabels = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']

function inCurrentMonth(d: Date): boolean {
  return d.getMonth() === displayMonth.value.getMonth()
}

function isDisabledCell(d: Date): boolean {
  return !isSameDay(d, minDate.value) && d.getTime() < minDate.value.getTime()
    || d.getTime() > maxDate.value.getTime()
}

function selectDay(d: Date) {
  if (props.disabled) return
  if (isDisabledCell(d)) return
  // Clamp to bounds just in case (defensive — UI should already enforce).
  const clamped = clampDate(d, minDate.value, maxDate.value)
  emit('update:modelValue', clamped)
}

function shiftMonth(amount: number) {
  if (amount < 0 && isPrevMonthDisabled.value) return
  if (amount > 0 && isNextMonthDisabled.value) return
  const next = new Date(displayMonth.value.getFullYear(), displayMonth.value.getMonth() + amount, 1)
  displayMonth.value = next
}
</script>

<template>
  <div class="w-full" :class="{ 'opacity-60 pointer-events-none': disabled }">
    <!-- Month header + nav -->
    <div class="flex items-center justify-between mb-3">
      <button
        type="button"
        @click="shiftMonth(-1)"
        class="btn-day-nav"
        :disabled="isPrevMonthDisabled"
        aria-label="Previous month"
      >
        <svg class="w-6 h-6 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
        </svg>
      </button>
      <h3 class="text-lg font-semibold text-muted">{{ monthLabel }}</h3>
      <button
        type="button"
        @click="shiftMonth(1)"
        class="btn-day-nav"
        :disabled="isNextMonthDisabled"
        aria-label="Next month"
      >
        <svg class="w-6 h-6 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
        </svg>
      </button>
    </div>

    <!-- Day-of-week header -->
    <div class="grid grid-cols-7 gap-1 mb-1">
      <div v-for="d in dayLabels" :key="d" class="text-center text-xs font-medium text-muted">
        {{ d }}
      </div>
    </div>

    <!-- Day grid -->
    <div class="grid grid-cols-7 gap-1">
      <template v-for="(cell, idx) in grid" :key="idx">
        <button
          v-if="inCurrentMonth(cell)"
          type="button"
          @click="selectDay(cell)"
          :disabled="isDisabledCell(cell)"
          :class="[
            'aspect-square rounded-lg text-sm flex items-center justify-center transition-colors',
            isDisabledCell(cell)
              ? 'text-gray-300 cursor-not-allowed'
              : isSameDay(cell, modelValue)
                ? 'bg-primary text-white font-semibold'
                : isSameDay(cell, todayMidnight) && !isSameDay(cell, modelValue)
                  ? 'border border-primary text-primary font-semibold hover:bg-primary hover:text-white cursor-pointer'
                  : 'text-muted hover:bg-primary hover:text-white cursor-pointer',
          ]"
          :aria-label="cell.toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })"
          :aria-pressed="isSameDay(cell, modelValue)"
        >
          {{ cell.getDate() }}
        </button>
        <div
          v-else
          class="aspect-square rounded-lg"
          aria-hidden="true"
        ></div>
      </template>
    </div>
  </div>
</template>
