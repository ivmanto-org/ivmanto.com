# Cards

Common: `--surface`, 1px `--border`, radius `--radius-lg`, padding 24px. Hover (only when the whole card is a link): border becomes `--border-input` plus `--shadow-md` (light theme only). The whole card is clickable through one link (`.card__link::after` overlay); never nest other links inside a clickable card.

## PillarIndex (home, v1.3)
Approved by Nick 2026-11-12 (concept B "Ticket" in `../previews/cards-alternatives.html`). Replaces the three pillar cards on the home page. Markup: `ol.pillar-index` > `li.pillar-index__item`.
- No card surface. Three columns (they stack below about 800px), each with a 2px top rule in `--text`.
- Numeral "01" to "03": Montserrat 600, 88px, outlined (1.5px `--accent` stroke, no fill), `aria-hidden`. The order is Governance, Risk management, Compliance.
- Title h3 (24px) with the one link (`.card__link`, whole column clickable), then the pillar's `tagline` in muted text.
- Arrow in a 44px circle, top right, decorative (`aria-hidden`).
- Hover and keyboard focus: the top rule turns `--accent` from the left (400ms), the numeral fills with `--accent`, the arrow circle fills. Nothing moves except the rule; no transitions under reduced motion.

## PillarCard
Props: `pillar` (governance|risk|compliance), `title`, `summary`, `services[]` (max 3), `href`.
- Icon: rounded square (radius 9) with the pillar initial (G/R/C) in Montserrat, `--accent` on `--accent-subtle`.
- Title h3 (Montserrat), summary in muted text, then up to three service names separated by hairlines.
- Used for the cross-links on pillar pages. **The home page uses PillarIndex instead (v1.3).**

## PackageCard = ticket (v1.3)
Approved by Nick 2026-11-12 (concept B). "Fixed scope, fixed price" drawn as a ticket: the scope above a perforation, the price on the stub. Markup: `.ticket` > `.ticket__body` + `.ticket__stub`. It replaces the earlier package card everywhere (home, /packages, pillar pages).

Props: `id`, `pillar`, `title`, `price` (EUR, net), `duration`, `deliverable`, `includes[]` (max 4), `recommended` (bool), `compact` (bool), `href`.
- Shell: `--surface`, 1px `--border-input`, radius `--radius-lg`, no padding of its own. Not a clickable card: the button is the only link.
- **Body** (padding 24px): pillar tag (+ recommended badge) → title h3 (20px) → ticks list (`includes`, full variant only) → spec list, pinned to the bottom of the body.
- **Spec list** (`dl.ticket__spec`): label in JetBrains Mono 12px uppercase, muted ("Duration", "You get"), value in 14px. "You get" is the main `deliverable`; leave the row out if there is none.
- **Perforation**: 2px dashed `--border-input` line between two half-circle notches (radius 13px) cut into the left and right edges. The notches take the colour of the ground behind the ticket (`--ticket-ground`: `--bg` by default, `--surface-sunken` inside a sunken section). On any other ground set `--ticket-ground`.
- **Stub** (fixed height 124px): price left (Montserrat 32px, one line; small "fixed, excl. VAT" below it), button right: "View" + chevron, with `aria-label="View package: <title>"`.
- `recommended`: `--accent` border, stub filled with `--accent` (text `--on-accent`), button on `--surface`, plus the badge. The amber appears only in the badge. Otherwise a secondary button.
- `compact`: no ticks list. Home page only; /packages and pillar pages show the full ticket.
- Tickets in one row have equal height; the stubs line up because the stub height is fixed.
- Data comes from the **service catalogue** (epic #112). The same object feeds pages, cards and the advisor, so card text and advisor recommendations can't drift.
- Prices are always shown. If a package has no fixed price, it doesn't belong in the package grid.
- The buy box on the package page keeps the plain card (`.card.package`); it is not a ticket.

## PostCard (insights)
Meta line (topic · reading time, uppercase 12px), Montserrat h3 title, one-line summary. On the home page (v1.1): meta line and title only, no summary. No images by default; an optional 16:9 image goes above the meta line.
