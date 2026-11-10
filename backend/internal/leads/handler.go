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

// leadRequest is the frontend LeadForm payload. Consent is a required,
// explicitly-ticked checkbox (Art. 6(1)(a)); the handler refuses to send
// without it. Message is the editable prefill of the conversation summary.
type leadRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Company string `json:"company"`
	Message string `json:"message"`
	Consent bool   `json:"consent"`
}

const (
	// maxLeadBody guards against oversized payloads (defence in depth).
	maxLeadBody = 16 << 10
	// maxMessageLength caps the conversation summary (mirrors the advisor's
	// per-message cap; the form textarea already limits to 600 chars but a
	// crafty client can bypass the textarea).
	maxMessageLength = 4000
)

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

	if req.Name == "" || req.Email == "" || req.Company == "" {
		writeError(w, http.StatusBadRequest, "name, email and company are required", "bad_request")
		return
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		writeError(w, http.StatusBadRequest, "a valid email is required", "bad_email")
		return
	}
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
		Name:    req.Name,
		Email:   req.Email,
		Company: req.Company,
		Message: req.Message,
		Consent: req.Consent,
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
