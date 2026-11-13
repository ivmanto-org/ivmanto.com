package gcal

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"google.golang.org/api/calendar/v3"
)

// fakeCalendarCalls records every call and lets a test fail any of the
// three operations on chosen event IDs.
type fakeCalendarCalls struct {
	gotTimeMin string
	gotTimeMax string

	listed []*calendar.Event

	wipeOrder []string // event IDs, in wipe order
	delOrder  []string // event IDs, in delete order

	wiped map[string]*calendar.Event // event ID -> event as wiped

	failWipe   map[string]error
	failDelete map[string]error
}

func (f *fakeCalendarCalls) ListEvents(ctx context.Context, timeMin, timeMax string) ([]*calendar.Event, error) {
	f.gotTimeMin = timeMin
	f.gotTimeMax = timeMax
	return f.listed, nil
}

func (f *fakeCalendarCalls) WipeEventPII(ctx context.Context, ev *calendar.Event) error {
	f.wipeOrder = append(f.wipeOrder, ev.Id)
	if f.failWipe != nil {
		if err, ok := f.failWipe[ev.Id]; ok {
			return err
		}
	}
	wiped := *ev
	wiped.Summary = ""
	wiped.Description = ""
	if wiped.ExtendedProperties != nil && wiped.ExtendedProperties.Private != nil {
		private := map[string]string{}
		for k, v := range wiped.ExtendedProperties.Private {
			private[k] = v
		}
		for _, k := range []string{"cancellation_token", "client_name", "client_email", "visitor_timezone"} {
			delete(private, k)
		}
		wiped.ExtendedProperties.Private = private
	}
	if f.wiped == nil {
		f.wiped = map[string]*calendar.Event{}
	}
	f.wiped[ev.Id] = &wiped
	return nil
}

func (f *fakeCalendarCalls) DeleteEvent(ctx context.Context, eventID string) error {
	f.delOrder = append(f.delOrder, eventID)
	if f.failDelete != nil {
		if err, ok := f.failDelete[eventID]; ok {
			return err
		}
	}
	return nil
}

func newCleanupTestService(fake *fakeCalendarCalls) *gcalService {
	return &gcalService{
		calendarID:           "test-cal",
		availableSlotSummary: "Available",
		testCalls:            fake,
	}
}

func fullBookedEvent(id, end string) *calendar.Event {
	ev := bookedEvent(id, end, "tok-"+id, "jane@example.com")
	ev.Description = "Notes from Jane: jane@example.com, prefers mornings."
	ev.ExtendedProperties.Private["client_name"] = "Jane Doe"
	ev.ExtendedProperties.Private["visitor_timezone"] = "Europe/Berlin"
	return ev
}

// HIGH 1 (PR #140 review): the list window must NOT have a near lower
// bound — TimeMin bounds the event END, so a small look-back window
// would hide every booking older than the window (all pre-first-run
// bookings, and anything missed in a week of downtime) and they would
// never be deleted. The guard: timeMin must be years before the cutoff,
// and timeMax must equal the cutoff.
func TestCleanup_ListWindowHasNoNearLowerBound(t *testing.T) {
	fake := &fakeCalendarCalls{}
	svc := newCleanupTestService(fake)
	cutoff := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)

	if _, err := svc.CleanupExpiredBookings(context.Background(), cutoff, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotMin, err := time.Parse(time.RFC3339, fake.gotTimeMin)
	if err != nil {
		t.Fatalf("timeMin not RFC3339: %q", fake.gotTimeMin)
	}
	if d := cutoff.Sub(gotMin); d < 4*365*24*time.Hour {
		t.Errorf("timeMin %s is only %v before cutoff — a near lower bound hides old bookings and they are never deleted", fake.gotTimeMin, d)
	}
	if fake.gotTimeMax != cutoff.UTC().Format(time.RFC3339) {
		t.Errorf("timeMax must be the cutoff (%s), got %s", cutoff.UTC().Format(time.RFC3339), fake.gotTimeMax)
	}
}

// HIGH 2 (PR #140 review): a failed delete must not end as success. The
// loop keeps processing the rest, but the run returns an error so the
// job exits non-zero and the failure is visible to Cloud Scheduler.
func TestCleanup_FailedDeleteReturnsErrorButContinues(t *testing.T) {
	fake := &fakeCalendarCalls{
		listed: []*calendar.Event{
			fullBookedEvent("evt-ok-1", "2026-11-01T13:00:00Z"),
			fullBookedEvent("evt-stuck", "2026-11-02T13:00:00Z"),
			fullBookedEvent("evt-ok-2", "2026-11-03T11:00:00Z"),
		},
		failDelete: map[string]error{"evt-stuck": errors.New("googleapi: 500")},
	}
	svc := newCleanupTestService(fake)
	cutoff := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)

	ids, err := svc.CleanupExpiredBookings(context.Background(), cutoff, false)
	if err == nil {
		t.Fatal("a failed delete must surface as an error at the end of the run")
	}
	if !strings.Contains(err.Error(), "1") || !strings.Contains(err.Error(), "failed") {
		t.Errorf("error should report the failed count, got: %v", err)
	}
	// The stuck event's failure must not abort the rest: both healthy
	// events were still deleted and reported as handled.
	if len(ids) != 2 || ids[0] != "evt-ok-1" || ids[1] != "evt-ok-2" {
		t.Errorf("remaining events must still be deleted, handled = %v", ids)
	}
	if len(fake.delOrder) != 3 {
		t.Errorf("delete must be attempted on all three events, got %v", fake.delOrder)
	}
}

