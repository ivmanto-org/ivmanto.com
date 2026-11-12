# Changelog

## v1.2.1 · 2026-11-12
- `pages/home.html` rebuilt to v1.2: approved copy, hero motion with pause button, compact package cards, icon diagram, dark band, tinted closing CTA. Template: `_build/templates/pages/home.html`.
- `components/ivmanto.css`: new blocks for hero motion (`.hero-motion`), the band (`.band`), the tinted section (`.section--tint`) and the icon diagram (`.flow`). Wave script: `_build/templates/_hero_motion.js`.
- Fix: inside the band `--accent-hover` maps to `--band-text`. The page's own hover colour (`#05636f`) is too dark on the navy band in the light theme.
- `components/layout.md`: the diagram line is `--border-input` (as in the approved preview), not `--border`, which disappears inside the band.
- `components.html`: "How it works" shows the icon diagram. It has no examples yet for the band or the compact package card.
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
