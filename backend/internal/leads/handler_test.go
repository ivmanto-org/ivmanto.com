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
	return `{"name":"Ada Lovelace","email":"ada@example.com","company":"Analytical Engines","message":"We need an AI governance framework.","consent":true}`
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
	resp, _ = post(t, srv, `{"name":"a","email":"b@c.d","company":"x","message":"m"}`)
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
		{"missing fields", `{"email":"a@b.c","consent":true}`, "bad_request", http.StatusBadRequest},
		{"bad email", `{"name":"a","email":"not-an-email","company":"x","message":"m","consent":true}`, "bad_email", http.StatusBadRequest},
		{"empty message", `{"name":"a","email":"a@b.c","company":"x","message":"  ","consent":true}`, "bad_request", http.StatusBadRequest},
		{"invalid json", `not json`, "bad_request", http.StatusBadRequest},
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
	if m.calls != 0 {
		t.Fatalf("mailer must not fire on invalid input, calls=%d", m.calls)
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
