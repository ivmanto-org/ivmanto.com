// Staging server: serves the generated static site behind HTTP Basic auth and
// keeps it out of search engines. It has no other job; the production site
// does not use it.
package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	defaultPort    = "8080"
	defaultSiteDir = "/site"
	notFoundPage   = "404.html"
	robotsBody     = "User-agent: *\nDisallow: /\n"
	authRealm      = `Basic realm="ivmanto staging", charset="UTF-8"`
)

// parseUsers reads "user:password" entries separated by newlines or commas.
// One entry per group of viewers, so access can be withdrawn per group.
func parseUsers(raw string) (map[string][32]byte, error) {
	users := make(map[string][32]byte)
	for _, entry := range strings.FieldsFunc(raw, func(r rune) bool { return r == '\n' || r == ',' }) {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		name, password, found := strings.Cut(entry, ":")
		if !found || name == "" || password == "" {
			return nil, errors.New("each entry must be user:password")
		}
		users[name] = sha256.Sum256([]byte(password))
	}
	if len(users) == 0 {
		return nil, errors.New("no users configured")
	}
	return users, nil
}

// authorised compares in constant time and does the same work for unknown users.
func authorised(users map[string][32]byte, name, password string) bool {
	want, known := users[name]
	got := sha256.Sum256([]byte(password))
	match := subtle.ConstantTimeCompare(got[:], want[:]) == 1
	return known && match
}

func newHandler(users map[string][32]byte, site fs.FS) http.Handler {
	files := http.FileServerFS(site)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Robots-Tag", "noindex, nofollow, noarchive")
		if r.URL.Path == "/robots.txt" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = w.Write([]byte(robotsBody))
			return
		}
		name, password, ok := r.BasicAuth()
		if !ok || !authorised(users, name, password) {
			w.Header().Set("WWW-Authenticate", authRealm)
			http.Error(w, "Sign in to view the ivmanto staging site.", http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "This staging site is read-only.", http.StatusMethodNotAllowed)
			return
		}
		// Viewers should see a new deploy without clearing their cache.
		w.Header().Set("Cache-Control", "no-cache")
		if !exists(site, r.URL.Path) {
			serveNotFound(w, site)
			return
		}
		files.ServeHTTP(w, r)
	})
}

func exists(site fs.FS, urlPath string) bool {
	name := strings.TrimPrefix(path.Clean("/"+urlPath), "/")
	if name == "" {
		name = "."
	}
	info, err := fs.Stat(site, name)
	if err != nil {
		return false
	}
	if !info.IsDir() {
		return true
	}
	_, err = fs.Stat(site, path.Join(name, "index.html"))
	return err == nil
}

func serveNotFound(w http.ResponseWriter, site fs.FS) {
	page, err := fs.ReadFile(site, notFoundPage)
	if err != nil {
		http.Error(w, "Page not found.", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write(page)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// Fail closed: without credentials the site must not be served at all.
	users, err := parseUsers(os.Getenv("STAGING_USERS"))
	if err != nil {
		logger.Error("STAGING_USERS is missing or invalid; refusing to start", "error", err.Error())
		os.Exit(1)
	}
	siteDir := envOr("SITE_DIR", defaultSiteDir)
	if _, err := os.Stat(path.Join(siteDir, "index.html")); err != nil {
		logger.Error("site directory has no index.html; refusing to start", "dir", siteDir, "error", err.Error())
		os.Exit(1)
	}
	addr := ":" + envOr("PORT", defaultPort)
	logger.Info("staging server starting", "addr", addr, "groups", len(users))
	if err := http.ListenAndServe(addr, newHandler(users, os.DirFS(siteDir))); err != nil {
		logger.Error("server stopped", "error", err.Error())
		os.Exit(1)
	}
}
