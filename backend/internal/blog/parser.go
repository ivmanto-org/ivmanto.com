package blog

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gmparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
)

// Parser converts markdown content to HTML and extracts frontmatter metadata.
type Parser struct {
	md goldmark.Markdown
}

// NewParser creates a new markdown parser with frontmatter and GFM support.
func NewParser() *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			&frontmatter.Extender{},
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	return &Parser{md: md}
}

// Parse takes raw markdown bytes and returns extracted metadata and rendered HTML.
// The slug is derived from the filename, not from frontmatter.
func (p *Parser) Parse(source []byte, slug string) (*Article, error) {
	var buf bytes.Buffer

	ctx := gmparser.NewContext()
	if err := p.md.Convert(source, &buf, gmparser.WithContext(ctx)); err != nil {
		return nil, fmt.Errorf("goldmark convert: %w", err)
	}

	var meta ArticleMeta
	d := frontmatter.Get(ctx)
	if d != nil {
		if err := d.Decode(&meta); err != nil {
			return nil, fmt.Errorf("frontmatter decode: %w", err)
		}
	}

	meta.Slug = slug

	meta.ReadingMinutes = readingMinutes(buf.String())

	return &Article{
		ArticleMeta: meta,
		Content:     buf.String(),
	}, nil
}

// readingMinutes estimates reading time from a rendered HTML body:
// words/200 rounded up, minimum 1 (a body with zero words still counts
// as one minute). Tags are stripped so markup never inflates the count.
func readingMinutes(htmlBody string) int {
	var text strings.Builder
	inTag := false
	for _, r := range htmlBody {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag && r != '\n':
			text.WriteRune(r)
		}
	}
	words := len(strings.Fields(text.String()))
	minutes := (words + 199) / 200
	if minutes < 1 {
		minutes = 1
	}
	return minutes
}
