package email

import (
	"strings"
	"testing"
	"time"
)

// TestBookingConfirmationHTML_AthensInJune covers the bug-report
// scenario end-to-end at the render layer: a 15:30 CEST event in
// June (i.e. an instant where Athens is on EEST) is shown to an
// Athens visitor as 4:30 PM - 5:00 PM (EEST), not the organiser's
// 3:30 PM - 4:00 PM CEST.
//
// This is the test the original PR plan called for but didn't
// ship; without it, the DST bug caught in PR review (the helper
// probing at a fixed Jan 1 instant) would not have been caught.
func TestBookingConfirmationHTML_AthensInJune(t *testing.T) {
	athens, err := time.LoadLocation("Europe/Athens")
	if err != nil {
		t.Skipf("tzdata not available for Europe/Athens: %v", err)
	}

	// 15:30 CEST = 13:30 UTC, same absolute instant as 16:30 EEST.
	// The handler passes startTime.In(visitorLoc) so the formatter
	// prints in Athens' wall-clock. We construct the value the
	// way the handler does: as an absolute UTC instant.
	startUTC := time.Date(2026, 6, 15, 13, 30, 0, 0, time.UTC)
	endUTC := time.Date(2026, 6, 15, 14, 0, 0, 0, time.UTC)
	start := startUTC.In(athens)
	end := endUTC.In(athens)
	// EEST is what the call site would compute via start.In(loc).Format("MST").
	label := start.Format("MST")

	body := buildBookingConfirmationHTML(BookingConfirmationDetails{
		ToName:     "Marina Muchakova",
		ToEmail:    "marina@example.com",
		StartTime:  start,
		EndTime:    end,
		Timezone:   label,
		MeetLink:   "https://meet.google.com/abc-defg-hij",
		IcsUID:     "test-uid@ivmanto.com",
		IcsSummary: "Consultation: Marina Muchakova",
	})

	// Must contain the visitor's wall-clock and abbreviation, not
	// the organiser's.
	wantSubs := []string{
		"4:30 PM",
		"5:00 PM",
		"(EEST)",
		"Marina Muchakova",
		"https://meet.google.com/abc-defg-hij",
	}
	for _, s := range wantSubs {
		if !strings.Contains(body, s) {
			t.Errorf("expected body to contain %q, body was:\n%s", s, body)
		}
	}

	// Must NOT contain the organiser's wall-clock or the old empty
	// parenthesised bug.
	notWantSubs := []string{
		"3:30 PM", // organiser's CEST start
		"4:00 PM", // organiser's CEST end
		"()",      // the original empty-parentheses bug
		"(CEST)",  // we expect the visitor's TZ, not the organiser's
	}
	for _, s := range notWantSubs {
		if strings.Contains(body, s) {
			t.Errorf("expected body NOT to contain %q, body was:\n%s", s, body)
		}
	}
}

// TestBookingConfirmationHTML_AthensInWinter covers the other half of
// the DST split: a 15:30 CET event in January is shown to an Athens
// visitor as 4:30 PM - 5:00 PM (EET). Pairs with the June test to
// prove the abbreviation follows the actual event date, not a fixed
// probe instant.
func TestBookingConfirmationHTML_AthensInWinter(t *testing.T) {
	athens, err := time.LoadLocation("Europe/Athens")
	if err != nil {
		t.Skipf("tzdata not available for Europe/Athens: %v", err)
	}

	// 15:30 CET = 14:30 UTC = 16:30 EET in Athens.
	startUTC := time.Date(2026, 1, 15, 14, 30, 0, 0, time.UTC)
	endUTC := time.Date(2026, 1, 15, 15, 0, 0, 0, time.UTC)
	start := startUTC.In(athens)
	end := endUTC.In(athens)
	label := start.Format("MST")

	body := buildBookingConfirmationHTML(BookingConfirmationDetails{
		ToName:    "Marina Muchakova",
		ToEmail:   "marina@example.com",
		StartTime: start,
		EndTime:   end,
		Timezone:  label,
	})

	if !strings.Contains(body, "4:30 PM") {
		t.Errorf("expected body to contain Athens winter start '4:30 PM', body was:\n%s", body)
	}
	if !strings.Contains(body, "5:00 PM") {
		t.Errorf("expected body to contain Athens winter end '5:00 PM', body was:\n%s", body)
	}
	if !strings.Contains(body, "(EET)") {
		t.Errorf("expected body to contain '(EET)' for Athens in January, body was:\n%s", body)
	}
	if strings.Contains(body, "(EEST)") {
		t.Errorf("expected body NOT to contain '(EEST)' in winter, body was:\n%s", body)
	}
}

