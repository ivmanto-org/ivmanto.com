package blog

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/idtoken"
)

// DefaultPushAudience is the OIDC audience the Pub/Sub push subscription's
// oidc_token is configured with (infra/terraform/envs/prod/blog_gcs.tf,
// local.blog_push_endpoint). Overridable via PUBSUB_PUSH_AUDIENCE for
// non-prod environments.
const DefaultPushAudience = "https://ivmanto.com/api/_internal/pubsub/blog-refresh"

// oidcTokenValidator validates a Google-signed OIDC JWT (the Bearer token a
// Pub/Sub push subscription sends) against an audience. Implementations must
// verify signature, issuer, audience and expiry, and return the token's
// email claim. An interface so tests can stub the network fetch.
type oidcTokenValidator func(ctx context.Context, rawToken, audience string) (email string, err error)

// googleOIDCValidator is the production validator: it verifies the JWT
// signature against Google's public certs plus iss/aud/exp, using
// google.golang.org/api/idtoken.
func googleOIDCValidator(ctx context.Context, rawToken, audience string) (string, error) {
	payload, err := idtoken.Validate(ctx, rawToken, audience)
	if err != nil {
		return "", err
	}
	email, _ := payload.Claims["email"].(string)
	return email, nil
}

// bearerToken extracts the token from an "Authorization: Bearer <token>"
// header. Returns "" when the header is absent or malformed. Never logs the
// token value.
func bearerToken(r *http.Request) string {
	const prefix = "Bearer "
	auth := r.Header.Get("Authorization")
	if len(auth) > len(prefix) && strings.EqualFold(auth[:len(prefix)], prefix) {
		return strings.TrimSpace(auth[len(prefix):])
	}
	return ""
}

// Handler serves blog article endpoints.
type Handler struct {
	logger       *slog.Logger
	cache        *Cache
	statusToken  string // optional shared secret for the operator-only status endpoint (Bearer header)
	pushSAEmail  string // optional expected OIDC signer email for Pub/Sub pushes (empty = skip in local dev)
	pushAudience string // expected OIDC audience for Pub/Sub pushes
	oidcValidate oidcTokenValidator
	webhookURL   string // optional Cloud Build webhook URL for frontend rebuilds
}

// NewHandler creates a new blog handler.
func NewHandler(logger *slog.Logger, cache *Cache, statusToken, pushSAEmail, pushAudience string, webhookURL string) *Handler {
	if pushAudience == "" {
		pushAudience = DefaultPushAudience
	}
	return &Handler{
		logger:       logger,
		cache:        cache,
		statusToken:  statusToken,
		pushSAEmail:  pushSAEmail,
		pushAudience: pushAudience,
		oidcValidate: googleOIDCValidator,
		webhookURL:   webhookURL,
	}
}

// RegisterRoutes sets up the routing for blog endpoints.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/articles", h.handleListArticles)
	mux.HandleFunc("GET /api/articles/{slug}", h.handleGetArticle)
	mux.HandleFunc("POST /api/_internal/pubsub/blog-refresh", h.handlePubSubPush)
	mux.HandleFunc("GET /api/_internal/articles/status", h.handleArticlesStatus)
	mux.HandleFunc("GET /api/sitemap-blog.xml", h.handleBlogSitemap)
}

// handleListArticles returns metadata for all published articles.
func (h *Handler) handleListArticles(w http.ResponseWriter, r *http.Request) {
	articles := h.cache.GetAllPublished()
	h.respondJSON(w, http.StatusOK, articles)
}

// handleGetArticle returns a single article's metadata and HTML content.
func (h *Handler) handleGetArticle(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		h.respondError(w, http.StatusBadRequest, "Article slug is required")
		return
	}

	article := h.cache.GetBySlug(slug)
	if article == nil {
		h.respondError(w, http.StatusNotFound, "Article not found")
		return
	}

	h.respondJSON(w, http.StatusOK, article)
}

