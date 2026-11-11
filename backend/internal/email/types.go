package email

import "time"

// ContactMessage holds the data from the contact form.
type ContactMessage struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Message        string `json:"message"`
	SendCopyToSelf bool   `json:"sendCopyToSelf"`
}

// BookingConfirmationDetails holds all the necessary information
// to send a booking confirmation email.
type BookingConfirmationDetails struct {
	ToName          string
	ToEmail         string
	StartTime       time.Time
	EndTime         time.Time
	Timezone        string
	MeetLink        string
	CancellationURL string
	IcsUID          string
	IcsSummary      string
	IcsDescription  string
	// IcsTimezone is the visitor's IANA timezone (e.g. "Europe/Athens"). When
	// non-empty, the .ics attachment emits an X-WR-TIMEZONE header at the
	// VCALENDAR level so older Outlook/iOS clients render the event in the
	// visitor's zone. Empty means the header is omitted.
	IcsTimezone string
}

// GeneratedIdea holds the data for a single AI-generated idea.
type GeneratedIdea struct {
	Title   string
	Summary string
}

// LeadInfo holds the data from the lead-capture forms (epic #112, spec:
// LeadForm). Source is which surface produced the lead (advisor | order |
// contact | ivgil); PackageID is set only on the advisor order flow (rendered
// into the email, never used to look anything up server-side). Role and
// Interest are the ivGIL form's optional fields (visitor's role; early
// access vs design partner) — forwarded to Nick, never acted on
// server-side. Name, Email, Company and Message are required; Consent
// documents the Art. 6(1)(a) conversation-sharing consent and must be true
// for the email to be sent. InsightsOptIn is the visitor's insights
// double-opt-in flag — captured and forwarded to Nick, never acted on
// server-side. The lead is emailed to Nick and NOT stored.
type LeadInfo struct {
	Source        string
	PackageID     string
	Role          string
	Interest      string
	Name          string
	Email         string
	Company       string
	Message       string
	Consent       bool
	InsightsOptIn bool
}
