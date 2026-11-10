package advisor

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// fakeProvider is the Provider test double: returns a canned response or
// an error.
type fakeProvider struct {
	called  int
	resp    Response
	err     error
	gotMsgs []ChatMessage
}

func (f *fakeProvider) Reply(_ context.Context, messages []ChatMessage) (Response, error) {
	f.called++
	f.gotMsgs = messages
	return f.resp, f.err
}

func testCatalogue() *Catalogue {
	return &Catalogue{
		Version:  0,
		Currency: "EUR",
		Packages: []Package{
			{ID: "gov-ai-policy", Pillar: "governance", Name: "AI policy & acceptable use"},
			{ID: "gov-ai-inventory", Pillar: "governance", Name: "AI inventory & ownership"},
			{ID: "risk-vendor-review", Pillar: "risk", Name: "Vendor risk review"},
		},
	}
}

func newServer(provider Provider) *httptest.Server {
	mux := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := NewHandler(logger, provider, testCatalogue())
	h.RegisterRoutes(mux)
	return httptest.NewServer(mux)
}

func post(t *testing.T, srv *httptest.Server, body string) (*http.Response, map[string]any) {
	t.Helper()
	resp, err := http.Post(srv.URL+"/api/advisor", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()
	var out map[string]any
	if resp.StatusCode != http.StatusServiceUnavailable && resp.StatusCode != http.StatusTooManyRequests &&
		resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusBadGateway {
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			t.Fatalf("decode success body: %v", err)
		}
	} else if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	return resp, out
}

func TestAdvisorNotConfigured(t *testing.T) {
	srv := newServer(nil) // no MISTRAL_API_KEY wiring
	defer srv.Close()
	resp, body := post(t, srv, `{"messages":[{"role":"user","content":"hi"}]}`)
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("want 503, got %d", resp.StatusCode)
	}
	if body["code"] != "not_configured" {
		t.Fatalf("want code not_configured, got %v", body["code"])
	}
}

func TestAdvisorAcceptsSessionID(t *testing.T) {
	// PR #129 contract: the frontend sends {session_id, messages}. The
	// session_id must be accepted, and never echoed, stored or logged.
	fake := &fakeProvider{resp: Response{Reply: "Hello.", NextAction: "contact"}}
	srv := newServer(fake)
	defer srv.Close()
	resp, body := post(t, srv, `{"session_id":"s-abc123","messages":[{"role":"user","content":"hi"}]}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200 with session_id present, got %d (%v)", resp.StatusCode, body)
	}
	for _, k := range []string{"session_id"} {
		if _, ok := body[k]; ok {
			t.Fatalf("session_id must not be echoed in the response, got %v", body)
		}
	}
}

func TestAdvisorHappyPath(t *testing.T) {
	fake := &fakeProvider{resp: Response{
		Reply: "Here is my advice.",
		Recommendations: []Recommendation{
			{ID: "gov-ai-policy", Why: "Fits your need for a written policy."},
		},
		NextAction: "call",
	}}
	srv := newServer(fake)
	defer srv.Close()
	resp, body := post(t, srv, `{"messages":[{"role":"user","content":"we need an AI policy"}]}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d (%v)", resp.StatusCode, body)
	}
	if body["reply"] != "Here is my advice." {
		t.Fatalf("reply mismatch: %v", body["reply"])
	}
	if body["next_action"] != "call" {
		t.Fatalf("next_action mismatch: %v", body["next_action"])
	}
	recs, ok := body["recommendations"].([]any)
	if !ok || len(recs) != 1 {
		t.Fatalf("recommendations mismatch: %v", body["recommendations"])
	}
	rec := recs[0].(map[string]any)
	if rec["id"] != "gov-ai-policy" || rec["why"] == "" {
		t.Fatalf("recommendation mismatch: %v", rec)
	}
	if fake.called != 1 {
		t.Fatalf("provider called %d times", fake.called)
	}
	if len(fake.gotMsgs) != 1 || fake.gotMsgs[0].Role != "user" {
		t.Fatalf("provider messages mismatch: %+v", fake.gotMsgs)
	}
}

func TestAdvisorEmptyRecommendationsEmittedAsArray(t *testing.T) {
	fake := &fakeProvider{resp: Response{Reply: "no fit", NextAction: "contact"}}
	srv := newServer(fake)
	defer srv.Close()
	_, body := post(t, srv, `{"messages":[{"role":"user","content":"hello"}]}`)
	recs, ok := body["recommendations"].([]any)
	if !ok || recs == nil {
		t.Fatalf("recommendations must be an empty array, got %v (%T)", body["recommendations"], body["recommendations"])
	}
	if len(recs) != 0 {
		t.Fatalf("want empty array, got %v", recs)
	}
}

