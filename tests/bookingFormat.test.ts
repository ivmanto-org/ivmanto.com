import { test } from 'node:test'
import assert from 'node:assert/strict'
import { buildMondayGrid, formatLongDate, formatMonth, formatTime24 } from '../utils/bookingFormat.ts'

test('long dates read "Monday 16 November 2026"', () => {
  const date = new Date(2026, 10, 16)
  assert.equal(formatLongDate(date, true), 'Monday 16 November 2026')
  assert.equal(formatLongDate(date), 'Monday 16 November')
  assert.equal(formatMonth(date), 'November 2026')
})

test('times are 24-hour, midnight is 00:00', () => {
  assert.equal(formatTime24('2026-11-16T12:30:00Z', 'Europe/Berlin'), '13:30')
  assert.equal(formatTime24('2026-11-16T23:00:00Z', 'Europe/Berlin'), '00:00')
  assert.equal(formatTime24('2026-11-16T07:05:00Z', 'UTC'), '07:05')
  assert.equal(formatTime24('nonsense'), '')
})

test('the month grid starts on Monday and holds whole weeks', () => {
  const nov = buildMondayGrid(new Date(2026, 10, 1)) // 1 Nov 2026 is a Sunday
  assert.equal(nov.length, 42)
  assert.equal(nov[0]!.getDay(), 1)
  assert.equal(nov[0]!.getDate(), 26)
  assert.equal(nov[6]!.getDate(), 1)
  const feb = buildMondayGrid(new Date(2027, 1, 1)) // 1 Feb 2027 is a Monday, 28 days
  assert.equal(feb.length, 28)
  assert.equal(feb[0]!.getDate(), 1)
  const jun = buildMondayGrid(new Date(2026, 5, 15)) // 1 Jun 2026 is a Monday
  assert.equal(jun.length, 35)
  assert.ok(jun.every((d, i) => d.getDay() === (i + 1) % 7))
})
