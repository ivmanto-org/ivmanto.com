package advisor

import (
	"encoding/json"
	"fmt"
	"os"
)

// Catalogue is the in-memory form of data/catalogue.json — the single
// source for packages on the site and for the advisor (#112). The advisor
// only ever recommends package ids that exist here.
type Catalogue struct {
	Version  int       `json:"version"`
	Currency string    `json:"currency"`
	Pillars  []Pillar  `json:"pillars"`
	Packages []Package `json:"packages"`
}

// Pillar is a top-level service pillar.
type Pillar struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Summary string `json:"summary"`
}

// Package is one bookable service package.
type Package struct {
	ID       string `json:"id"`
	Pillar   string `json:"pillar"`
	Name     string `json:"name"`
	Outcome  string `json:"outcome"`
	PriceEur *int64 `json:"price_eur"`
	// PillarRecommended is Brio's optional per-package flag (set on a few
	// packages once #126's catalogue corrections land). Ignored by the
	// advisor logic today; captured so the embedded-copy test tolerates it
	// and a later slice can use it without a schema change.
	PillarRecommended *bool `json:"pillar_recommended"`
}

// DefaultCataloguePath is where the backend looks for the shared service
// catalogue when CATALOGUE_PATH is unset. Producing a build that carries it
// next to the binary is a deploy concern (Cloud Run working dir); env var
// overrides for other layouts.
const DefaultCataloguePath = "data/catalogue.json"

// LoadCatalogue reads and parses the catalogue JSON file at path.
func LoadCatalogue(path string) (*Catalogue, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read catalogue: %w", err)
	}
	var c Catalogue
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse catalogue %s: %w", path, err)
	}
	if len(c.Packages) == 0 {
		return nil, fmt.Errorf("catalogue %s contains no packages", path)
	}
	return &c, nil
}

// ValidID reports whether id is a package id present in the catalogue.
func (c *Catalogue) ValidID(id string) bool {
	for _, p := range c.Packages {
		if p.ID == id {
			return true
		}
	}
	return false
}

// promptContext is the compact, model-facing view of the catalogue embedded
// in the system prompt: id, pillar, name and outcome only. Prices are
// deliberately excluded — the model must never quote prices (the catalogue
// is a placeholder and prices render from catalogue data by id).
func (c *Catalogue) promptContext() string {
	type entry struct {
		ID      string `json:"id"`
		Pillar  string `json:"pillar"`
		Name    string `json:"name"`
		Outcome string `json:"outcome"`
	}
	entries := make([]entry, 0, len(c.Packages))
	for _, p := range c.Packages {
		entries = append(entries, entry{ID: p.ID, Pillar: p.Pillar, Name: p.Name, Outcome: p.Outcome})
	}
	b, err := json.Marshal(entries)
	if err != nil {
		// entry fields are all strings; Marshal cannot fail here.
		return "[]"
	}
	return string(b)
}
