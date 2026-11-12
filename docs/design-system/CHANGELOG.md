# Changelog

## v1.3.1 · 2026-11-12 (decided by Nick)
- **Home headline on two lines on desktop** (option C of `previews/hero-headline.html`): text column `1.15fr`, advisor `1fr`; home hero headline at most 48px, `clamp(2.25rem, calc(4.4vw - 7px), 3rem)` from 1000px up. Rule in `pages/home-copy.md`; styles are page-local in `_build/templates/pages/home.html`. `pages/home.html` rebuilt.

## v1.3.0 · 2026-11-12 (approved by Nick)
- **Package card becomes a ticket** (concept B from `previews/cards-alternatives.html`): scope above a perforation, price on the stub; recommended = teal outline and filled stub. New `.ticket` block in `components/ivmanto.css`; spec in `components/cards.md`. Used on home (compact), pillar pages and /packages (full).
- **Home pillars become a numbered index** (`.pillar-index`): outlined numerals 01 to 03 that fill on hover or focus, no card surface. The pillar card stays for cross-links on pillar pages.
- The recommended package no longer has an amber border; the amber is in the badge only.
- New build icon `{{ICON:arrow}}`. Templates changed: `pages/home.html`, `pages/pillar.html`, `components.html`. The old `.package` styles remain for the buy box on the package page.
- New preview `previews/cards-alternatives.html` (the three concepts Nick chose from: A Cell, B Ticket, C Rows). Prices in it are placeholders.

## v1.2.3 · 2026-11-12 (pause button: decided by Nick)
- **Hero pause button moves to the bottom-left corner** of the home hero. Reason: on a first visit the cookie banner sits bottom right and covered it. `components/hero-motion.md`, `components/ivmanto.css` (`.hero-motion__pause`).
- **The pause button is hidden under `prefers-reduced-motion: reduce`**: nothing moves there, so there is nothing to pause.
- New preview `previews/hero-headline.html`: three ways to bring the home headline back to two lines on desktop (today it breaks into four from 1000 px up). **Not decided**; no spec or stylesheet change for the headline yet.
- All previews regenerated with the new stylesheet.

## v1.2.2 · 2026-11-12
- `components.html`: new examples for the compact package card, the dark band and the tinted section. No change to tokens, stylesheet or specs.
- Home about block: the Cloud Architect badge carries "certified 2019, 2022 and 2024" (`pages/home-copy.md`, `pages/home.html`, `DECISIONS.md`).

## v1.2.1 · 2026-11-12
- `pages/home.html` rebuilt to v1.2: approved copy, hero motion with pause button, compact package cards, icon diagram, dark band, tinted closing CTA. Template: `_build/templates/pages/home.html`.
- `components/ivmanto.css`: new blocks for hero motion (`.hero-motion`), the band (`.band`), the tinted section (`.section--tint`) and the icon diagram (`.flow`). Wave script: `_build/templates/_hero_motion.js`.
- Fix: inside the band `--accent-hover` maps to `--band-text`. The page's own hover colour (`#05636f`) is too dark on the navy band in the light theme.
- `components/layout.md`: the diagram line is `--border-input` (as in the approved preview), not `--border`, which disappears inside the band.
- `components.html`: "How it works" shows the icon diagram.
- All other previews regenerated with the new tokens and stylesheet; their content is unchanged.

## v1.2.0 · 2026-11-12 (approved by Nick)
- **Home backgrounds** (option C): dark band around sections 4b and 5, tinted closing CTA. New spec `components/backgrounds.md`, live preview `previews/home-backgrounds.html` (also shows the whole home page with the v1.1 copy).
- Band tokens (`--band-*`, nine roles, light and dark) added to `_build/build_tokens.py`, `tokens/color.tokens.json`, `tokens/tokens.css`; their contrast pairs are in `tokens/contrast.json`.

## v1.1.0 · 2026-11-11 (approved by Nick)
- **Home hero motion** (variant C): glow, one-time entrance, hex-grid wave. New spec `components/hero-motion.md`, live preview `previews/hero-motion.html`.
- **Home text cut**: approved copy and rules in `pages/home-copy.md`.
- `components/cards.md`: PillarCard tagline on home, PackageCard `compact`, PostCard without summary on home.
- `components/layout.md`: HowItWorks becomes an icon diagram; new icons `assets/icons/step-{describe,choose,deliver,handover}.svg`.
- README principle 1 and the motion rule in `components/README.md` now name the hero exception.
- `DECISIONS.md` updated. Not rebuilt yet: `pages/home.html` and `components.html` still show v1.0.1.

## v1.0.1 · 2026-11-11
- `DECISIONS.md` added: one log of every decision (approved, changed, team rules, placeholders, open, rejected).
- Price placeholder is now `€ [price]` (the repo's dash check bans em dashes in site copy).
- Governance pillar preview matches the catalogue: "AI policy & acceptable use" (3 weeks), "Governance operating model" (6 weeks).
- `pages/README.md`: `data/catalogue.json` in the site repo is the single source for package content.
- `components/advisor.md`: privacy-notice status (draft with Ivgil) and the backend facts it relies on.
- README: site scope widened by Nick (services + ivGIL + GCP Data & AI section); the decisions section moved to `DECISIONS.md`.

## v1.0 · 2026-11-10 (approved by Nick)
- Header shows the lockup (symbol + wordmark, 30px); footer the wordmark.
- Logo approved: tick-v wordmark (Tally), hexagon symbol with the wider upper chevron, lockup, short form "iv" + tick.
- `assets/` holds ready-to-use files: SVG, PNG, JPG in four colourways; favicon set (SVG, ICO, PNG, apple-touch, manifest icons); OG image.
- Headlines in Montserrat; the serif token is removed.
- Pillar icons are rounded squares; the advisor avatar is the tick on a teal disc.
- Removed from this folder (archived in Brio's workspace): logo concepts A/B/C, the logotype studies, the font comparison page and the commercial font samples.

## v0.1 · 2026-11-09
- First complete draft: tokens (light/dark), components, page templates, logo concepts.
