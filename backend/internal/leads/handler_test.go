package leads

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ivmanto.com/backend/internal/email"
)

type fakeMailer struct {
	sent  []email.LeadInfo
	fail  bool
	calls int
}

func (f *fakeMailer) SendLeadNotification(lead email.LeadInfo) error {
	f.calls++
	if f.fail {
		return io.EOF
	}
	f.sent = append(f.sent, lead)
	return nil
}

func newServer(m *fakeMailer) *httptest.Server {
	mux := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	NewHandler(logger, m).RegisterRoutes(mux)
	return httptest.NewServer(mux)
}

func post(t *testing.T, srv *httptest.Server, body string) (*http.Response, map[string]any) {
	t.Helper()
	resp, err := http.Post(srv.URL+"/api/leads", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()
	var out map[string]any
	json.NewDecoder(resp.Body).Decode(&out)
	return resp, out
}

func validBody() string {
	return `{"source":"advisor","name":"Ada Lovelace","email":"ada@example.com","company":"Analytical Engines","message":"We need an AI governance framework.","consent":true,"insights_optin":true}`
}

func TestLeadIvgilSourceContract(t *testing.T) {
	// ivGIL launch contract (#112, Nick 2026-11-11): source 'ivgil' plus
	// optional role (free string) and interest (early_access |
	// design_partner). Standalone per the suite convention: each
	// contract point asserted independently.
	m := &fakeMailer{}
	srv := newServer(m)
	defer srv.Close()

	// ivgil with role + interest is accepted and forwarded to Nick.
	resp, body := post(t, srv, `{"source":"ivgil","role":"Head of Data","interest":"early_access","name":"Ada","email":"ada@example.com","company":"Analytical Engines","message":"We want early access to ivGIL.","consent":true}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200 for ivgil lead, got %d (%v)", resp.StatusCode, body)
	}
	if m.calls != 1 || len(m.sent) != 1 {
		t.Fatalf("mailer must be called once, got %d", m.calls)
	}
	got := m.sent[0]
	if got.Source != "ivgil" || got.Role != "Head of Data" || got.Interest != "early_access" {
		t.Fatalf("ivgil lead not forwarded verbatim: %+v", got)
	}

	// ivgil without role/interest is fine — both are optional.
	resp, _ = post(t, srv, `{"source":"ivgil","name":"Ada","email":"ada@example.com","company":"Analytical Engines","message":"Design partner, please.","consent":true}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200 for ivgil lead without optional fields, got %d", resp.StatusCode)
	}

	// Invalid interest is rejected, not silently dropped.
	resp, body = post(t, srv, `{"source":"ivgil","interest":"newsletter","name":"Ada","email":"ada@example.com","company":"Analytical Engines","message":"x","consent":true}`)
	if resp.StatusCode != http.StatusBadRequest || body["code"] != "bad_interest" {
		t.Fatalf("invalid interest must 400 bad_interest, got %d (%v)", resp.StatusCode, body)
	}
	if m.calls != 2 {
		t.Fatalf("mailer must not fire for the rejected lead, got %d calls", m.calls)
	}
}

func TestLeadHappyPath(t *testing.T) {
	m := &fakeMailer{}
	srv := newServer(m)
	defer srv.Close()
	resp, body := post(t, srv, validBody())
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d (%v)", resp.StatusCode, body)
	}
	if m.calls != 1 {
		t.Fatalf("mailer calls: %d", m.calls)
	}
	if len(m.sent) != 1 {
		t.Fatalf("sent: %d", len(m.sent))
	}
	lead := m.sent[0]
	if lead.Name != "Ada Lovelace" || lead.Email != "ada@example.com" || lead.Company != "Analytical Engines" || !lead.Consent {
		t.Fatalf("lead mismatch: %+v", lead)
	}
	if lead.Source != "advisor" || !lead.InsightsOptIn {
		t.Fatalf("source/insights passthrough mismatch: %+v", lead)
	}
	if strings.Contains(body["message"].(string), "within one working day") == false {
		t.Fatalf("success message mismatch: %v", body["message"])
	}
}

func TestLeadRequiresConsent(t *testing.T) {
	m := &fakeMailer{}
	srv := newServer(m)
	defer srv.Close()
	body := validBody()
	// consent explicitly false
	resp, out := post(t, srv, body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("consent:true must pass, got %d", resp.StatusCode)
	}
	_ = out

	resp, out = post(t, srv, strings.Replace(body, `"consent":true`, `"consent":false`, 1))
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("want 422 without consent, got %d (%v)", resp.StatusCode, out)
	}
	if out["code"] != "consent_required" {
		t.Fatalf("want code consent_required, got %v", out["code"])
	}
	if m.calls != 1 {
		t.Fatalf("mailer must not fire without consent, calls=%d", m.calls)
	}

	// consent omitted entirely: default false -> refused
	resp, _ = post(t, srv, `{"source":"contact","name":"a","email":"b@c.d","company":"x","message":"m"}`)
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("want 422 when consent omitted, got %d", resp.StatusCode)
	}
}

func TestLeadValidation(t *testing.T) {
	m := &fakeMailer{}
	srv := newServer(m)
	defer srv.Close()

	cases := []struct {
		name string
		body string
		code string
		stat int
	}{
		{"missing fields", `{"email":"a@b.c","source":"contact","consent":true}`, "bad_request", http.StatusBadRequest},
		{"bad email", `{"source":"contact","name":"a","email":"not-an-email","company":"x","message":"m","consent":true}`, "bad_email", http.StatusBadRequest},
		{"empty message", `{"source":"order","name":"a","email":"a@b.c","company":"x","message":"  ","consent":true}`, "bad_request", http.StatusBadRequest},
		{"invalid json", `not json`, "bad_request", http.StatusBadRequest},
		{"missing source", `{"name":"a","email":"a@b.c","company":"x","message":"m","consent":true}`, "bad_source", http.StatusBadRequest},
		{"invalid source", `{"source":"newsletter","name":"a","email":"a@b.c","company":"x","message":"m","consent":true}`, "bad_source", http.StatusBadRequest},
		{"package id passthrough ok", `{"source":"order","package_id":"gov-ai-policy","name":"a","email":"a@b.c","company":"x","message":"m","consent":true,"insights_optin":false}`, "", http.StatusOK},
	}
	for _, tc := range cases {
		resp, out := post(t, srv, tc.body)
		if resp.StatusCode != tc.stat {
			t.Errorf("%s: want %d, got %d", tc.name, tc.stat, resp.StatusCode)
		}
		if tc.code != "" && out["code"] != tc.code {
			t.Errorf("%s: want code %s, got %v", tc.name, tc.code, out["code"])
		}
	}
	// The passthrough-ok case legitimately sends one lead.
	if m.calls != 1 {
		t.Fatalf("mailer calls: want 1 (passthrough-ok case), got %d", m.calls)
	}
}

func TestLeadSendFailureIs500NotPanic(t *testing.T) {
	m := &fakeMailer{fail: true}
	srv := newServer(m)
	defer srv.Close()
	resp, body := post(t, srv, validBody())
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("want 500, got %d (%v)", resp.StatusCode, body)
	}
	if body["code"] != "send_failed" {
		t.Fatalf("want code send_failed, got %v", body)
	}
}
