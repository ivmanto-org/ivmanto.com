package blog

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	testSAEmail     = "ivmanto-blog-push-sa@ivmanto-com-prod-514215.iam.gserviceaccount.com"
	testAudience    = "https://test.example/api/_internal/pubsub/blog-refresh"
	testStatusToken = "status-shared-secret"
)

// stubOIDC is a test double for the google OIDC validator: it "validates" any
// token of the form valid:<email> and returns the embedded email; everything
// else fails validation.
func stubOIDC(ctx context.Context, rawToken, audience string) (string, error) {
	if audience != testAudience {
		return "", errors.New("audience mismatch")
	}
	if strings.HasPrefix(rawToken, "valid:") {
		return strings.TrimPrefix(rawToken, "valid:"), nil
	}
	return "", errors.New("invalid token")
}

// newTestHandler builds a handler with the stub OIDC validator and no cache
// (auth-rejection paths return before the cache is touched).
func newTestHandler(saEmail, statusToken string) *Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := NewHandler(logger, nil, statusToken, saEmail, testAudience, "")
	h.oidcValidate = stubOIDC
	return h
}

func pubsubPushRequest(t *testing.T, authHeader string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/_internal/pubsub/blog-refresh", strings.NewReader(`{"message":{"attributes":{"objectId":"post.md","eventType":"OBJECT_FINALIZE"}}}`))
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	return req
}

// --- Pub/Sub push endpoint (OIDC bearer auth, #115) ---

func TestPubSubPush_NoSAConfigured_SkipsAuth(t *testing.T) {
	h := newTestHandler("", "")
	// No bearer token at all; with no SA configured (local dev) the request
	// must pass auth. The nil test cache panics on the first post-auth use,
	// which is the sentinel proving auth passed (a 403 would mean it didn't).
	rec := httptest.NewRecorder()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected nil-cache panic proving auth passed, got response %d", rec.Code)
		} else if rec.Code == http.StatusForbidden {
			t.Error("request was rejected at auth despite no SA configured")
		}
	}()
	h.handlePubSubPush(rec, pubsubPushRequest(t, ""))
}

func TestPubSubPush_MissingBearerToken_Rejected(t *testing.T) {
	h := newTestHandler(testSAEmail, "")
	rec := httptest.NewRecorder()
	h.handlePubSubPush(rec, pubsubPushRequest(t, ""))
	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403 for missing bearer token, got %d", rec.Code)
	}
}

func TestPubSubPush_InvalidJWT_Rejected(t *testing.T) {
	h := newTestHandler(testSAEmail, "")
	for _, tok := range []string{"Bearer garbage", "Bearer valid:someone-else@other.iam.gserviceaccount.com"} {
		rec := httptest.NewRecorder()
		h.handlePubSubPush(rec, pubsubPushRequest(t, tok))
		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403 for auth %q, got %d", tok, rec.Code)
		}
	}
}

func TestPubSubPush_ValidOIDCTokenFromPushSA_Accepted(t *testing.T) {
	h := newTestHandler(testSAEmail, "")
	// Use a non-.md object: the handler accepts the push and returns 200
	// without touching the (nil) test cache.
	req := httptest.NewRequest(http.MethodPost, "/api/_internal/pubsub/blog-refresh", strings.NewReader(`{"message":{"attributes":{"objectId":"image.png","eventType":"OBJECT_FINALIZE"}}}`))
	req.Header.Set("Authorization", "Bearer valid:"+testSAEmail)
	rec := httptest.NewRecorder()
	h.handlePubSubPush(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 for valid push-SA token, got %d", rec.Code)
	}
}

func TestPubSubPush_TokenInQueryString_NotAccepted(t *testing.T) {
	// Regression: the token used to ride in the URL; it must no longer be
	// read from there (#115, LB logging of full URLs).
	h := newTestHandler(testSAEmail, "")
	req := httptest.NewRequest(http.MethodPost, "/api/_internal/pubsub/blog-refresh?token="+testSAEmail, strings.NewReader(`{}`))
	rec := httptest.NewRecorder()
	h.handlePubSubPush(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403 when token is only in the query string, got %d", rec.Code)
	}
}

func TestPubSubPush_RejectionDoesNotEchoToken(t *testing.T) {
	h := newTestHandler(testSAEmail, "")
	rec := httptest.NewRecorder()
	h.handlePubSubPush(rec, pubsubPushRequest(t, "Bearer super-secret-value"))
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "super-secret-value") {
		t.Error("response body echoed the token value")
	}
}

// --- Articles status endpoint (shared bearer token, constant-time, #115) ---

func statusRequest(authHeader string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/api/_internal/articles/status", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	return req
}

func TestArticlesStatus_MissingOrWrongToken_Rejected(t *testing.T) {
	h := newTestHandler("", testStatusToken)
	for _, auth := range []string{"", "Bearer wrong", "Bearer " + testStatusToken + "x"} {
		rec := httptest.NewRecorder()
		h.handleArticlesStatus(rec, statusRequest(auth))
		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403 for auth %q, got %d", auth, rec.Code)
		}
	}
}

func TestArticlesStatus_CorrectHeaderToken_PassesAuth(t *testing.T) {
	// nil cache: passing auth must reach the cache call, which panics on a
	// nil receiver - so use a sentinel to prove auth passed.
	h := newTestHandler("", testStatusToken)
	rec := httptest.NewRecorder()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected nil-cache panic proving auth passed, got a normal response")
		}
	}()
	h.handleArticlesStatus(rec, statusRequest("Bearer "+testStatusToken))
}

func TestArticlesStatus_TokenInQueryString_NotAccepted(t *testing.T) {
	// Regression (#115): the status token must no longer be read from the URL.
	h := newTestHandler("", testStatusToken)
	req := httptest.NewRequest(http.MethodGet, "/api/_internal/articles/status?token="+testStatusToken, nil)
	rec := httptest.NewRecorder()
	h.handleArticlesStatus(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403 when token is only in the query string, got %d", rec.Code)
	}
}

func TestArticlesStatus_LengthMismatchNotLeakedByCompare(t *testing.T) {
	// The compare runs on sha-256 digests, so a wrong-length token takes the
	// same path as a same-length one. Both must simply be rejected.
	h := newTestHandler("", testStatusToken)
	for _, tok := range []string{"short", strings.Repeat("x", 4096)} {
		rec := httptest.NewRecorder()
		h.handleArticlesStatus(rec, statusRequest("Bearer "+tok))
		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403 for token of len %d, got %d", len(tok), rec.Code)
		}
	}
}

func TestBearerToken(t *testing.T) {
	cases := []struct {
		header string
		want   string
	}{
		{"", ""},
		{"Basic abc", ""},
		{"Bearer", ""},
		{"Bearer ", ""},
		{"Bearer abc", "abc"},
		{"bearer abc", "abc"}, // scheme is case-insensitive
		{"Bearer  a  ", "a"},  // whitespace trimmed
	}
	for _, c := range cases {
		req := statusRequest("")
		if c.header != "" {
			req.Header.Set("Authorization", c.header)
		}
		if got := bearerToken(req); got != c.want {
			t.Errorf("bearerToken(%q) = %q, want %q", c.header, got, c.want)
		}
	}
}
