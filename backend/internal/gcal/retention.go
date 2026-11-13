package gcal

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/api/calendar/v3"
)

// Retention deletion for booking data (Nick, 2026-11-13): booking data
// from /booking — name, work email, notes — is kept 10 days and then
// deleted. The Google Calendar event IS the booking record; there is no
// database copy (see BookSlot: name/email/notes go into the event
// summary, description and private extended properties). Deleting the
// event therefore deletes the record.
//
// Cancellations need no retention run: CancelBooking already wipes the
// client PII in place (summary and description are reverted, and
// cancellation_token, client_name, client_email and visitor_timezone are
// deleted from the private properties).
//
// Events.Delete is not erasure by itself: a deleted event sits in the
// calendar's trash, restorable with all its data, for 30 days (Google
// Calendar Help, "Delete an event", support.google.com/calendar/answer/
// 37113). The cleanup therefore wipes the personal data IN PLACE first —
// exactly as CancelBooking does — and only then deletes, so the trash
// copy carries no name, email or notes either.

// BookedEventMarkerProps are the private extended-property keys that only
// a held booking carries. Available slots never set them, and
// CancelBooking deletes both, so the presence of either is the reliable
// discriminator of an event that still holds booking PII. client_email is
// a marker too because a booking whose cancellation-token UUID generation
// failed (logged Warn, booking continues) still stores client_email.
func BookedEventMarkerProps() []string {
	return []string{"cancellation_token", "client_email"}
}

// isBookedEvent reports whether the event still carries booking PII.
func isBookedEvent(ev *calendar.Event) bool {
	if ev == nil || ev.ExtendedProperties == nil || ev.ExtendedProperties.Private == nil {
		return false
	}
	for _, key := range BookedEventMarkerProps() {
		if v, ok := ev.ExtendedProperties.Private[key]; ok && v != "" {
			return true
		}
	}
	return false
}

// IsExpiredBookedEvent reports whether ev is a held booking whose call
// ended before cutoff, i.e. it is due for retention deletion. Pure
// function; unit-tested in retention_test.go. The comparison is against
// the call's END time — deleting N days after the booking was MADE could
// remove a call that has not happened yet (booking window is 90 days).
// Fail-closed: an unparsable or missing end time never selects an event.
func IsExpiredBookedEvent(ev *calendar.Event, cutoff time.Time) bool {
	if !isBookedEvent(ev) {
		return false
	}
	if ev.End == nil || ev.End.DateTime == "" {
		return false
	}
	end, err := time.Parse(time.RFC3339, ev.End.DateTime)
	if err != nil {
		return false
	}
	return !end.After(cutoff)
}

// calendarCalls abstracts the three Google Calendar operations the
// retention cleanup needs, so the deletion loop is unit-testable behind
// a fake instead of the real API.
type calendarCalls interface {
	// ListExpiredWindow lists all events in [timeMin, timeMax]
	// (RFC3339, both bounds inclusive per the Calendar API).
	ListEvents(ctx context.Context, timeMin, timeMax string) ([]*calendar.Event, error)
	// WipeEventPII overwrites the event's personal data in place,
	// mirroring CancelBooking (summary, description, private props).
	WipeEventPII(ctx context.Context, ev *calendar.Event) error
	// DeleteEvent deletes the event from the calendar.
	DeleteEvent(ctx context.Context, eventID string) error
}

// prodCalendarCalls is the real implementation on top of the calendar
// service client.
type prodCalendarCalls struct {
	s *gcalService
}

func (p *prodCalendarCalls) ListEvents(ctx context.Context, timeMin, timeMax string) ([]*calendar.Event, error) {
	var events []*calendar.Event
	pageToken := ""
	for {
		call := p.s.calSvc.Events.List(p.s.calendarID).
			TimeMin(timeMin).
			TimeMax(timeMax).
			SingleEvents(true).
			MaxResults(250)
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}
		page, err := call.Do()
		if err != nil {
			return events, fmt.Errorf("failed to list calendar events for retention cleanup: %w", err)
		}
		events = append(events, page.Items...)
		if page.NextPageToken == "" {
			break
		}
		pageToken = page.NextPageToken
	}
	return events, nil
}

