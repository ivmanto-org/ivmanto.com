# Cards

Common: `--surface`, 1px `--border`, radius `--radius-lg`, padding 24px. Hover (only when the whole card is a link): border becomes `--border-input` plus `--shadow-md` (light theme only). The whole card is clickable through one link (`.card__link::after` overlay); never nest other links inside a clickable card.

## PillarCard
Props: `pillar` (governance|risk|compliance), `title`, `summary`, `services[]` (max 3), `href`.
- Icon: rounded square (radius 9) with the pillar initial (G/R/C) in Montserrat, `--accent` on `--accent-subtle`.
- Title h3 (Montserrat), summary in muted text, then up to three service names separated by hairlines.
- **Home (v1.1):** icon, title and the pillar's `tagline` only (at most seven words). No summary, no service names. Pillar pages keep `summary`.

## PackageCard
Props: `id`, `pillar`, `title`, `price` (EUR, net), `duration`, `deliverable`, `includes[]` (max 4), `recommended` (bool), `href`.
- Order: pillar tag (+ recommended badge) → title → price (Montserrat 32px; small "fixed, excl. VAT") → meta row (clock + duration, doc + main deliverable) → ticks list → button (bottom-aligned).
- `recommended`: border `--signal`, badge, primary button. Otherwise a secondary button.
- Data comes from the **service catalogue** (epic #112). The same object feeds pages, cards and the advisor, so card text and advisor recommendations can't drift.
- Prices are always shown. If a package has no fixed price, it doesn't belong in the package grid.
- **`compact` (bool, v1.1):** hides the ticks list. Used on the home page only; /packages and pillar pages show the full card.

## PostCard (insights)
Meta line (topic · reading time, uppercase 12px), Montserrat h3 title, one-line summary. On the home page (v1.1): meta line and title only, no summary. No images by default; an optional 16:9 image goes above the meta line.
