# Section backgrounds (home)

Approved by Nick 2026-11-12 (option C, "Dark band"). Live preview: `../previews/home-backgrounds.html` (Today, A, B, C; light and dark).
The idea: the page stays flat and calm, with **one** strong change of pace halfway down and a tinted ending.

## What changes on the home page
| Section | Background |
|---|---|
| Hero | as in `hero-motion.md` |
| Frameworks, Pillars, Packages, About, Insights, FAQ | unchanged: flat `--bg` / `--surface-sunken`, alternating |
| 4b "We build what we govern" + 5 "How it works" | **one band** (`.band` wraps both sections) |
| Closing CTA | `--accent-subtle` fill, one soft glow from the top, the still hex-grid texture as before |

## The band
- One wrapper around both sections; the sections inside have no background of their own, and 4b has no bottom padding (the two read as one block).
- Fill `--band-bg`. Inside the band the semantic tokens are re-pointed, so cards, text, tags and the flow diagram need no special styles:

| Token inside `.band` | Light theme | Dark theme |
|---|---|---|
| `--bg`, `--surface-sunken` | `#0a1220` | `#062f38` |
| `--surface` | `#0f1a2b` | `#0a3f4a` |
| `--border` | `#233149` | `#14606b` |
| `--border-input` | `#6b7a93` | `#4f8f98` |
| `--text` | `#e6edf5` | `#e6edf5` |
| `--text-muted` | `#9aa8bc` | `#a9c4c8` |
| `--accent` | `#62d5c3` | `#a0e8db` |
| `--accent-subtle` | `#0d2b33` | `#05636f` |
| `--on-accent` | `#062f38` | `#062f38` |
| `--accent-hover` | `#e6edf5` (= band text) | `#e6edf5` |

  In the light theme the band uses the dark theme's own palette (a navy band on a light page). In the dark theme it turns deep teal, so it still stands out from the navy page.
- Glow: `--brand` at 26%, centred 88% / 12% (40% × 55%); `--data` at 22%, centred 6% / 92% (36% × 50%); both fade to transparent at 70%. Still, no drift.
- Hex-grid texture, still, top right: mask `radial-gradient(ellipse 60% 70% at 85% 20%, #000, transparent 72%)`.

## Closing CTA
Fill `--accent-subtle`; glow `--brand` at the hero's glow strength (16% light, 20% dark), centred 50% / 0% (50% × 90%).

## Contrast (WCAG AA, computed 2026-11-12)
| Pair | Light theme | Dark theme |
|---|---|---|
| Band text on band fill | 15.9 | 12.1 |
| Band muted text on band fill / on band card | 7.8 / 7.2 | 7.8 / 6.3 |
| Band accent on band fill / on its subtle fill | 10.6 / 8.4 | 10.2 / 5.0 |
| Closing: text / muted text on `--accent-subtle` | 16.8 / 7.1 | 12.6 / 6.2 |

The band values are tokens: `--band-bg`, `--band-surface`, `--band-border`, `--band-border-input`, `--band-text`, `--band-text-muted`, `--band-accent`, `--band-accent-subtle`, `--band-on-accent` (`tokens/color.tokens.json`, `tokens/tokens.css`). `.band` maps each to the role of the same name, for example `--text: var(--band-text)`. All pairs above are in `tokens/contrast.json`. Reference CSS: `.band` and `.section--tint` in `components/ivmanto.css`; built example: `pages/home.html`. The preview file uses two older names (`--band-muted`, `--band-border-strong`); the token files are the reference.

## Rules
- One band per page. Home only until Nick decides otherwise.
- Nothing in the band or the closing section moves. Motion stays in the hero.
- No amber inside the band.
- The band holds no form and no long text.

## Not chosen
A "Soft light" (fading grey sections with corner glows) and B "Grid echo" (hex grid in section corners, teal hairlines).
