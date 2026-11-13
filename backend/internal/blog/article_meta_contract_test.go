package blog

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"testing"
)

// fakeStorage is an in-memory Storage for cache-level contract tests.
type fakeStorage struct {
	files map[string][]byte
}

func (f *fakeStorage) ListMarkdownFiles(ctx context.Context) ([]string, error) {
	names := make([]string, 0, len(f.files))
	for name := range f.files {
		names = append(names, name)
	}
	return names, nil
}

func (f *fakeStorage) ReadFile(ctx context.Context, name string) ([]byte, error) {
	data, ok := f.files[name]
	if !ok {
		return nil, io.ErrUnexpectedEOF
	}
	return data, nil
}

func (f *fakeStorage) WriteMetadataCache(ctx context.Context, cache *MetadataCache) error {
	return nil
}

func articleFrontmatter(extra string) string {
	return "---\ntitle: Test\ndate: 2026-11-13\npublished: true\n" + extra + "---\n\n# Head\n\nBody text.\n"
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// TestParse_DigitalSourceType_Passthrough pins the spec-3 contract: the
// backend passes the IPTC frontmatter value through untouched (vocabulary
// enforcement lives in the emitting frontend), and non-AI articles omit it.
func TestParse_DigitalSourceType_Passthrough(t *testing.T) {
	p := NewParser()

	a, err := p.Parse([]byte(articleFrontmatter("digital_source_type: trainedAlgorithmicMedia\n")), "slug")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if a.DigitalSourceType != "trainedAlgorithmicMedia" {
		t.Errorf("digital_source_type = %q, want trainedAlgorithmicMedia (raw passthrough)", a.DigitalSourceType)
	}

	// An unknown code must also pass through untouched, not be dropped here.
	b, err := p.Parse([]byte(articleFrontmatter("digital_source_type: someFutureCode\n")), "slug")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if b.DigitalSourceType != "someFutureCode" {
		t.Errorf("digital_source_type = %q, want raw someFutureCode (frontend owns the allowlist)", b.DigitalSourceType)
	}

	// Absent field omits from JSON (omitempty).
	c, err := p.Parse([]byte(articleFrontmatter("")), "slug")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	j, _ := json.Marshal(c.ArticleMeta)
	if strings.Contains(string(j), "digital_source_type") {
		t.Errorf("digital_source_type should be omitted when absent, got %s", j)
	}
}

// TestParse_ReadingMinutes_Arithmetic pins words/200 rounded up, min 1.
func TestParse_ReadingMinutes_Arithmetic(t *testing.T) {
	p := NewParser()

	words := func(n int) string {
		return "---\ntitle: T\ndate: 2026-11-13\npublished: true\n---\n\n" + strings.Repeat("word ", n)
	}

	cases := []struct {
		name  string
		words int
		want  int
	}{
		{"empty body floors at 1", 0, 1},
		{"200 words is 1 minute", 200, 1},
		{"201 words rounds up to 2", 201, 2},
		{"400 words is 2", 400, 2},
		{"4500 words is 23", 4500, 23},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a, err := p.Parse([]byte(words(tc.words)), "slug")
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if a.ReadingMinutes != tc.want {
				t.Errorf("ReadingMinutes = %d, want %d", a.ReadingMinutes, tc.want)
			}
		})
	}
}

// TestParse_ReadingMinutes_TagsDoNotInflateCount guards against HTML markup
// counting as words.
func TestParse_ReadingMinutes_TagsDoNotInflateCount(t *testing.T) {
	p := NewParser()
	src := "---\ntitle: T\ndate: 2026-11-13\npublished: true\n---\n\n" +
		strings.Repeat("<div class=\"x\">word</div> ", 10)
	a, err := p.Parse([]byte(src), "slug")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if a.ReadingMinutes != 1 {
		t.Errorf("ReadingMinutes = %d, want 1 (10 markup-wrapped words)", a.ReadingMinutes)
	}
}

// TestCacheRefresh_TopicVocabulary pins that valid topics survive to the API
// meta list and off-vocabulary topics are dropped (fail closed).
func TestCacheRefresh_TopicVocabulary(t *testing.T) {
	storage := &fakeStorage{files: map[string][]byte{
		"good.md": []byte(articleFrontmatter("topic: governance\n")),
		"bad.md":  []byte(articleFrontmatter("topic: philosophy\n")),
		"none.md": []byte(articleFrontmatter("")),
	}}

	cache, err := NewCache(context.Background(), storage, NewParser(), testLogger())
	if err != nil {
		t.Fatalf("new cache: %v", err)
	}
	defer cache.Stop()

	all := cache.GetAllPublished()
	bySlug := map[string]ArticleMeta{}
	for _, m := range all {
		bySlug[m.Slug] = m
	}

	if got := bySlug["good"].Topic; got != "governance" {
		t.Errorf("good.md topic = %q, want governance", got)
	}
	if got := bySlug["bad"].Topic; got != "" {
		t.Errorf("bad.md topic = %q, want \"\" (off-vocabulary dropped)", got)
	}
	if got := bySlug["none"].Topic; got != "" {
		t.Errorf("none.md topic = %q, want \"\"", got)
	}
	if bySlug["good"].ReadingMinutes < 1 {
		t.Errorf("good.md reading_minutes = %d, want >= 1", bySlug["good"].ReadingMinutes)
	}

	// The single-article endpoint carries the same fields.
	full := cache.GetBySlug("good")
	if full == nil {
		t.Fatal("GetBySlug(good) = nil")
	}
	if full.Topic != "governance" || full.ReadingMinutes < 1 {
		t.Errorf("GetBySlug fields: topic=%q reading_minutes=%d", full.Topic, full.ReadingMinutes)
	}
}

// TestCacheRefresh_AllTopicsValid pins the full closed vocabulary.
func TestCacheRefresh_AllTopicsValid(t *testing.T) {
	files := map[string][]byte{}
	for _, topic := range []string{"governance", "risk", "compliance", "engineering"} {
		files[topic+".md"] = []byte(articleFrontmatter("topic: "+topic+"\n"))
	}
	cache, err := NewCache(context.Background(), &fakeStorage{files: files}, NewParser(), testLogger())
	if err != nil {
		t.Fatalf("new cache: %v", err)
	}
	defer cache.Stop()

	for _, m := range cache.GetAllPublished() {
		if !validTopics[m.Topic] {
			t.Errorf("topic %q should have survived as valid", m.Topic)
		}
	}
	if got := len(cache.GetAllPublished()); got != 4 {
		t.Errorf("published count = %d, want 4", got)
	}
}