func (p *prodCalendarCalls) WipeEventPII(ctx context.Context, ev *calendar.Event) error {
	// Mirror CancelBooking's wipe so the trash copy is clean too: the
	// event must read as a plain available slot before it is deleted.
	// Attendees are untouched — BookSlot never stores the client as an
	// attendee (see calendar.go), so the list carries no client PII.
	ev.Summary = p.s.availableSlotSummary
	ev.Description = "This slot is now available for booking."
	if ev.ExtendedProperties == nil {
		ev.ExtendedProperties = &calendar.EventExtendedProperties{}
	}
	if ev.ExtendedProperties.Private == nil {
		ev.ExtendedProperties.Private = map[string]string{}
	}
	delete(ev.ExtendedProperties.Private, "cancellation_token")
	delete(ev.ExtendedProperties.Private, "client_name")
	delete(ev.ExtendedProperties.Private, "client_email")
	delete(ev.ExtendedProperties.Private, "visitor_timezone")
	_, err := p.s.calSvc.Events.Update(p.s.calendarID, ev.Id, ev).Do()
	if err != nil {
		return fmt.Errorf("failed to wipe personal data from event %s before deletion: %w", ev.Id, err)
	}
	return nil
}

func (p *prodCalendarCalls) DeleteEvent(ctx context.Context, eventID string) error {
	if err := p.s.calSvc.Events.Delete(p.s.calendarID, eventID).Do(); err != nil {
		return fmt.Errorf("failed to delete expired booking event %s: %w", eventID, err)
	}
	return nil
}

// cleanupListLookback bounds how far back the retention run lists events.
// It must NOT be a small window: TimeMin is a lower bound on the event
// END, so a booking whose call ended more than lookback+retention before
// a run would never be listed and never deleted — every pre-existing
// booking before the first run, and anything missed during downtime
// (review finding on PR #140, 2026-11-13). Five years predates the
// calendar's use for bookings entirely; IsExpiredBookedEvent applies the
// exact cutoff anyway.
const cleanupListLookback = 5 * 365 * 24 * time.Hour

// CleanupExpiredBookings deletes every held booking whose call ended
// before cutoff and returns the handled event IDs. IDs only are returned
// and logged — never names, emails or notes. In dry-run mode nothing is
// modified or deleted; the events that would be deleted are logged and
// returned.
//
// Each event's personal data is wiped in place BEFORE deletion, because
// Events.Delete keeps the event restorable in the calendar trash for ~30
// days (see the package comment). A wipe or delete failure does not abort
// the run — the remaining events are still processed — but the run
// returns an error at the end so the job exits non-zero and the failure
// is visible (Cloud Scheduler / job alerts).
func (s *gcalService) CleanupExpiredBookings(ctx context.Context, cutoff time.Time, dryRun bool) ([]string, error) {
	calls := calendarCalls(&prodCalendarCalls{s: s})
	if s.testCalls != nil {
		calls = s.testCalls
	}

	// The list window covers everything from far in the past up to the
	// cutoff. TimeMax keeps future slots out of the result set (the
	// booking window is 90 days ahead); there is no near lower bound.
	timeMin := cutoff.Add(-cleanupListLookback).UTC().Format(time.RFC3339)
	timeMax := cutoff.UTC().Format(time.RFC3339)

	events, err := calls.ListEvents(ctx, timeMin, timeMax)
	if err != nil {
		return nil, err
	}

	var handled []string
	var failed int
	for _, ev := range events {
		if !IsExpiredBookedEvent(ev, cutoff) {
			continue
		}
		if dryRun {
			slog.Info("dry-run: would delete expired booking event", "event_id", ev.Id, "action", "delete_skipped")
			handled = append(handled, ev.Id)
			continue
		}
		// Wipe first: the trash copy must not carry the client PII.
		// A wipe failure skips the delete — deleting an unwiped event
		// would move the PII into the trash for 30 more days.
		if err := calls.WipeEventPII(ctx, ev); err != nil {
			slog.Error("retention cleanup: skipping delete, could not wipe personal data", "event_id", ev.Id, "error", err)
			failed++
			continue
		}
		if err := calls.DeleteEvent(ctx, ev.Id); err != nil {
			// One stuck event must not abort the whole run; log (id
			// only), count it and keep deleting the rest.
			slog.Error("failed to delete expired booking event", "event_id", ev.Id, "error", err)
			failed++
			continue
		}
		slog.Info("deleted expired booking event", "event_id", ev.Id)
		handled = append(handled, ev.Id)
	}

	if failed > 0 {
		return handled, fmt.Errorf("booking retention cleanup: %d event(s) failed to wipe/delete and remain on the calendar", failed)
	}
	return handled, nil
}
