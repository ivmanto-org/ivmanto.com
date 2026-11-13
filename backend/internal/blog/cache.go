package blog

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"
	"time"
)

const defaultRefreshInterval = 5 * time.Minute

const defaultDebounceDuration = 5 * time.Second

// Cache holds the in-memory article cache with periodic refresh from GCS.
type Cache struct {
	storage Storage
	parser  *Parser
	logger  *slog.Logger

	mu         sync.RWMutex
	articles   map[string]*Article // keyed by slug
	metaList   []ArticleMeta       // published only, sorted by date desc
	skipped    []SkippedArticle    // articles skipped during last refresh
	totalFiles int                 // total .md files found in last refresh

	stopCh chan struct{}

	// Debounce: coalesce multiple rapid Refresh() calls into one.
	debounceMu    sync.Mutex
	debounceTimer *time.Timer
	debounceDur   time.Duration
}

// NewCache creates a new blog cache, performs the initial load, and starts
// a background goroutine that refreshes the cache periodically.
func NewCache(ctx context.Context, storage Storage, parser *Parser, logger *slog.Logger) (*Cache, error) {
	c := &Cache{
		storage:     storage,
		parser:      parser,
		logger:      logger.With("service", "blog_cache"),
		articles:    make(map[string]*Article),
		stopCh:      make(chan struct{}),
		debounceDur: defaultDebounceDuration,
	}

	if err := c.refresh(ctx); err != nil {
		return nil, fmt.Errorf("initial blog cache load: %w", err)
	}

	go c.backgroundRefresh()

	return c, nil
}

// GetAllPublished returns metadata for all published articles, sorted newest first.
func (c *Cache) GetAllPublished() []ArticleMeta {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]ArticleMeta, len(c.metaList))
	copy(result, c.metaList)
	return result
}

// CacheStatus holds diagnostic information about the article cache.
type CacheStatus struct {
	TotalFiles int              `json:"total_files"`
	Published  int              `json:"published"`
	Skipped    []SkippedArticle `json:"skipped"`
}

// GetCacheStatus returns diagnostic information about all articles,
// including those that were skipped during the last refresh.
func (c *Cache) GetCacheStatus() CacheStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()

	skipped := make([]SkippedArticle, len(c.skipped))
	copy(skipped, c.skipped)

	return CacheStatus{
		TotalFiles: c.totalFiles,
		Published:  len(c.metaList),
		Skipped:    skipped,
	}
}

// GetBySlug returns a full article (metadata + HTML) by slug, or nil if not found.
func (c *Cache) GetBySlug(slug string) *Article {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.articles[slug]
}

// refresh re-reads all .md files from GCS, parses them, and updates the cache.
func (c *Cache) refresh(ctx context.Context) error {
	files, err := c.storage.ListMarkdownFiles(ctx)
	if err != nil {
		return fmt.Errorf("listing markdown files: %w", err)
	}

	newArticles := make(map[string]*Article, len(files))
	var newMeta []ArticleMeta
	var newSkipped []SkippedArticle

	for _, filename := range files {
		data, err := c.storage.ReadFile(ctx, filename)
		if err != nil {
			c.logger.Error("failed to read article file", "file", filename, "error", err)
			continue
		}

		slug := strings.TrimSuffix(filename, ".md")
		article, err := c.parser.Parse(data, slug)
		if err != nil {
			c.logger.Error("failed to parse article", "file", filename, "error", err)
			newSkipped = append(newSkipped, SkippedArticle{
				Slug:   slug,
				Reason: "parse error: " + err.Error(),
			})
			continue
		}

		// Warn about missing required fields for debugging.
		if article.Title == "" {
			c.logger.Warn("article missing title", "slug", slug, "file", filename)
		}
		if article.Summary == "" {
			c.logger.Warn("article missing summary", "slug", slug, "file", filename)
		}
		if article.Date == "" {
			c.logger.Warn("article missing date", "slug", slug, "file", filename)
		}

		if !article.Published {
			c.logger.Info("skipping unpublished article (published != true)", "slug", slug)
			newSkipped = append(newSkipped, SkippedArticle{
				Slug:       slug,
				Reason:     "published field is false or missing",
				HasTitle:   article.Title != "",
				HasSummary: article.Summary != "",
				HasDate:    article.Date != "",
			})
			continue
		}

		newArticles[slug] = article
		newMeta = append(newMeta, article.ArticleMeta)
	}

	// Sort by date descending (newest first).
	sort.Slice(newMeta, func(i, j int) bool {
		return newMeta[i].Date > newMeta[j].Date
	})

	c.mu.Lock()
	c.articles = newArticles
	c.metaList = newMeta
	c.skipped = newSkipped
	c.totalFiles = len(files)
	c.mu.Unlock()

	c.logger.Info("blog cache refreshed", "article_count", len(newMeta))

	// Write derived metadata.json to GCS (best-effort).
	cache := &MetadataCache{
		GeneratedAt: time.Now().UTC(),
		Articles:    newMeta,
	}
	if err := c.storage.WriteMetadataCache(ctx, cache); err != nil {
		c.logger.Error("failed to write metadata cache to GCS", "error", err)
	}

	return nil
}

// Refresh schedules a cache refresh with debouncing. Multiple calls within
// the debounce window (5s) are coalesced into a single GCS read. This is the
// public entry point for Pub/Sub push notifications.
func (c *Cache) Refresh() {
	c.debounceMu.Lock()
	defer c.debounceMu.Unlock()

	if c.debounceTimer != nil {
		c.debounceTimer.Stop()
	}

	c.debounceTimer = time.AfterFunc(c.debounceDur, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := c.refresh(ctx); err != nil {
			c.logger.Error("pubsub-triggered cache refresh failed", "error", err)
		}
	})

	c.logger.Debug("cache refresh debounced/scheduled")
}

// backgroundRefresh runs refresh on a ticker until Stop is called.
func (c *Cache) backgroundRefresh() {
	ticker := time.NewTicker(defaultRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			if err := c.refresh(ctx); err != nil {
				c.logger.Error("background cache refresh failed", "error", err)
			}
			cancel()
		case <-c.stopCh:
			return
		}
	}
}

// Stop shuts down the background refresh goroutine.
func (c *Cache) Stop() {
	close(c.stopCh)
}