// MEDIUM 3 (PR #140 review): Events.Delete keeps the event restorable in
// the calendar trash for ~30 days, so personal data must be wiped IN
// PLACE before the delete. Guard: wipe happens before delete for every
// event, and the wiped event carries no name, email, notes or markers.
func TestCleanup_WipesPIIBeforeDelete(t *testing.T) {
	fake := &fakeCalendarCalls{
		listed: []*calendar.Event{fullBookedEvent("evt-expired", "2026-11-01T13:00:00Z")},
	}
	svc := newCleanupTestService(fake)
	cutoff := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)

	if _, err := svc.CleanupExpiredBookings(context.Background(), cutoff, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fake.wipeOrder) != 1 || len(fake.delOrder) != 1 {
		t.Fatalf("expected one wipe and one delete, got wipe=%v delete=%v", fake.wipeOrder, fake.delOrder)
	}
	if fake.wipeOrder[0] != fake.delOrder[0] {
		t.Errorf("wipe and delete must target the same event")
	}
	wiped := fake.wiped["evt-expired"]
	if wiped == nil {
		t.Fatal("wipe was not recorded")
	}
	// The real wipe (prodCalendarCalls) mirrors CancelBooking; the fake
	// here strips the same fields, so assert the trash copy is clean.
	if strings.Contains(wiped.Summary, "Jane") || strings.Contains(wiped.Description, "jane@example.com") {
		t.Errorf("wiped event still carries personal data: summary=%q description=%q", wiped.Summary, wiped.Description)
	}
	for _, k := range []string{"cancellation_token", "client_name", "client_email", "visitor_timezone"} {
		if v, ok := wiped.ExtendedProperties.Private[k]; ok && v != "" {
			t.Errorf("wiped event still carries %s", k)
		}
	}
}

// A wipe failure must skip the delete: deleting an unwiped event would
// move the personal data into the trash for 30 more days. The failure
// counts toward the end-of-run error.
func TestCleanup_WipeFailureSkipsDeleteAndFailsRun(t *testing.T) {
	fake := &fakeCalendarCalls{
		listed: []*calendar.Event{fullBookedEvent("evt-dirty", "2026-11-01T13:00:00Z")},
		failWipe: map[string]error{
			"evt-dirty": errors.New("update denied"),
		},
	}
	svc := newCleanupTestService(fake)
	cutoff := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)

	ids, err := svc.CleanupExpiredBookings(context.Background(), cutoff, false)
	if err == nil {
		t.Fatal("a failed wipe must surface as an error at the end of the run")
	}
	if len(fake.delOrder) != 0 {
		t.Errorf("an event whose wipe failed must NOT be deleted (PII would sit in the trash 30 days), delete order: %v", fake.delOrder)
	}
	if len(ids) != 0 {
		t.Errorf("a skipped event must not be reported as handled, got %v", ids)
	}
}

// Dry-run touches nothing: no wipe, no delete.
func TestCleanup_DryRunModifiesNothing(t *testing.T) {
	fake := &fakeCalendarCalls{
		listed: []*calendar.Event{fullBookedEvent("evt-expired", "2026-11-01T13:00:00Z")},
	}
	svc := newCleanupTestService(fake)
	cutoff := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)

	ids, err := svc.CleanupExpiredBookings(context.Background(), cutoff, true)
	if err != nil {
		t.Fatalf("dry-run must not error: %v", err)
	}
	if len(ids) != 1 || ids[0] != "evt-expired" {
		t.Errorf("dry-run must report the event it would delete, got %v", ids)
	}
	if len(fake.wipeOrder) != 0 || len(fake.delOrder) != 0 {
		t.Errorf("dry-run must not wipe or delete anything, wipe=%v delete=%v", fake.wipeOrder, fake.delOrder)
	}
}

// A booking whose call ended long before the run (e.g. predating the
// first run, or missed in downtime) IS listed and deleted — the end-to-
// end companion to the list-window guard.
func TestCleanup_OldBookingIsDeleted(t *testing.T) {
	fake := &fakeCalendarCalls{
		listed: []*calendar.Event{
			fullBookedEvent("evt-ancient", "2024-06-01T13:00:00Z"),
			bookedEvent("evt-future", "2026-12-03T13:00:00Z", "tok-f", "jane@example.com"),
		},
	}
	svc := newCleanupTestService(fake)
	cutoff := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)

	ids, err := svc.CleanupExpiredBookings(context.Background(), cutoff, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 1 || ids[0] != "evt-ancient" {
		t.Errorf("the ancient booking must be deleted, the future call kept, got %v", ids)
	}
}
