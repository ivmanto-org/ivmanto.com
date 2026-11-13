package gcal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
	"ivmanto.com/backend/internal/config"
)

var (
	// ErrSlotNotFound is returned when a requested booking slot cannot be found or is already booked.
	ErrSlotNotFound = errors.New("slot not found or already booked")
)

// Service defines the interface for interacting with Google Calendar.
type Service interface {
	GetAvailability(day time.Time) ([]*calendar.Event, error)
	BookSlot(details BookingDetails) (*calendar.Event, error)
	CancelBooking(ctx context.Context, token string) (*calendar.Event, error)
	Location() *time.Location
	// CleanupExpiredBookings deletes held bookings whose call ended
	// before cutoff (booking-data retention; see retention.go).
	CleanupExpiredBookings(ctx context.Context, cutoff time.Time, dryRun bool) ([]string, error)
}

// gcalService implements the Service interface for Google Calendar.
type gcalService struct {
	calSvc               *calendar.Service
	calendarID           string
	location             *time.Location
	availableSlotSummary string

	// testCalls, when non-nil, replaces the real calendar calls in
	// CleanupExpiredBookings (unit tests only; see retention.go).
	testCalls calendarCalls
}

// BookingDetails contains information for a new booking.
type BookingDetails struct {
	EventID         string
	Name            string
	Email           string
	Notes           string
	VisitorTimezone string // IANA zone, e.g. "Europe/Athens"; empty string falls back to the calendar's zone
}

// NewService creates a new calendar service client using Domain-Wide Delegation.
// The runtime principal (ADC) impersonates the configured service account, which in turn
// impersonates a Workspace user via DWD. This is required so that Google Meet conferences
// can be created on the calendar events.
func NewService(ctx context.Context, cfg *config.Config) (Service, error) {
	slog.Info("Authenticating for Google Calendar via DWD user impersonation",
		"sa", cfg.GCal.ServiceAccountEmail,
		"subject", cfg.GCal.ImpersonateUser)

	ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: cfg.GCal.ServiceAccountEmail,
		Scopes:          []string{calendar.CalendarScope},
		Subject:         cfg.GCal.ImpersonateUser,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create impersonated token source: %w", err)
	}

	srv, err := calendar.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Calendar client: %w", err)
	}

	// Fetch calendar details to get its timezone. This is crucial for correctly
	// interpreting date-only queries from the frontend.
	cal, err := srv.Calendars.Get(cfg.GCal.CalendarID).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve calendar details for ID %s: %w", cfg.GCal.CalendarID, err)
	}
	loc, err := time.LoadLocation(cal.TimeZone)
	if err != nil {
		// Fallback to UTC if the location is not found, but log a warning.
		// This can happen in minimal container environments.
		slog.Warn("could not load timezone, falling back to UTC", "timezone", cal.TimeZone, "error", err)
		loc = time.UTC
	}

	return &gcalService{
		calSvc:               srv,
		calendarID:           cfg.GCal.CalendarID,
		location:             loc,
		availableSlotSummary: strings.TrimSpace(cfg.GCal.AvailableSlotSummary),
	}, nil
}

// GetAvailability fetches available time slots for a given day.
// It now lists events with a specific summary to find available slots.
func (s *gcalService) GetAvailability(day time.Time) ([]*calendar.Event, error) {
	loc := s.location
	startOfDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	events, err := s.calSvc.Events.List(s.calendarID).
		TimeMin(startOfDay.Format(time.RFC3339)).
		TimeMax(endOfDay.Format(time.RFC3339)).
		Q(s.availableSlotSummary). // Search for events with the "Available" summary
		SingleEvents(true).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve availability events: %w", err)
	}

	// The original implementation returned []*calendar.TimePeriod, but the handler
	// expects []*calendar.Event. We will return the events directly.
	return events.Items, nil
}

