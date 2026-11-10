package advisor

import "context"

// Conversation and response contract for the home-page AI service advisor
// (epic #112, spec docs/design-system/components/advisor.md).
//
// The HTTP response shape is fixed and consumed verbatim by the frontend
// (PR #126):
//
//	{"reply": "...", "recommendations": [{"id": "...", "why": "..."}], "next_action": "buy|call|schedule|contact"}
//
// recommendation ids are validated against data/catalogue.json — the model
// never invents ids; unknown ids are dropped server-side and logged.

// ChatMessage is one turn of the conversation. Role is "user" or "assistant".
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Recommendation is one catalogue package the advisor recommends. Only Why
// comes from the model; Id must exist in the catalogue and the frontend
// renders name, price and link from catalogue data by id.
type Recommendation struct {
	ID  string `json:"id"`
	Why string `json:"why"`
}

// Response is the advisor's answer for one turn.
type Response struct {
	Reply           string           `json:"reply"`
	Recommendations []Recommendation `json:"recommendations"`
	NextAction      string           `json:"next_action"`
}

// NextAction values (spec: CtaSet). Invalid model output falls back to
// Call — steering to a human is the safe default.
const (
	NextActionBuy      = "buy"
	NextActionCall     = "call"
	NextActionSchedule = "schedule"
	NextActionContact  = "contact"
)

// ValidNextAction reports whether s is one of the four CtaSet actions.
func ValidNextAction(s string) bool {
	switch s {
	case NextActionBuy, NextActionCall, NextActionSchedule, NextActionContact:
		return true
	}
	return false
}

// Provider is an LLM backend that answers one advisor conversation turn.
// Implementations: MistralProvider (mistral.go), fakes in tests. The handler
// is provider-agnostic and applies all validation/limits itself.
type Provider interface {
	// Reply answers the conversation. messages is the full conversation so
	// far, ending with the new user message. Implementations must return the
	// model's answer with the catalogue ids it recommends; the handler — not
	// the provider — is responsible for dropping unknown ids and normalising
	// next_action, so a provider may return them unvalidated.
	Reply(ctx context.Context, messages []ChatMessage) (Response, error)
}

// Conversation limits (epic #112 abuse limits; spec state 6).
const (
	// MaxUserTurns is the maximum number of user messages in one
	// conversation, including the message being sent. The frontend steers
	// to a person at this limit.
	MaxUserTurns = 4

	// MaxMessageChars caps a single user or assistant message.
	MaxMessageChars = 600
)
