# Cards

Common: `--surface`, 1px `--border`, radius `--radius-lg`, padding 24px. Hover (only when the whole card is a link): border becomes `--border-input` plus `--shadow-md` (light theme only). The whole card is clickable through one link (`.card__link::after` overlay); never nest other links inside a clickable card.

## PillarIndex (home, v1.4: opens in place)
Approved by Nick 2026-11-13 (`../previews/open-in-place.html`); the numbered index itself is from v1.3. Flat markup inside `div.pillar-index`: `.pillar-index__head`, `.pillar-index__panel`, three times. Each head and its panel card carry `style="--i:0|1|2"`.
- **Head:** no card surface, 2px top rule in `--text`. Numeral "01" to "03" left, pillar icon (44px, `--accent`) right, then title h3 (24px), the pillar's `tagline`, and a mono label "Packages and details" with a chevron ("Close" when open; decorative, `aria-hidden`).
- **Numerals:** Fraunces 600, 88px, digits only (self-hosted subset on the site; OFL). Outline without inner contour lines: `-webkit-text-stroke: 3px var(--accent); paint-order: stroke fill; color: var(--ground)`. `--ground` is the colour behind the index (`--bg`; `--surface-sunken` inside a sunken section). Fraunces is used for these numerals only, nowhere else.
- **Control:** the title is a `button.pillar-index__btn` with `aria-expanded` and `aria-controls`; its hit area is the whole head. The pillar page is no longer one click from the head: the link is inside the panel.
- **Panel:** `role="region"`, labelled by the title. One panel is open at a time. From 900px it spans the full row under the three heads, with a small pointer under the open column; below 900px it opens directly under its head. Content: the pillar's `summary` (18px) and a ghost link "More about <pillar>" on the left; the pillar's packages on the right, each row = name (link, whole row clickable), "duration · deliverable" in muted text, arrow.
- **States:** hover, focus and open turn the top rule `--accent` from the left and fill the numeral. On hover or open the icon draws itself once (700ms).
- **Motion:** the panel opens over 360ms (grid rows), content fades up. Under reduced motion everything is instant and the icon does not draw. A closed panel is `visibility: hidden`, so it is out of the tab order.
- Script: `_build/templates/_disclosure.js` (no dependency).

### Pillar icons (v1.4)
`assets/icons/pillar-governance.svg` (decision tree), `pillar-risk.svg` (radar), `pillar-compliance.svg` (checklist). 24px grid, round caps, same family as the step icons. Shown at 44px with stroke 1.3 in the pillar index and at 14px with stroke 1.9 inside the pillar tag on a ticket. Every shape has `pathLength="1"` for the draw animation.

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
- `compact`: no ticks list on the face. Home page only; /packages and pillar pages show the full ticket.
- **Compact ticket opens in place (v1.4, approved by Nick 2026-11-13):** under the spec list a `button.ticket__toggle` "What’s included" with a plus in a 28px circle (`aria-expanded`, `aria-controls`). It unfolds `.ticket__more`: the ticks list (`includes`) and a muted "Not included" line (`not_included`). The plus turns into a cross. From 900px the tickets of one row open and close together, so the row stays level and the scopes can be compared; below 900px each ticket opens on its own. Reduced motion: instant.
- The pillar tag on a ticket carries the pillar icon (14px).
- Tickets in one row have equal height; the stubs line up because the stub height is fixed.
- Data comes from the **service catalogue** (epic #112). The same object feeds pages, cards and the advisor, so card text and advisor recommendations can't drift.
- Prices are always shown. If a package has no fixed price, it doesn't belong in the package grid.
- The buy box on the package page keeps the plain card (`.card.package`); it is not a ticket.

## PostCard (insights)
v1.5: one meta line (`TOPIC · 18 SEP 2026 · 6 MIN`, JetBrains Mono 12px uppercase, topic in `--accent`), Montserrat h3 title (18px; the only link, covers the card), summary clamped to two lines. Grid: `.post-grid` (3 columns from 1000px, 2 from 640px). On the home page (v1.1): date and title only, no summary. No images by default; an optional 16:9 image goes above the meta line. Full rules: `../pages/blog.md`.