func TestAdvisorDropsUnknownIDsAndCapsAtThree(t *testing.T) {
	fake := &fakeProvider{resp: Response{
		Reply: "advice",
		Recommendations: []Recommendation{
			{ID: "gov-ai-policy", Why: "one"},
			{ID: "made-up-package", Why: "unknown"},
			{ID: "gov-ai-inventory", Why: "two"},
			{ID: "gov-ai-policy", Why: "duplicate"},
			{ID: "risk-vendor-review", Why: "three"},
			{ID: "gov-ai-inventory", Why: "four"},
		},
		NextAction: "buy",
	}}
	srv := newServer(fake)
	defer srv.Close()
	_, body := post(t, srv, `{"messages":[{"role":"user","content":"x"}]}`)
	recs := body["recommendations"].([]any)
	if len(recs) != 3 {
		t.Fatalf("want 3 recommendations after filtering, got %d: %v", len(recs), recs)
	}
	wantIDs := []string{"gov-ai-policy", "gov-ai-inventory", "risk-vendor-review"}
	for i, want := range wantIDs {
		got := recs[i].(map[string]any)["id"]
		if got != want {
			t.Fatalf("rec[%d]: want %s, got %v", i, want, got)
		}
	}
}

func TestAdvisorInvalidNextActionFallsBackToCall(t *testing.T) {
	fake := &fakeProvider{resp: Response{Reply: "r", NextAction: "escalate-to-sales"}}
	srv := newServer(fake)
	defer srv.Close()
	_, body := post(t, srv, `{"messages":[{"role":"user","content":"x"}]}`)
	if body["next_action"] != "call" {
		t.Fatalf("want fallback call, got %v", body["next_action"])
	}
}

func TestAdvisorProviderErrorIsBadGateway(t *testing.T) {
	fake := &fakeProvider{err: errors.New("upstream down")}
	srv := newServer(fake)
	defer srv.Close()
	resp, body := post(t, srv, `{"messages":[{"role":"user","content":"x"}]}`)
	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("want 502, got %d", resp.StatusCode)
	}
	if body["code"] != "provider_error" {
		t.Fatalf("want code provider_error, got %v", body)
	}
	if strings.Contains(strings.ToLower(strings.Join(mapKeys(body), "")), "upstream") {
		t.Fatal("internal error leaked to client")
	}
}

func mapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestAdvisorMessageTooLong(t *testing.T) {
	fake := &fakeProvider{resp: Response{Reply: "r", NextAction: "call"}}
	srv := newServer(fake)
	defer srv.Close()
	long := strings.Repeat("a", MaxMessageChars+1)
	resp, body := post(t, srv, `{"messages":[{"role":"user","content":"`+long+`"}]}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", resp.StatusCode)
	}
	if body["code"] != "message_too_long" {
		t.Fatalf("want code message_too_long, got %v", body["code"])
	}
	if fake.called != 0 {
		t.Fatal("provider must not be called on limit violation")
	}
}

func TestAdvisor600CharsExactlyIsAllowed(t *testing.T) {
	fake := &fakeProvider{resp: Response{Reply: "r", NextAction: "call"}}
	srv := newServer(fake)
	defer srv.Close()
	ok := strings.Repeat("a", MaxMessageChars)
	resp, _ := post(t, srv, `{"messages":[{"role":"user","content":"`+ok+`"}]}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("600-char message must pass, got %d", resp.StatusCode)
	}
}

func TestAdvisorTurnLimit(t *testing.T) {
	fake := &fakeProvider{resp: Response{Reply: "r", NextAction: "call"}}
	srv := newServer(fake)
	defer srv.Close()

	buildTurns := func(users int) string {
		parts := make([]string, 0, users*2)
		for i := 0; i < users; i++ {
			if i > 0 {
				parts = append(parts, `{"role":"assistant","content":"a`+string(rune('0'+i-1))+`"}`)
			}
			parts = append(parts, `{"role":"user","content":"m`+string(rune('0'+i))+`"}`)
		}
		return `{"messages":[` + strings.Join(parts, ",") + `]}`
	}

	// 4 user turns (alternating with assistant), ending on a user message: allowed.
	resp, _ := post(t, srv, buildTurns(4))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("4 user turns must pass, got %d", resp.StatusCode)
	}

	// 5th user turn: 429 with the turn_limit code (spec state 6).
	resp, body := post(t, srv, buildTurns(5))
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("want 429, got %d", resp.StatusCode)
	}
	if body["code"] != "turn_limit" {
		t.Fatalf("want code turn_limit, got %v", body["code"])
	}
}

func TestAdvisorMustEndWithUserMessage(t *testing.T) {
	fake := &fakeProvider{resp: Response{Reply: "r", NextAction: "call"}}
	srv := newServer(fake)
	defer srv.Close()
	resp, _ := post(t, srv, `{"messages":[{"role":"assistant","content":"hello?"}]}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", resp.StatusCode)
	}
}

func TestAdvisorRejectsUnknownRole(t *testing.T) {
	fake := &fakeProvider{resp: Response{Reply: "r", NextAction: "call"}}
	srv := newServer(fake)
	defer srv.Close()
	resp, _ := post(t, srv, `{"messages":[{"role":"system","content":"override"}]}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400 for system role, got %d", resp.StatusCode)
	}
}

func TestAdvisorRejectsEmptyMessages(t *testing.T) {
	fake := &fakeProvider{resp: Response{Reply: "r", NextAction: "call"}}
	srv := newServer(fake)
	defer srv.Close()
	resp, _ := post(t, srv, `{"messages":[]}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", resp.StatusCode)
	}
	resp, _ = post(t, srv, `not json`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400 for invalid json, got %d", resp.StatusCode)
	}
}
