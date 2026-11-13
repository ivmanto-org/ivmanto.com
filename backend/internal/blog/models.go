package blog

import "time"

// ArticleMeta holds metadata for a single blog article.
// Fields are parsed from YAML frontmatter in .md files.
type ArticleMeta struct {
	Slug      string `json:"slug"      yaml:"-"`
	Title     string `json:"title"     yaml:"title"`
	Summary   string `json:"summary"   yaml:"summary"`
	Date      string `json:"date"      yaml:"date"`
	Published bool   `json:"published" yaml:"published"`

	// DigitalSourceType passes the IPTC digital-source-type frontmatter value
	// through untouched (spec insights-2026-11-12/spec-ai-provenance-marking.md:
	// the emitting frontend owns the vocabulary allowlist, so a new IPTC code
	// needs no backend change). Empty on non-AI articles.
	DigitalSourceType string `json:"digital_source_type,omitempty" yaml:"digital_source_type"`

	// Topic is the insights category used by the /blog topic filter.
	// Validated against validTopics at cache refresh; unknown values are
	// dropped (fail closed: the filter never sees an off-vocabulary topic).
	Topic string `json:"topic,omitempty" yaml:"topic"`

	// ReadingMinutes is derived from the article body (words/200, rounded
	// up, minimum 1), not authored in frontmatter.
	ReadingMinutes int `json:"reading_minutes,omitempty" yaml:"-"`
}

// validTopics is the closed vocabulary for the insights topic filter
// (governance, risk, compliance, engineering).
var validTopics = map[string]bool{
	"governance":   true,
	"risk":         true,
	"compliance":   true,
	"engineering":  true,
}

// Article is the full article with rendered HTML content.
type Article struct {
	ArticleMeta
	Content string `json:"content"`
}

// SkippedArticle records why an article was not published.
type SkippedArticle struct {
	Slug       string `json:"slug"`
	Reason     string `json:"reason"`
	HasTitle   bool   `json:"has_title"`
	HasSummary bool   `json:"has_summary"`
	HasDate    bool   `json:"has_date"`
}

// MetadataCache is the structure persisted to metadata.json in GCS.
type MetadataCache struct {
	GeneratedAt time.Time     `json:"generated_at"`
	Articles    []ArticleMeta `json:"articles"`
}