// TestBookingConfirmationHTML_NoMeetLinkOmitsLine guards the layout:
// when there's no Meet link, the bullet point shouldn't render with
// empty content.
func TestBookingConfirmationHTML_NoMeetLinkOmitsLine(t *testing.T) {
	body := buildBookingConfirmationHTML(BookingConfirmationDetails{
		ToName:    "Test",
		ToEmail:   "test@example.com",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(30 * time.Minute),
		Timezone:  "UTC",
	})
	if strings.Contains(body, "Google Meet Link") {
		t.Errorf("expected no Meet-link line when MeetLink is empty, body was:\n%s", body)
	}
}

func TestBuildLeadNotificationHTMLEscapesUserInput(t *testing.T) {
	lead := LeadInfo{
		Name:    "<b>Ada</b>",
		Email:   "a@b.c",
		Company: `<script>alert("x")</script>`,
		Message: `Contact: a</p><script>evil()</script>`,
		Consent: true,
	}
	html := buildLeadNotificationHTML(lead)
	for _, bad := range []string{"<script>", "</script>", "<b>Ada</b>"} {
		if strings.Contains(html, bad) {
			t.Errorf("user input leaked into email HTML unescaped: %q", bad)
		}
	}
	for _, good := range []string{"&lt;script&gt;", "&lt;b&gt;Ada&lt;/b&gt;", "a&lt;/p&gt;"} {
		if !strings.Contains(html, good) {
			t.Errorf("expected escaped form %q in html:\n%s", good, html)
		}
	}
}

func TestBuildMessageReplyToPresent(t *testing.T) {
	msg := string(buildMessage("IVMANTO", "noreply@ivmanto.com", []string{"nick@ivmanto.com"}, nil, "visitor@example.com", "New advisor lead", "b123"))
	if !strings.Contains(msg, "Reply-To: visitor@example.com\r\n") {
		t.Errorf("expected Reply-To header with visitor address, message was:\n%s", msg)
	}
}

func TestBuildMessageReplyToOmittedWhenEmpty(t *testing.T) {
	msg := string(buildMessage("IVMANTO", "noreply@ivmanto.com", []string{"nick@ivmanto.com"}, nil, "", "Subject", "b123"))
	if strings.Contains(msg, "Reply-To:") {
		t.Errorf("expected no Reply-To header when replyTo is empty, message was:\n%s", msg)
	}
}

func TestBuildMessageReplyToAndSubjectStripCRLFInjection(t *testing.T) {
	msg := string(buildMessage(
		"IVMANTO", "noreply@ivmanto.com",
		[]string{"nick@ivmanto.com"}, nil,
		"visitor@example.com\r\nBcc: attacker@evil.example",
		"Subject\r\nX-Spam: yes",
		"b123"))
	for _, bad := range []string{"\r\nBcc:", "\r\nX-Spam:"} {
		if strings.Contains(msg, bad) {
			t.Errorf("CR/LF injection started a new header line %q in message:\n%s", bad, msg)
		}
	}
	if !strings.Contains(msg, "Reply-To: visitor@example.comBcc: attacker@evil.example\r\n") {
		t.Errorf("expected Reply-To with CR/LF stripped, message was:\n%s", msg)
	}
	if !strings.Contains(msg, "Subject: SubjectX-Spam: yes\r\n") {
		t.Errorf("expected Subject with CR/LF stripped, message was:\n%s", msg)
	}
}
