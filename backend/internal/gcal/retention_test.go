package gcal

import (
	"testing"
	"time"

	"google.golang.org/api/calendar/v3"
)

func bookedEvent(id, end string, token, email string) *calendar.Event {
	ev := &calendar.Event{
		Id:      id,
		End:     &calendar.EventDateTime{DateTime: end},
		Summary: "Consultation: Jane Doe",
	}
	if token != "" || email != "" {
		ev.ExtendedProperties = &calendar.EventExtendedProperties{
			Private: map[string]string{
				"cancellation_token": token,
				"client_email":       email,
			},
		}
	}
	return ev
}

// The canonical scenario: a call that ended more than 10 days ago is
// due for deletion.
func TestIsExpiredBookedEvent_ExpiredCallSelected(t *testing.T) {
	cutoff := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	ev := bookedEvent("evt-expired", "2026-11-01T13:00:00Z", "tok-1", "jane@example.com")
	if !IsExpiredBookedEvent(ev, cutoff) {
		t.Error("held booking whose call ended 12 days before cutoff must be selected")
	}
}

// Nick's edge case (2026-11-13): a booking made today for a call 3 weeks
// out must NOT be deleted — the 10 days count from the call, not from the
// booking. This is the guard that makes a "10 days after booking" bug
// impossible to reintroduce silently.
func TestIsExpiredBookedEvent_FutureCallNeverSelected(t *testing.T) {
	cutoff := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	// Booked today (props written now), call 20 days in the future.
	ev := bookedEvent("evt-future-call", "2026-12-03T13:00:00Z", "tok-2", "jane@example.com")
	if IsExpiredBookedEvent(ev, cutoff) {
		t.Error("a booking whose call has not happened yet must never be selected for deletion")
	}
}

// A call that ended exactly at the cutoff (end + 10 days) IS expired —
// "kept 10 days, then deleted" — while one ending a second later is not.
func TestIsExpiredBookedEvent_BoundaryIsExact(t *testing.T) {
	// now = 2026-11-13 12:00 UTC, so cutoff = now - 10 days = Nov 3 12:00.
	cutoff := time.Date(2026, 11, 3, 12, 0, 0, 0, time.UTC)
	ev := bookedEvent("evt-boundary", "2026-11-03T12:00:00Z", "tok-3", "jane@example.com")
	if !IsExpiredBookedEvent(ev, cutoff) {
		t.Error("end == cutoff (exactly 10 days) must be selected")
	}
	ev2 := bookedEvent("evt-just-after", "2026-11-03T12:00:01Z", "tok-4", "jane@example.com")
	if IsExpiredBookedEvent(ev2, cutoff) {
		t.Error("end one second after cutoff must not be selected; 10 full days have not passed")
	}
}

// An available (never booked, or cancelled-and-reverted) slot carries no
// booking PII and must never be deleted.
func TestIsExpiredBookedEvent_AvailableSlotNotSelected(t *testing.T) {
	cutoff := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	slot := &calendar.Event{
		Id:      "evt-available",
		Summary: "Available",
		End:     &calendar.EventDateTime{DateTime: "2026-10-01T13:00:00Z"},
	}
	if IsExpiredBookedEvent(slot, cutoff) {
		t.Error("available slot must never be selected")
	}
	// A cancelled booking: CancelBooking deleted the private props, so
	// the reverted slot must not be selected even though the event still
	// exists on the calendar.
	cancelled := &calendar.Event{
		Id:                 "evt-cancelled",
		Summary:            "Available",
		Description:        "This slot is now available for booking.",
		End:                &calendar.EventDateTime{DateTime: "2026-10-01T13:00:00Z"},
		ExtendedProperties: &calendar.EventExtendedProperties{Private: map[string]string{}},
	}
	if IsExpiredBookedEvent(cancelled, cutoff) {
		t.Error("cancelled booking (PII wiped at cancellation) must not be selected")
	}
}

// A booking whose cancellation-token UUID generation failed still carries
// client_email, and must still be selected for deletion.
func TestIsExpiredBookedEvent_EmailOnlyMarker(t *testing.T) {
	cutoff := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	ev := bookedEvent("evt-email-only", "2026-11-01T13:00:00Z", "", "jane@example.com")
	if !IsExpiredBookedEvent(ev, cutoff) {
		t.Error("booking carrying client_email (token generation failed) must be selected")
	}
}

// Fail-closed: unparsable, missing end time or nil event never deletes.
func TestIsExpiredBookedEvent_FailClosed(t *testing.T) {
	cutoff := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	cases := map[string]*calendar.Event{
		"nil event":        nil,
		"missing end":      bookedEvent("e1", "", "tok", "jane@example.com"),
		"unparsable end":   bookedEvent("e2", "not-a-time", "tok", "jane@example.com"),
		"empty marker key": {Id: "e3", End: &calendar.EventDateTime{DateTime: "2026-11-01T13:00:00Z"}, ExtendedProperties: &calendar.EventExtendedProperties{Private: map[string]string{"cancellation_token": ""}}},
	}
	for name, ev := range cases {
		if IsExpiredBookedEvent(ev, cutoff) {
			t.Errorf("%s: must not be selected (fail-closed)", name)
		}
	}
}
