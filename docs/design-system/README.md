# ivmanto design system

Owner: Brio (design). Implementation: Clauco (Nuxt 3 + Tailwind). Review: Ivmo. Sign-off: Nick.
Status: **v1.0, approved by Nick** (2026-11-10); Later versions (now v1.4.0: home hero motion, text cut, dark band, ticket package cards, pillars and tickets that open in place) are in `CHANGELOG.md` and `DECISIONS.md`. Handed to Clauco for the new site.

**What is decided, what is a placeholder, and what was rejected: see [`DECISIONS.md`](DECISIONS.md).** It wins over any other file or chat.

## Start here
| Open | What |
|---|---|
| `brand.html` | **The logo**: wordmark, symbol, lockup, short form, favicon, rules, file list |
| `foundations.html` | Palette, contrast checks, type scale |
| `components.html` | Component library, including all advisor states |
| `pages/home.html` | The home page with the AI advisor hero |
| `pages/pillar.html`, `pages/package.html`, `pages/order.html` | Remaining page templates |
| `DECISIONS.md` | **Every decision**: approved, placeholder, rejected |
| `assets/` | **Ready-to-use files**: logo (SVG/PNG/JPG), favicon set, OG image. See `assets/README.md` |

Every preview has a light/dark toggle. Specs: `components/*.md`, `pages/README.md`.

## Folder structure
```
design-system/
├── README.md, DECISIONS.md, CHANGELOG.md
├── brand.html, foundations.html, components.html   previews (generated)
├── pages/          page templates (generated) + README with page specs
├── tokens/         design tokens: JSON (source of truth) + CSS
├── components/     ivmanto.css + one spec per component (*.md)
├── assets/         logo, favicon, social image (ready to use)
└── _build/         generators and HTML templates (edit these, not the generated pages)
```

## Who it's for
C-level executives at SMEs: CEOs, Chief AI Officers, Chief Data Officers. They have little time,
they are wary of hype, and they want to know what they get, from whom, by when, and at what price.

## Principles
1. **Calm precision.** Trust comes from structure and restraint. Flat surfaces, 1px borders, soft
   shadows only on hover. No glass, gradient text or moving text. One exception, approved by Nick
   2026-11-11: the home hero has a soft glow, a one-time entrance and a slow wave in the hexagon
   grid (`components/hero-motion.md`). Nowhere else. Since 2026-11-12 the home page also has one dark
   band and a tinted closing section, both still (`components/backgrounds.md`).
2. **One idea, used consistently.** The tick (approved, checked) and the hexagon with rising
   chevrons are the only motifs: logo, avatar, list ticks, funnel steps, a faint hexagon grid at most.
   They are not decoration.
3. **Answer first.** Every section opens with the outcome in one sentence. Scope, deliverable,
   duration and price come before the explanation.
4. **One accent.** Teal does the work. Amber appears once per page at most ("recommended").
   Status colours appear only where a risk level is shown.
5. **Accessible by default.** WCAG 2.2 AA in both themes (checked in `tokens/contrast.json`).
   Visible focus rings, `prefers-reduced-motion` respected, minimum body text 16px (18px on
   marketing pages).

## Voice (AI governance / GRC)
- Plain, specific, unexcited. "We document owners, controls and evidence for each AI system",
  not "unlock responsible innovation".
- Avoid the category clichés: agentic, AI-ready, one platform, unlock, supercharge, fear copy.
- Name the frameworks exactly: ISO/IEC 42001, EU AI Act, NIST AI RMF, GDPR. No claims of
  certification, partnership or results we can't prove (UWG). No vendor comparisons.
- Scope (Nick, 2026-11-10 late): the site sells AI governance **services** delivered by Nick, is a selling point for
  **ivGIL**, and has a section for Nick's **GCP Data & AI Architect / Data Management** work. The information
  architecture for the three is being proposed (see `DECISIONS.md`, section 2).

## Tokens
| File | Contents |
|---|---|
| `tokens/color.tokens.json` | W3C design tokens. Primitives (`color.teal.500`…) plus semantic roles for `theme.light` and `theme.dark`. **The source of truth.** |
| `tokens/tokens.css` | CSS variables generated from the JSON. Light by default, dark via `[data-theme="dark"]` or the OS preference. |
| `tokens/foundation.tokens.json` / `foundation.css` | Typography, spacing (4px base), radii, shadows, motion. |
| `tokens/contrast.json` | Computed contrast ratios for the key pairs. |

**Rule for code:** components use **semantic roles only** (`var(--accent)`, `var(--surface)`),
never primitives or hex values. Then the dark theme comes for free.

### Colour decisions (v0.1)
- **Brand teal kept:** `#00a896` (500), `#02c39a` (400). `#028090` is adjusted to **`#027a89`**
  (700), visually the same, so text and buttons on light backgrounds pass AA (5.1:1 on white).
- **Navy kept:** `#0f172a` for text on light; the dark theme uses deeper navy `#0a1220`.
- **Added:**
  - a cool slate neutral scale;
  - `dusk` (muted slate-blue) for diagrams and data;
  - calm green, amber and red for risk levels.
- **Amber demoted** from accent to rare signal.
- **Removed:** the glow, glass and gradient tokens from the old site.

## Decisions and open items
All in [`DECISIONS.md`](DECISIONS.md): approved, changed, team rules, placeholders, open, rejected.

## Build
Everything except the specs is generated:
```
python3 _build/build.py          # tokens, logo + favicon SVGs, all preview pages
node _build/export_assets.mjs    # PNG, JPG, ICO and OG image from the SVGs
```
Edit `_build/templates/`, `_build/build_logo.py`, `components/ivmanto.css` or the token JSON, never the
generated HTML or asset files.

## Handoff to Clauco
Map `tokens.css` roles into `tailwind.config.js` as `colors: { bg: 'var(--bg)', … }` and set
`darkMode: ['selector', '[data-theme="dark"]']`. Self-host the fonts (e.g. `@fontsource`): Montserrat 500/600, Inter 400/500/600, JetBrains Mono 400. The logo needs no font; it is outlined SVG.
Don't use the Google Fonts CDN. Details come with the component specs.
