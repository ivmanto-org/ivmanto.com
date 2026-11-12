package advisor

import (
	"bytes"
	"os"
	"testing"
)

func TestEmbeddedCatalogueMatchesSiteCatalogue(t *testing.T) {
	// Arrange: the site's single source of truth, relative to this package.
	site, err := os.ReadFile("../../../data/catalogue.json")
	if err != nil {
		t.Fatalf("read site catalogue: %v", err)
	}

	// Assert
	if !bytes.Equal(site, embeddedCatalogue) {
		t.Fatal("backend/internal/advisor/catalogue_embedded.json differs from data/catalogue.json; copy the file again")
	}
}

func TestEmbeddedCatalogueParsesAndKnowsItsPackages(t *testing.T) {
	c, err := EmbeddedCatalogue()

	if err != nil {
		t.Fatalf("EmbeddedCatalogue: %v", err)
	}
	if !c.ValidID("comp-eu-ai-act-readiness") {
		t.Error("embedded catalogue should contain comp-eu-ai-act-readiness")
	}
	if c.ValidID("no-such-package") {
		t.Error("unknown id must not validate")
	}
}

func TestEmbeddedCataloguePromptContextCarriesNoPrices(t *testing.T) {
	c, err := EmbeddedCatalogue()
	if err != nil {
		t.Fatalf("EmbeddedCatalogue: %v", err)
	}

	if ctx := c.promptContext(); bytes.Contains([]byte(ctx), []byte("price")) {
		t.Errorf("prompt context must not mention prices: %s", ctx)
	}
}
