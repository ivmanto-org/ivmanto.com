package advisor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSystemPromptContainsCatalogueIDs(t *testing.T) {
	c := testCatalogue()
	p, err := NewMistralProvider(MistralOptions{APIKey: "k", Catalogue: c})
	if err != nil {
		t.Fatalf("NewMistralProvider: %v", err)
	}
	for _, pkg := range c.Packages {
		if !strings.Contains(p.system, pkg.ID) {
			t.Errorf("system prompt missing catalogue id %q", pkg.ID)
		}
	}
}

func TestMistralReplyRoundTrip(t *testing.T) {
	var gotAuth, gotModel string
	var gotBody chatRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"content":"{\"reply\":\"r\",\"recommendations\":[{\"id\":\"gov-quick-check\",\"why\":\"w\"}],\"next_action\":\"buy\"}"}}]}`))
	}))
	defer srv.Close()

	p, err := NewMistralProvider(MistralOptions{
		APIKey:    "test-key",
		Endpoint:  srv.URL,
		Catalogue: testCatalogue(),
	})
	if err != nil {
		t.Fatalf("NewMistralProvider: %v", err)
	}
	resp, err := p.Reply(context.Background(), []ChatMessage{
		{Role: "user", Content: "we need a policy"},
	})
	if err != nil {
		t.Fatalf("Reply: %v", err)
	}
	if gotAuth != "Bearer test-key" {
		t.Errorf("auth header: %q", gotAuth)
	}
	gotModel = gotBody.Model
	if gotModel != DefaultMistralModel {
		t.Errorf("model: %q", gotModel)
	}
	if len(gotBody.Messages) != 2 || gotBody.Messages[0].Role != "system" || gotBody.Messages[1].Content != "we need a policy" {
		t.Errorf("messages: %+v", gotBody.Messages)
	}
	if gotBody.ResponseFormat.Type != "json_object" {
		t.Errorf("response_format: %+v", gotBody.ResponseFormat)
	}
	if resp.Reply != "r" || len(resp.Recommendations) != 1 || resp.Recommendations[0].ID != "gov-quick-check" || resp.NextAction != "buy" {
		t.Errorf("resp: %+v", resp)
	}
}

func TestMistralAPIErrorPropagates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"unauthorized"}`))
	}))
	defer srv.Close()
	p, _ := NewMistralProvider(MistralOptions{APIKey: "k", Endpoint: srv.URL, Catalogue: testCatalogue()})
	if _, err := p.Reply(context.Background(), []ChatMessage{{Role: "user", Content: "x"}}); err == nil {
		t.Fatal("want error on 401")
	}
}

func TestMistralRejectsEmptyKeyAndNilCatalogue(t *testing.T) {
	if _, err := NewMistralProvider(MistralOptions{Catalogue: testCatalogue()}); err == nil {
		t.Fatal("want error on empty key")
	}
	if _, err := NewMistralProvider(MistralOptions{APIKey: "k"}); err == nil {
		t.Fatal("want error on nil catalogue")
	}
}

func TestLoadCatalogueAndValidID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "catalogue.json")
	// One real catalogue v1 package with the new optional fields
	// (line, price_unit, price_note, duration_label) — LoadCatalogue
	// must decode them and ValidID must work against the new ids.
	if err := os.WriteFile(path, []byte(`{
		"version": 1, "currency": "EUR",
		"pillars": [{"id":"governance","name":"Governance","summary":"s"}],
		"packages": [{"id":"gov-quick-check","pillar":"governance","line":"governance","name":"AI governance quick check","outcome":"o","price_eur":1900,"price_note":"credited later","duration_label":"1 day","pillar_recommended":true},
		             {"id":"gov-retainer","pillar":"governance","line":"governance","name":"AI governance retainer","outcome":"o","price_eur":1500,"price_unit":"month"}]
	}`), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := LoadCatalogue(path)
	if err != nil {
		t.Fatalf("LoadCatalogue: %v", err)
	}
	if !c.ValidID("gov-quick-check") {
		t.Error("gov-quick-check must be valid")
	}
	if !c.ValidID("gov-retainer") {
		t.Error("gov-retainer must be valid")
	}
	for _, p := range c.Packages {
		if p.Line != "governance" {
			t.Errorf("package %s: line not decoded: %q", p.ID, p.Line)
		}
	}
	if c.ValidID("nope") {
		t.Error("unknown id must be invalid")
	}

	// No packages -> error (advisor must be grounded).
	if err := os.WriteFile(path, []byte(`{"packages":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadCatalogue(path); err == nil {
		t.Fatal("want error for empty catalogue")
	}
}
