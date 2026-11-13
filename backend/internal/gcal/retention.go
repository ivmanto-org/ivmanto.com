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

// CleanupExpiredBookings deletes every held booking whose call ended
// before cutoff and returns the deleted event IDs. IDs only are returned
// and logged — never names, emails or notes. In dry-run mode nothing is
// deleted; the events that would be deleted are logged and returned.
func (s *gcalService) CleanupExpiredBookings(ctx context.Context, cutoff time.Time, dryRun bool) ([]string, error) {
	// List window: calls are 30 minutes, so a 7-day look-back before the
	// cutoff covers any conceivable event duration overlap. Events ending
	// after the cutoff are filtered out by IsExpiredBookedEvent anyway;
	// the window only bounds the API result set.
	timeMin := cutoff.Add(-7 * 24 * time.Hour).UTC().Format(time.RFC3339)

	var handled []string
	pageToken := ""
	for {
		call := s.calSvc.Events.List(s.calendarID).
			TimeMin(timeMin).
			MaxResults(250)
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}
		events, err := call.Do()
		if err != nil {
			return handled, fmt.Errorf("failed to list calendar events for retention cleanup: %w", err)
		}

		for _, ev := range events.Items {
			if !IsExpiredBookedEvent(ev, cutoff) {
				continue
			}
			if dryRun {
				slog.Info("dry-run: would delete expired booking event", "event_id", ev.Id, "action", "delete_skipped")
			} else {
				if err := s.calSvc.Events.Delete(s.calendarID, ev.Id).Do(); err != nil {
					// One stuck event must not abort the whole run;
					// log (id only) and keep deleting the rest.
					slog.Error("failed to delete expired booking event", "event_id", ev.Id, "error", err)
					continue
				}
				slog.Info("deleted expired booking event", "event_id", ev.Id)
			}
			handled = append(handled, ev.Id)
		}

		if events.NextPageToken == "" {
			break
		}
		pageToken = events.NextPageToken
	}
	return handled, nil
}
