<script setup lang="ts">
// Month calendar of the booking panel (DS v1.5, pages/booking.md): weeks from Monday,
// 44px day buttons, the full date as each button's label. The parent owns the chosen day;
// this component only holds the month on display.
import { isSameDay, clampDate } from '~/composables/useDate'
import { WEEKDAY_LABELS, buildMondayGrid, formatLongDate, formatMonth } from '~/utils/bookingFormat'

const DEFAULT_WINDOW_DAYS = 30

const props = withDefaults(
  defineProps<{
    modelValue: Date
    minDate?: Date
    maxDate?: Date
    disabled?: boolean
  }>(),
  { minDate: undefined, maxDate: undefined, disabled: false },
)

const emit = defineEmits<{ 'update:modelValue': [Date] }>()

// Captured once, so "today" and the bounds do not drift while the page stays open.
const now = new Date()
const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
const minDate = computed(() => props.minDate ?? today)
const maxDate = computed(
  () => props.maxDate ?? new Date(today.getFullYear(), today.getMonth(), today.getDate() + DEFAULT_WINDOW_DAYS),
)

const displayMonth = ref(new Date(props.modelValue.getFullYear(), props.modelValue.getMonth(), 1))

// If the bound value jumps to another month (parent reset), follow it.
watch(
  () => props.modelValue,
  (v) => {
    if (v.getFullYear() !== displayMonth.value.getFullYear() || v.getMonth() !== displayMonth.value.getMonth()) {
      displayMonth.value = new Date(v.getFullYear(), v.getMonth(), 1)
    }
  },
)

const grid = computed(() => buildMondayGrid(displayMonth.value))

const monthIndex = (d: Date) => d.getFullYear() * 12 + d.getMonth()
const isPrevMonthDisabled = computed(() => monthIndex(displayMonth.value) <= monthIndex(minDate.value))
const isNextMonthDisabled = computed(() => monthIndex(displayMonth.value) >= monthIndex(maxDate.value))

function inCurrentMonth(d: Date): boolean {
  return d.getMonth() === displayMonth.value.getMonth()
}

function isDisabledCell(d: Date): boolean {
  return (!isSameDay(d, minDate.value) && d.getTime() < minDate.value.getTime()) || d.getTime() > maxDate.value.getTime()
}

function dayLabel(d: Date): string {
  return isSameDay(d, today) ? `${formatLongDate(d)}, today` : formatLongDate(d)
}

function selectDay(d: Date) {
  if (props.disabled || isDisabledCell(d)) return
  emit('update:modelValue', clampDate(d, minDate.value, maxDate.value))
}

function shiftMonth(amount: number) {
  if (amount < 0 && isPrevMonthDisabled.value) return
  if (amount > 0 && isNextMonthDisabled.value) return
  displayMonth.value = new Date(displayMonth.value.getFullYear(), displayMonth.value.getMonth() + amount, 1)
}
</script>

<template>
  <div class="cal">
    <div class="cal__head">
      <button class="icon-btn" type="button" :disabled="isPrevMonthDisabled" aria-label="Previous month" @click="shiftMonth(-1)">
        <svg viewBox="0 0 16 16" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M10 3L5 8l5 5" /></svg>
      </button>
      <h2 aria-live="polite">{{ formatMonth(displayMonth) }}</h2>
      <button class="icon-btn" type="button" :disabled="isNextMonthDisabled" aria-label="Next month" @click="shiftMonth(1)">
        <svg viewBox="0 0 16 16" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M6 3l5 5-5 5" /></svg>
      </button>
    </div>
    <div class="cal__grid">
      <div v-for="d in WEEKDAY_LABELS" :key="d" class="cal__dow" aria-hidden="true">{{ d }}</div>
      <template v-for="cell in grid" :key="cell.getTime()">
        <button
          v-if="inCurrentMonth(cell)"
          type="button"
          class="cal__day"
          :class="{ 'is-today': isSameDay(cell, today) }"
          :disabled="disabled || isDisabledCell(cell)"
          :aria-label="dayLabel(cell)"
          :aria-pressed="isSameDay(cell, modelValue)"
          @click="selectDay(cell)"
        >
          {{ cell.getDate() }}
        </button>
        <span v-else aria-hidden="true"></span>
      </template>
    </div>
  </div>
</template>
