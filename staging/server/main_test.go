package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func testSite() fstest.MapFS {
	return fstest.MapFS{
		"index.html":          {Data: []byte("home")},
		"packages/index.html": {Data: []byte("packages")},
		"404.html":            {Data: []byte("custom not found")},
		"empty/keep.txt":      {Data: []byte("x")},
	}
}

func request(t *testing.T, method, target, user, password string) *httptest.ResponseRecorder {
	t.Helper()
	users, err := parseUsers("family:fam-secret\npeers:peer-secret")
	if err != nil {
		t.Fatalf("parseUsers: %v", err)
	}
	req := httptest.NewRequest(method, target, nil)
	if user != "" {
		req.SetBasicAuth(user, password)
	}
	rec := httptest.NewRecorder()
	newHandler(users, testSite()).ServeHTTP(rec, req)
	return rec
}

func TestRejectsRequestsWithoutCredentials(t *testing.T) {
	rec := request(t, http.MethodGet, "/", "", "")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
	if rec.Header().Get("WWW-Authenticate") == "" {
		t.Error("missing WWW-Authenticate header")
	}
	if strings.Contains(rec.Body.String(), "home") {
		t.Error("page content leaked without credentials")
	}
}

func TestRejectsWrongPasswordAndUnknownUser(t *testing.T) {
	cases := [][2]string{{"family", "peer-secret"}, {"nobody", "fam-secret"}, {"family", ""}}
	for _, c := range cases {
		if rec := request(t, http.MethodGet, "/", c[0], c[1]); rec.Code != http.StatusUnauthorized {
			t.Errorf("%s/%s: status = %d, want 401", c[0], c[1], rec.Code)
		}
	}
}

func TestServesPagesToEachGroup(t *testing.T) {
	if rec := request(t, http.MethodGet, "/", "family", "fam-secret"); rec.Code != http.StatusOK || rec.Body.String() != "home" {
		t.Errorf("family home: %d %q", rec.Code, rec.Body.String())
	}
	rec := request(t, http.MethodGet, "/packages/", "peers", "peer-secret")
	if rec.Code != http.StatusOK || rec.Body.String() != "packages" {
		t.Errorf("peers packages: %d %q", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("Cache-Control") != "no-cache" {
		t.Error("pages must be served with Cache-Control: no-cache")
	}
}

func TestEveryResponseIsMarkedNoindex(t *testing.T) {
	for _, rec := range []*httptest.ResponseRecorder{
		request(t, http.MethodGet, "/", "", ""),
		request(t, http.MethodGet, "/", "family", "fam-secret"),
		request(t, http.MethodGet, "/missing", "family", "fam-secret"),
	} {
		if !strings.Contains(rec.Header().Get("X-Robots-Tag"), "noindex") {
			t.Errorf("status %d: missing noindex header", rec.Code)
		}
		if rec.Header().Get("X-Content-Type-Options") != "nosniff" || rec.Header().Get("X-Frame-Options") != "DENY" {
			t.Errorf("status %d: missing hardening headers", rec.Code)
		}
	}
}

func TestRobotsTxtBlocksCrawlersWithoutCredentials(t *testing.T) {
	rec := request(t, http.MethodGet, "/robots.txt", "", "")

	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Disallow: /") {
		t.Fatalf("robots.txt: %d %q", rec.Code, rec.Body.String())
	}
}

func TestMissingPathsGetTheSiteNotFoundPage(t *testing.T) {
	for _, target := range []string{"/missing", "/api/advisor", "/empty/", "/../etc/passwd"} {
		rec := request(t, http.MethodGet, target, "family", "fam-secret")
		if rec.Code != http.StatusNotFound || rec.Body.String() != "custom not found" {
			t.Errorf("%s: %d %q", target, rec.Code, rec.Body.String())
		}
	}
}

func TestWritesAreRefused(t *testing.T) {
	if rec := request(t, http.MethodPost, "/api/leads", "family", "fam-secret"); rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}

func TestParseUsersRejectsBadConfiguration(t *testing.T) {
	for _, raw := range []string{"", "  \n", "nopassword", "user:", ":pw"} {
		if _, err := parseUsers(raw); err == nil {
			t.Errorf("parseUsers(%q) should fail", raw)
		}
	}
	users, err := parseUsers("a:1,b:2\nc:with:colon\n")
	if err != nil || len(users) != 3 {
		t.Fatalf("got %d users, err %v", len(users), err)
	}
	if !authorised(users, "c", "with:colon") {
		t.Error("password containing a colon must work")
	}
}