// pubSubMessage represents the Pub/Sub push message envelope.
type pubSubMessage struct {
	Message struct {
		Attributes map[string]string `json:"attributes"`
		MessageID  string            `json:"messageId"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

// handlePubSubPush handles GCS notification push messages from Pub/Sub.
// It validates the shared token, filters non-.md files, and schedules a
// debounced cache refresh.
func (h *Handler) handlePubSubPush(w http.ResponseWriter, r *http.Request) {
	// 1. Validate the push identity if configured: Pub/Sub sends the
	// subscription's OIDC token as a Bearer token; verify it is signed by
	// Google, aimed at our audience, and issued to the expected service
	// account. The token value is never logged.
	if h.pushSAEmail != "" {
		raw := bearerToken(r)
		if raw == "" {
			h.logger.Warn("pubsub push rejected: missing bearer token")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		email, err := h.oidcValidate(r.Context(), raw, h.pushAudience)
		if err != nil || email != h.pushSAEmail {
			h.logger.Warn("pubsub push rejected: invalid push token")
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	// 2. Decode the Pub/Sub message.
	var msg pubSubMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		h.logger.Error("failed to decode pubsub message", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 3. Only refresh for .md file changes (ignore metadata.json writes, etc).
	objectID := msg.Message.Attributes["objectId"]
	eventType := msg.Message.Attributes["eventType"]

	if !strings.HasSuffix(objectID, ".md") {
		h.logger.Debug("pubsub notification ignored: not a .md file",
			"object", objectID, "event", eventType)
		w.WriteHeader(http.StatusOK)
		return
	}

	h.logger.Info("pubsub notification: triggering cache refresh",
		"object", objectID, "event", eventType, "messageId", msg.Message.MessageID)

	// 4. Schedule debounced refresh (returns immediately).
	h.cache.Refresh()

	// 5. Trigger frontend rebuild only for new/updated articles (not deletions).
	if eventType == "OBJECT_FINALIZE" {
		go h.triggerFrontendRebuild(objectID)
	}

	// 6. Acknowledge (200 OK). The actual refresh runs asynchronously.
	w.WriteHeader(http.StatusOK)
}

// triggerFrontendRebuild sends a POST request to the Cloud Build webhook URL
// to trigger a frontend-only rebuild so new articles are pre-rendered.
// This runs in a goroutine and never blocks the Pub/Sub acknowledgment.
func (h *Handler) triggerFrontendRebuild(objectID string) {
	if h.webhookURL == "" {
		return // No webhook configured, skip silently.
	}

	client := &http.Client{Timeout: 10 * time.Second}
	body := bytes.NewBufferString(fmt.Sprintf(`{"trigger_source":"pubsub","object":"%s"}`, objectID))

	resp, err := client.Post(h.webhookURL, "application/json", body)
	if err != nil {
		h.logger.Error("failed to trigger frontend rebuild webhook", "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		h.logger.Info("frontend rebuild triggered successfully", "object", objectID, "status", resp.StatusCode)
	} else {
		h.logger.Warn("frontend rebuild webhook returned non-2xx", "object", objectID, "status", resp.StatusCode)
	}
}

// handleArticlesStatus returns diagnostic information about all articles
// including skipped ones. Protected by a shared bearer token (Authorization
// header, never the URL) compared in constant time.
func (h *Handler) handleArticlesStatus(w http.ResponseWriter, r *http.Request) {
	if h.statusToken != "" {
		got := bearerToken(r)
		// Hash both sides first: ConstantTimeCompare alone leaks whether
		// the lengths match; comparing fixed-size digests hides the token
		// length too. Neither value is ever logged.
		gotSum := sha256.Sum256([]byte(got))
		wantSum := sha256.Sum256([]byte(h.statusToken))
		if subtle.ConstantTimeCompare(gotSum[:], wantSum[:]) != 1 {
			h.logger.Warn("articles status rejected: invalid token")
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	status := h.cache.GetCacheStatus()
	h.respondJSON(w, http.StatusOK, status)
}

// sitemapURLSet is the root element for a sitemap XML document.
type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}

// handleBlogSitemap generates a dynamic sitemap XML for all published articles.
func (h *Handler) handleBlogSitemap(w http.ResponseWriter, r *http.Request) {
	articles := h.cache.GetAllPublished()

	urls := make([]sitemapURL, 0, len(articles)+1)

	// Add the blog index page.
	urls = append(urls, sitemapURL{
		Loc:     "https://ivmanto.com/blog",
		LastMod: time.Now().UTC().Format("2006-01-02"),
	})

	for _, a := range articles {
		urls = append(urls, sitemapURL{
			Loc:     fmt.Sprintf("https://ivmanto.com/blog/%s", a.Slug),
			LastMod: a.Date,
		})
	}

	urlSet := sitemapURLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, xml.Header)
	if err := xml.NewEncoder(w).Encode(urlSet); err != nil {
		h.logger.Error("failed to write blog sitemap", "error", err)
	}
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			h.logger.Error("could not write JSON response", "error", err)
		}
	}
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"message": message})
}
