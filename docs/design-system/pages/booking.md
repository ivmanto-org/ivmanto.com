# Book a call (v1.5)

Preview: `booking.html` (the calendar works; times are made up). Styles: `.booking`, `.facts`, `.booking__panel`, `.cal`, `.slots`, `.slot`, `.picked`, `.done-mark` in `components/ivmanto.css`.
Routes stay: `/booking` and `/booking/cancel?token=…`. Booking logic, API calls and emails do not change; this is layout, copy and formats.

## What changes for the visitor
| Today | v1.5 |
|---|---|
| A calendar with the title "Book a Consultation" | Left: what the call is. Right: the booking panel |
| No word on length, format or who is on the call | Three facts: 30 minutes, video call, with Nikolay Tonev |
| Week starts on Sunday, times as "1:30 PM" | Week starts on Monday, times as "13:30" |
| "Friday, November 13, 2026" | "Friday 13 November 2026" |
| No sign of where you are in the flow | Three steps: Pick a time, Your details, Confirmed (`OrderStepper`) |

## Layout
Two columns from 960px (`1fr / 1.25fr`), one column below with the text first. Ground `--bg`, panel on `--surface` with a 1px `--border-input` border, no shadow.

### Left column
- Overline "Book a call". h1 "30 minutes to find the first step".
- Lead: "Tell me where you stand with AI. You leave the call knowing which service fits, or that you need none."
- Facts (`.facts`, icon + title + one line):
  1. **30 minutes** · `[Free and without obligation: Nick to confirm]`
  2. **Video call** · "The link comes with the confirmation email." `[Clauco: confirm the email carries a link, and which tool.]`
  3. **With Nikolay Tonev** · `[In English or German: Nick to confirm]`
- "Prefer to write? Send a message" (link to `/contact`).

### Panel, step 1: Pick a time
- Month head: previous and next as 44px icon buttons (`aria-label` "Previous month", "Next month"), month and year between them (`aria-live="polite"`). Previous is disabled in the current month, Next in the last month of the 90-day window.
- Weekday row Mon to Sun, JetBrains Mono 12px uppercase.
- Day cells: buttons 44px high. States: free to pick (`--text`); not bookable (past, or beyond the window: muted, `disabled`); today (1px `--border-input` outline, "today" in the `aria-label`); chosen (filled `--accent`, `aria-pressed="true"`). Each button's `aria-label` is the full date ("Monday 16 November").
- Line under the grid: "Times are shown in your time zone: Europe/Berlin".
- "Free times on Monday 16 November", then the times as `.slot` buttons (44px high, 24-hour format from `Intl` with `hour12: false`).
- States: loading (three grey placeholder slots, plus "Loading the times…" for screen readers); none ("No free times on this day. Please pick another date."); error (`.notice--error`: "That didn't work. Please try again, or send a message instead.").

### Panel, step 2: Your details
- Summary bar (`.picked`): "Monday 16 November 2026, 12:30 · 30 minutes" and a "Change" link back to step 1 (the chosen day stays).
- Fields (`.field`): **Name** (required), **Work email** (required; hint "The confirmation and the calendar file go here."), **What would you like to talk about?** (optional, 3 rows). A topic passed in from the advisor still pre-fills the last field.
- Note under the fields: "I use these details only to hold and prepare this call (Art. 6(1)(b) GDPR). Privacy notice" (link).
- Button "Book the call" (primary, full width); while sending "Booking…", same width.
- Errors: field errors under the field (`.field--error`); server error as `.notice--error` above the button; slot taken in the meantime: "This time has just been taken. Please pick another one." and back to step 1.
- Focus moves to the first field when the step opens.

### Panel, step 3: Confirmed
- Green tick, h2 "Your call is booked" (receives focus), the date line, then: "A confirmation with a calendar file is on its way to your inbox. It also has a link to cancel."
- Buttons: "Back to home" (secondary), "Read the insights" (ghost). "Book Another" is gone.

## Cancel page (`/booking/cancel`)
Same ground, one centred panel (max 560px).
- Working: "Cancelling your call…"
- Done: h1 "Your call is cancelled", "The time is free again. You can book a new one whenever it suits you.", button "Book a new time" (primary).
- Link invalid or used: h1 "This cancel link no longer works", "The call may already be cancelled. If not, reply to the confirmation email and I'll cancel it for you.", button "Book a call" (secondary).

## Data protection
- Legal basis for name, email and notes: Art. 6(1)(b) GDPR (steps before a contract, at the visitor's request). No consent box is needed and none is shown; the details are not used for marketing. A newsletter or follow-up sequence would need its own consent and is not part of this page.
- The booking request sends the Google Analytics client and session id. **Clauco to confirm** these are sent only after the visitor accepted analytics cookies; without consent they must be empty.
- The privacy policy must name the booking data, the calendar provider and the retention time. `[Retention time: Nick. Check by Ivgil with the privacy policy.]`

## Checks before it ships
Tab order: month buttons, days, times, then the form. Every target 44px high. At 390px the seven day columns fit without sideways scroll (the panel padding drops to 16px). Reduced motion: nothing here moves.