// BookSlot books a consultation by finding an "Available" event and updating it.
// This provides an atomic way to claim a slot.
func (s *gcalService) BookSlot(details BookingDetails) (*calendar.Event, error) {
	slog.Info("Attempting to book event", "eventID", details.EventID)

	// 1. Get the event directly by its unique ID. This is more reliable than searching.
	eventToBook, err := s.calSvc.Events.Get(s.calendarID, details.EventID).Do()
	if err != nil {
		// If the error is 404, it means the event doesn't exist, which we treat as a slot not found.
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
			return nil, ErrSlotNotFound
		}
		return nil, fmt.Errorf("unable to retrieve event to book with ID %s: %w", details.EventID, err)
	}

	slog.Info("Found available event to book", "eventID", eventToBook.Id)

	// 2. Verify the event is indeed an available slot and not already booked.
	// We trim the space from the calendar summary to be robust against accidental whitespace.
	if strings.TrimSpace(eventToBook.Summary) != s.availableSlotSummary {
		slog.Error("Slot verification failed", "eventSummary", strings.TrimSpace(eventToBook.Summary), "expectedSummary", s.availableSlotSummary)
		return nil, ErrSlotNotFound
	}

	// Generate a unique cancellation token for this booking.
	cancellationUUID, err := uuid.NewRandom()
	if err != nil {
		// This is a server-side issue, but we shouldn't fail the whole booking for it.
		// Log it and continue. The user just won't get a cancellation link.
		slog.Warn("could not generate cancellation token UUID", "error", err)
	} else {
		cancellationToken := cancellationUUID.String()
		if eventToBook.ExtendedProperties == nil {
			eventToBook.ExtendedProperties = &calendar.EventExtendedProperties{}
		}
		if eventToBook.ExtendedProperties.Private == nil {
			eventToBook.ExtendedProperties.Private = make(map[string]string)
		}
		// Store booking details for later use (e.g., cancellation notifications).
		eventToBook.ExtendedProperties.Private["cancellation_token"] = cancellationToken
		eventToBook.ExtendedProperties.Private["client_name"] = details.Name
		eventToBook.ExtendedProperties.Private["client_email"] = details.Email
		// VisitorTimezone is the IANA zone of the visitor's browser at booking
		// time. We persist it so the cancellation email can render the slot
		// in the visitor's local time even though the cancellation request
		// comes via an email link with no live client context. The booking
		// handler's resolveVisitorTimezone helper applies the same fallback
		// rules if this string is empty or unknown.
		eventToBook.ExtendedProperties.Private["visitor_timezone"] = details.VisitorTimezone
	}

	// 3. Update the event with the client's details.
	eventToBook.Summary = fmt.Sprintf("Consultation: %s", details.Name)
	eventToBook.Description = fmt.Sprintf(
		"Client Name: %s\nClient Email: %s\n\nNotes:\n%s",
		details.Name,
		details.Email,
		details.Notes,
	)
	// We do not add the client as an attendee directly, as this can require
	// domain-wide delegation. Instead, we send an .ics attachment in the
	// confirmation email. We will leave the existing attendees (i.e., the calendar owner) on the event.
	// eventToBook.Attendees = nil

	// Request Google Meet conference data to be added to the event.
	meetRequestID, err := uuid.NewRandom()
	if err != nil {
		slog.Warn("could not generate Meet request ID, skipping Meet link", "error", err)
	} else {
		eventToBook.ConferenceData = &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				RequestId: meetRequestID.String(),
				ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
					Type: "hangoutsMeet",
				},
			},
		}
	}

	// 4. Atomically update the event. The ETag mechanism handled by the client library
	// ensures that if the event was changed between our read and write, this will fail.
	// ConferenceDataVersion(1) tells the API to process the ConferenceData create request.
	updatedEvent, err := s.calSvc.Events.Update(s.calendarID, eventToBook.Id, eventToBook).ConferenceDataVersion(1).Do()

	if err != nil {
		// Check for a 409 Conflict or 412 Precondition Failed, which indicates the slot was just taken.
		if gerr, ok := err.(*googleapi.Error); ok && (gerr.Code == http.StatusConflict || gerr.Code == http.StatusPreconditionFailed) {
			return nil, ErrSlotNotFound
		}
		// This is a generic error for when the event update fails for reasons other than a conflict.
		return nil, fmt.Errorf("failed to update event during booking: %w", err)
	}
	slog.Info("Successfully updated event with booking details", "eventID", updatedEvent.Id)

	// Return the updated event. The Google Meet link will be in ConferenceData / HangoutLink.
	return updatedEvent, nil
}

// CancelBooking finds an event by its cancellation token and reverts it to an available slot.
// It returns the original event details for notification purposes.
func (s *gcalService) CancelBooking(ctx context.Context, token string) (*calendar.Event, error) {
	slog.Info("Searching for event with cancellation token", "tokenPrefix", token[:8])
	// 1. Find the event using the private extended property.
	query := fmt.Sprintf("cancellation_token=%s", token)
	events, err := s.calSvc.Events.List(s.calendarID).
		PrivateExtendedProperty(query).
		MaxResults(1).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to query for event with token: %w", err)
	}

	if len(events.Items) == 0 {
		slog.Warn("No event found for cancellation token", "tokenPrefix", token[:8])
		return nil, ErrSlotNotFound // Using existing error for "not found"
	}

	eventToCancel := events.Items[0]
	slog.Info("Found event to cancel", "eventID", eventToCancel.Id)

	// 2. Preserve original details for notifications before modifying.
	// We retrieve the client details from the private properties we stored during booking.
	var visitorTZ string
	if eventToCancel.ExtendedProperties != nil && eventToCancel.ExtendedProperties.Private != nil {
		visitorTZ = eventToCancel.ExtendedProperties.Private["visitor_timezone"]
	}
	originalEvent := &calendar.Event{
		Id:      eventToCancel.Id,
		Summary: eventToCancel.Summary,
		Start:   eventToCancel.Start,
		End:     eventToCancel.End,
		Attendees: []*calendar.EventAttendee{
			{
				DisplayName: eventToCancel.ExtendedProperties.Private["client_name"],
				Email:       eventToCancel.ExtendedProperties.Private["client_email"],
			},
		},
		ExtendedProperties: &calendar.EventExtendedProperties{
			Private: map[string]string{
				"visitor_timezone": visitorTZ,
			},
		},
	}

	// 3. Update the event to revert it to an "Available" slot.
	eventToCancel.Summary = s.availableSlotSummary
	eventToCancel.Description = "This slot is now available for booking."
	// Do not modify the attendees list to avoid permission errors trying to remove the calendar owner.
	// eventToCancel.Attendees = nil
	// Conference data modification has been removed to align with the booking logic.
	// The Meet link, if it was ever created, will remain on the reverted event.
	// eventToCancel.ConferenceData = nil
	delete(eventToCancel.ExtendedProperties.Private, "cancellation_token")
	delete(eventToCancel.ExtendedProperties.Private, "client_name")
	delete(eventToCancel.ExtendedProperties.Private, "client_email")
	delete(eventToCancel.ExtendedProperties.Private, "visitor_timezone")

	// 4. Persist the update to Google Calendar.
	_, err = s.calSvc.Events.Update(s.calendarID, eventToCancel.Id, eventToCancel).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to update event to available: %w", err)
	}
	slog.Info("Successfully reverted event to an available slot", "eventID", eventToCancel.Id)

	return originalEvent, nil
}

// Location returns the timezone of the calendar.
func (s *gcalService) Location() *time.Location {
	return s.location
}
