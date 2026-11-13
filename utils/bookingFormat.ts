// Date and time formats of the booking page (DS v1.5, pages/booking.md):
// "Monday 16 November 2026", "Monday 16 November", "13:30", weeks from Monday.

const WEEKDAYS = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
const MONTHS = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December']

export const WEEKDAY_LABELS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'] as const
export const CALL_MINUTES = 30

/** "Monday 16 November" in the visitor's local time; with the year on request. */
export function formatLongDate(date: Date, withYear = false): string {
  const base = `${WEEKDAYS[date.getDay()]} ${date.getDate()} ${MONTHS[date.getMonth()]}`
  return withYear ? `${base} ${date.getFullYear()}` : base
}

export function formatMonth(date: Date): string {
  return `${MONTHS[date.getMonth()]} ${date.getFullYear()}`
}

/** "13:30" in the visitor's time zone (or the given one). h23 keeps midnight at 00:00. */
export function formatTime24(value: string | Date, timeZone?: string): string {
  const date = typeof value === 'string' ? new Date(value) : value
  if (Number.isNaN(date.getTime())) return ''
  return new Intl.DateTimeFormat('en-GB', { hour: '2-digit', minute: '2-digit', hourCycle: 'h23', timeZone }).format(date)
}

/** Month grid from the Monday on or before the 1st; whole weeks, no week that lies fully outside the month. */
export function buildMondayGrid(displayMonth: Date): Date[] {
  const year = displayMonth.getFullYear()
  const month = displayMonth.getMonth()
  const lead = (new Date(year, month, 1).getDay() + 6) % 7
  const days = new Date(year, month + 1, 0).getDate()
  const cells = Math.ceil((lead + days) / 7) * 7
  return Array.from({ length: cells }, (_, i) => new Date(year, month, 1 - lead + i))
}
