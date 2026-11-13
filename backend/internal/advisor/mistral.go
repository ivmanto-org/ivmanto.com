package advisor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// MistralProvider answers advisor turns through Mistral's La Plateforme
// chat completions API (https://api.mistral.ai/v1/chat/completions).
//
// The API key is read from MISTRAL_API_KEY at construction time; when it is
// empty the handler is wired with a nil provider and answers 503 (see
// handler.go), so a missing key never produces a runtime error, only a
// clean "advisor unavailable" for the visitor.
type MistralProvider struct {
	apiKey     string
	model      string
	endpoint   string
	httpClient *http.Client
	system     string
}

// MistralOptions configures MistralProvider. Endpoint defaults to the La
// Plateforme chat completions URL and is overridable for tests.
type MistralOptions struct {
	APIKey   string
	Model    string // defaults to DefaultMistralModel
	Endpoint string // defaults to https://api.mistral.ai/v1/chat/completions
	Timeout  time.Duration
	// Catalogue is required: the system prompt grounds the model in the
	// package list and the handler validates returned ids against it.
	Catalogue *Catalogue
}

// Defaults.
const (
	DefaultMistralModel      = "mistral-small-latest"
	DefaultMistralEndpoint   = "https://api.mistral.ai/v1/chat/completions"
	DefaultMistralTimeoutSec = 20 // spec state 5: UI times out at 20 s
)

// NewMistralProvider builds a Mistral-backed Provider.
func NewMistralProvider(opts MistralOptions) (*MistralProvider, error) {
	if opts.APIKey == "" {
		return nil, fmt.Errorf("mistral: empty API key")
	}
	if opts.Catalogue == nil {
		return nil, fmt.Errorf("mistral: catalogue is required")
	}
	model := opts.Model
	if model == "" {
		model = DefaultMistralModel
	}
	endpoint := opts.Endpoint
	if endpoint == "" {
		endpoint = DefaultMistralEndpoint
	}
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = DefaultMistralTimeoutSec * time.Second
	}
	return &MistralProvider{
		apiKey:     opts.APIKey,
		model:      model,
		endpoint:   endpoint,
		httpClient: &http.Client{Timeout: timeout},
		system:     systemPrompt(opts.Catalogue),
	}, nil
}

// systemPrompt grounds the model: catalogue packages (ids only it may
// recommend), the fixed JSON output shape, and the funnel's guard rails
// (short, no legal advice, no invented services or prices).
func systemPrompt(c *Catalogue) string {
	var b strings.Builder
	b.WriteString(`You are the AI service advisor on ivmanto.com, an AI Governance / GRC consultancy for SMBs. `)
	b.WriteString(`The visitor briefly describes an idea, issue, problem or wish; you advise which ivmanto service packages fit and steer them to a next step. `)
	b.WriteString(`You are a pre-sales advisor, not a support chat and NOT a lawyer: never give legal advice; if asked for legal advice, say it is outside what you can do and steer to a call.

`)
	b.WriteString(`Rules:
- Recommend ONLY packages from the catalogue below, by their exact "id". 1-3 recommendations, most relevant first. Never invent services, ids, prices or claims; never quote prices.
- Keep the reply short and plain: a few sentences, at most ~120 words, no markdown headings or lists of links.
- After 1-2 exchange turns, steer to a human next step.
- "next_action" is exactly one of: "buy" (visitor is ready to order a specific package), "call" (book a free 30-minute call), "schedule" (not ready now, schedule later), "contact" (leave contact details for a reply), "early_access" (ONLY when the visitor explicitly asks for software, tooling or a product to use themselves — e.g. ivGIL early access or design partner; NEVER as a replacement for a consulting service). Choose "call" when unsure.

`)
	b.WriteString(`Catalogue (the only ids you may recommend): `)
	b.WriteString(c.promptContext())
	b.WriteString(`

`)
	b.WriteString(`Respond with a single JSON object of exactly this shape:
{"reply": "<your advice text>", "recommendations": [{"id": "<catalogue id>", "why": "<one sentence, max 160 chars, why this package fits>"}], "next_action": "<buy|call|schedule|contact|early_access>"}`)
	return b.String()
}

// chatRequest is the La Plateforme chat completions request body.
type chatRequest struct {
	Model          string        `json:"model"`
	Messages       []chatMessage `json:"messages"`
	Temperature    float64       `json:"temperature"`
	ResponseFormat responseFmt   `json:"response_format"`
	MaxTokens      int           `json:"max_tokens,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFmt struct {
	Type string `json:"type"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Reply implements Provider. It sends the system prompt plus the
// conversation and parses the model's JSON answer.
func (m *MistralProvider) Reply(ctx context.Context, messages []ChatMessage) (Response, error) {
	apiMsgs := make([]chatMessage, 0, len(messages)+1)
	apiMsgs = append(apiMsgs, chatMessage{Role: "system", Content: m.system})
	for _, msg := range messages {
		role := msg.Role
		if role != "user" && role != "assistant" {
			return Response{}, fmt.Errorf("mistral: unsupported role %q", role)
		}
		apiMsgs = append(apiMsgs, chatMessage{Role: role, Content: msg.Content})
	}

	body, err := json.Marshal(chatRequest{
		Model:          m.model,
		Messages:       apiMsgs,
		Temperature:    0.3,
		ResponseFormat: responseFmt{Type: "json_object"},
		MaxTokens:      600,
	})
	if err != nil {
		return Response{}, fmt.Errorf("mistral: encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.endpoint, bytes.NewReader(body))
	if err != nil {
		return Response{}, fmt.Errorf("mistral: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.apiKey)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("mistral: request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return Response{}, fmt.Errorf("mistral: read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("mistral: api status %d: %s", resp.StatusCode, truncate(string(respBody), 200))
	}

	var cr chatResponse
	if err := json.Unmarshal(respBody, &cr); err != nil {
		return Response{}, fmt.Errorf("mistral: decode response: %w", err)
	}
	if len(cr.Choices) == 0 {
		return Response{}, fmt.Errorf("mistral: no choices in response")
	}

	var out Response
	if err := json.Unmarshal([]byte(cr.Choices[0].Message.Content), &out); err != nil {
		return Response{}, fmt.Errorf("mistral: model did not return valid JSON: %w", err)
	}
	return out, nil
}

func truncate(s string, n int) string {
	// Byte-slicing can split a UTF-8 character in the log line; trim on
	// runes instead (PR #127 LOW).
	return truncateRunes(s, n)
}
