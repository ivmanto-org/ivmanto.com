package leads

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"ivmanto.com/backend/internal/email"
)

// Mailer is the narrow dependence the lead handler needs: send one lead
// notification. SmtpService satisfies it.
type Mailer interface {
	SendLeadNotification(lead email.LeadInfo) error
}

// Handler serves POST /api/leads — the advisor lead-capture endpoint
// (epic #112, spec: LeadForm). A lead is emailed to Nick and is NOT stored
// anywhere: no DB write, no disk, no analytics event, and the message body
// is never logged (personal data per GDPR Art. 4(1); minimal logging).
type Handler struct {
	logger  *slog.Logger
	emailer Mailer
}

// NewHandler creates the lead handler.
func NewHandler(logger *slog.Logger, emailer Mailer) *Handler {
	return &Handler{logger: logger, emailer: emailer}
}

// RegisterRoutes registers the lead routes with a mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/leads", h.handleLead)
}

// leadRequest is the frontend LeadForm payload (contract per PR #129/#126,
// Clauco 2026-11-10 21:36Z, extended for the ivGIL launch 2026-11-11 #112):
// {source:'advisor'|'order'|'contact'|'ivgil', package_id?, role?,
// interest?, name, email, company, message, consent:true,
// insights_optin:bool}. Source tells Nick which surface produced the lead;
// package_id is set only on the advisor order flow; insights_optin is the
// double-opt-in flag for future insights emails (Brio's privacy draft: the
// opt-in is captured now, the actual double opt-in comes later). Role and
// interest are the ivGIL form's optional fields: the visitor's role (free
// string) and whether they want early access or to become a design partner.
// Consent is the required confirmation tick (the processing basis for the
// lead mail is Art. 6(1)(b), steps at the visitor's request before a
// contract — Ivgil's verdict 2026-11-11; only the separate insights opt-in
// is consent, Art. 6(1)(a)); the handler refuses to send without the tick.
// Message is the editable prefill of the conversation summary.
type leadRequest struct {
	Source        string `json:"source"`
	PackageID     string `json:"package_id"`
	Role          string `json:"role"`
	Interest      string `json:"interest"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Company       string `json:"company"`
	Message       string `json:"message"`
	Consent       bool   `json:"consent"`
	InsightsOptIn bool   `json:"insights_optin"`
}

// validLeadSources are the four surfaces that submit to /api/leads
// ('ivgil' added with the ivGIL launch, #112 2026-11-11).
var validLeadSources = map[string]bool{
	"advisor": true,
	"order":   true,
	"contact": true,
	"ivgil":   true,
}

// validLeadInterests are the ivGIL form's interest choices.
var validLeadInterests = map[string]bool{
	"early_access":   true,
	"design_partner": true,
}

const (
	// maxLeadBody guards against oversized payloads (defence in depth).
	maxLeadBody = 16 << 10
	// maxMessageLength caps the conversation summary. Generous on
	// purpose: the summary can carry more than one advisor exchange
	// (the form textarea limits to 600 chars per message, but a crafty
	// client can bypass the textarea).
	maxMessageLength = 4000
	// maxRoleLength caps the ivGIL form's free-string role field.
	maxRoleLength = 200
	// maxShortFieldLength caps name, company and package_id (PR #128 M1):
	// they end up in the mail subject line and must stay sane.
	maxShortFieldLength = 200
)

// hasControlChars reports whether s contains any character that is
// illegal in a mail header value: anything below 0x20 (including CR/LF)
// and DEL. Rejecting them at the handler keeps CR/LF out of the Subject
// line that SendLeadNotification builds from company and name (PR #128
// H1: header injection).
func hasControlChars(s string) bool {
	for _, r := range s {
		if r < 0x20 || r == 0x7f {
			return true
		}
	}
	return false
}

func (h *Handler) handleLead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req leadRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxLeadBody)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", "bad_request")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Company = strings.TrimSpace(req.Company)
	req.Message = strings.TrimSpace(req.Message)
	req.PackageID = strings.TrimSpace(req.PackageID)
	req.Role = strings.TrimSpace(req.Role)
	req.Interest = strings.TrimSpace(req.Interest)

	if !validLeadSources[req.Source] {
		writeError(w, http.StatusBadRequest, `source must be "advisor", "order", "contact" or "ivgil"`, "bad_source")
		return
	}
	// Interest is optional, but when present it must be one of the two
	// ivGIL choices — a garbage value means the contract changed again,
	// and silently dropping it would mislabel the lead for Nick.
	if req.Interest != "" && !validLeadInterests[req.Interest] {
		writeError(w, http.StatusBadRequest, `interest must be "early_access" or "design_partner"`, "bad_interest")
		return
	}
	// Defence in depth for the free-string role (the form limits it, a
	// crafty client can bypass the input).
	if len([]rune(req.Role)) > maxRoleLength {
		writeError(w, http.StatusBadRequest, "role is too long", "bad_request")
		return
	}
	// Header-injection guard (PR #128 H1): name, company, email,
	// package_id and role all end up in mail headers (company and name in
	// the Subject). Control characters there would let a caller add
	// headers to the mail Nick receives.
	for field, value := range map[string]string{
		"name":       req.Name,
		"company":    req.Company,
		"email":      req.Email,
		"package_id": req.PackageID,
		"role":       req.Role,
	} {
		if hasControlChars(value) {
			writeError(w, http.StatusBadRequest, field+" must not contain control characters", "bad_request")
			return
		}
		if len([]rune(value)) > maxShortFieldLength {
			writeError(w, http.StatusBadRequest, field+" is too long", "bad_request")
			return
		}
	}

	if req.Name == "" || req.Email == "" || req.Company == "" {
		writeError(w, http.StatusBadRequest, "name, email and company are required", "bad_request")
		return
	}
	// Parse and use the bare address (PR #128 M2): mail.ParseAddress
	// accepts "Bob <bob@example.com>" and comments, and the raw string
	// would then flow into the mailto: link. Keep only the parsed
	// address, and reject input whose parsed form differs from what was
	// sent.
	addr, err := mail.ParseAddress(req.Email)
	if err != nil {
		writeError(w, http.StatusBadRequest, "a valid email is required", "bad_email")
		return
	}
	if addr.Address != req.Email {
		writeError(w, http.StatusBadRequest, "a valid email is required", "bad_email")
		return
	}
	req.Email = addr.Address
	if req.Message == "" {
		writeError(w, http.StatusBadRequest, "message is required", "bad_request")
		return
	}
	if len([]rune(req.Message)) > maxMessageLength {
		writeError(w, http.StatusBadRequest, "message is too long", "message_too_long")
		return
	}
	// Consent is mandatory: without it the lead must not be processed
	// (no storage, no email). 422 = unprocessable per GDPR intent.
	if !req.Consent {
		writeError(w, http.StatusUnprocessableEntity, "consent is required to submit a lead", "consent_required")
		return
	}

	// Send to Nick. Minimal logging: count and recipient only, never the
	// message/name/email (PII).
	if err := h.emailer.SendLeadNotification(email.LeadInfo{
		Source:        req.Source,
		PackageID:     req.PackageID,
		Role:          req.Role,
		Interest:      req.Interest,
		Name:          req.Name,
		Email:         req.Email,
		Company:       req.Company,
		Message:       req.Message,
		Consent:       req.Consent,
		InsightsOptIn: req.InsightsOptIn,
	}); err != nil {
		h.logger.Error("failed to send lead notification", "error", err.Error())
		writeError(w, http.StatusInternalServerError, "failed to submit lead", "send_failed")
		return
	}
	h.logger.Info("lead notification sent", "ts", time.Now().UTC().Format(time.RFC3339))

	// Spec: success message "Thanks. Nick will reply within one working day."
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Thanks. Nick will reply within one working day.",
	})
}

func writeError(w http.ResponseWriter, status int, msg, code string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg, "code": code})
}
