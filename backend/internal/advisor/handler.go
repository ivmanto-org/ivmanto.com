package advisor

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Handler serves the advisor HTTP endpoints. It is provider-agnostic: the
// Provider interface (provider.go) isolates the LLM backend, and all
// conversation limits and output validation live here.
type Handler struct {
	logger    *slog.Logger
	provider  Provider // nil => 503 (MISTRAL_API_KEY not set)
	catalogue *Catalogue
}

// NewHandler creates an advisor handler. A nil provider means the advisor
// is not configured (no MISTRAL_API_KEY): requests get a clean 503 instead
// of a runtime error, per the deployment contract.
func NewHandler(logger *slog.Logger, provider Provider, catalogue *Catalogue) *Handler {
	return &Handler{logger: logger, provider: provider, catalogue: catalogue}
}

// RegisterRoutes registers the advisor routes with a mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/advisor", h.handleAdvisor)
}

// advisorRequest is what the frontend (spec: AdvisorChat, PR #129 contract
// per Clauco 2026-11-10T21:36Z) sends: the client-generated session_id plus
// the full conversation so far, ending with the new user message.
//
// SessionID is accepted but deliberately NOT used, persisted or logged:
// per the privacy notice (Brio's draft) the session id is a client-side
// random value that must never be linked to an identity — the endpoint
// never sees name/email, and it stays that way. Decoding it into the
// struct documents the contract instead of silently ignoring the field.
type advisorRequest struct {
	SessionID string        `json:"session_id"`
	Messages  []ChatMessage `json:"messages"`
}

// errorBody is the uniform error response. Code is machine-readable so the
// frontend can distinguish turn-limit steering (spec state 6) from errors.
type errorBody struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

func (h *Handler) handleAdvisor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if h.provider == nil {
		// Advisor not configured (MISTRAL_API_KEY unset). Clean 503.
		writeError(w, http.StatusServiceUnavailable, "advisor unavailable", "not_configured")
		return
	}

	var req advisorRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", "bad_request")
		return
	}
	if len(req.Messages) == 0 {
		writeError(w, http.StatusBadRequest, "messages must not be empty", "bad_request")
		return
	}

	userTurns := 0
	for _, msg := range req.Messages {
		switch msg.Role {
		case "user":
			userTurns++
		case "assistant":
			// allowed: prior advisor answers
		default:
			writeError(w, http.StatusBadRequest, `role must be "user" or "assistant"`, "bad_request")
			return
		}
		if len([]rune(msg.Content)) > MaxMessageChars {
			writeError(w, http.StatusBadRequest, "message exceeds 600 characters", "message_too_long")
			return
		}
		if msg.Content == "" {
			writeError(w, http.StatusBadRequest, "message must not be empty", "bad_request")
			return
		}
	}

	// The conversation must end with the new user message.
	last := req.Messages[len(req.Messages)-1]
	if last.Role != "user" {
		writeError(w, http.StatusBadRequest, `conversation must end with a "user" message`, "bad_request")
		return
	}

	// 4-turn limit (spec state 6: steer to a person). The message being
	// sent counts, so this fires on the 5th user message.
	if userTurns > MaxUserTurns {
		writeError(w, http.StatusTooManyRequests, "conversation limit reached — book a call or leave contact", "turn_limit")
		return
	}

	resp, err := h.provider.Reply(r.Context(), req.Messages)
	if err != nil {
		// Spec state 5: never leak internals; the frontend shows its
		// error notice and keeps the visitor's text.
		h.logger.Error("advisor provider failed", "error", err.Error())
		writeError(w, http.StatusBadGateway, "advisor could not answer", "provider_error")
		return
	}

	out := h.validate(resp)
	if err := json.NewEncoder(w).Encode(out); err != nil {
		h.logger.Error("advisor encode response", "error", err.Error())
	}
}

// validate enforces the contract on provider output: catalogue-validated
// ids (unknown ids dropped and logged server-side, per the spec), at most 3
// recommendations, and a normalised next_action.
func (h *Handler) validate(resp Response) Response {
	recs := make([]Recommendation, 0, len(resp.Recommendations))
	seen := make(map[string]bool, len(resp.Recommendations))
	for _, rec := range resp.Recommendations {
		if rec.ID == "" || seen[rec.ID] || !h.catalogue.ValidID(rec.ID) {
			h.logger.Warn("advisor recommendation dropped", "id", rec.ID, "in_catalogue", h.catalogue.ValidID(rec.ID))
			continue
		}
		seen[rec.ID] = true
		recs = append(recs, rec)
		if len(recs) == 3 {
			break
		}
	}

	next := resp.NextAction
	if !ValidNextAction(next) {
		h.logger.Warn("advisor invalid next_action, falling back to call", "next_action", resp.NextAction)
		next = NextActionCall
	}

	return Response{
		Reply:           resp.Reply,
		Recommendations: recs, // always non-nil so JSON emits []
		NextAction:      next,
	}
}

func writeError(w http.ResponseWriter, status int, msg, code string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorBody{Error: msg, Code: code})
}
