# Layout and structure

## Grid
- Container 1200px; gutters 20px (mobile) / 32px (≥ 768px). Sections: 64px vertical (mobile) / 96px (desktop).
- Alternate sections between `--bg` and `--surface-sunken` instead of using dividers. Max line length 68ch.

## SiteHeader
Sticky, 68px, surface at 92% opacity with a light blur, hairline bottom border.
Left: the lockup (symbol + wordmark), 30px high, inline SVG so it follows the theme (`assets/logo/lockup/`). Footer: the wordmark alone, 22px. Nav (≥ 960px): Services, Packages, Insights, About; active item gets a 2px `--accent` underline.
Right: theme toggle (icon button), "Book a call" (primary sm; moves into the menu below 480px), menu button (< 960px).

## SiteFooter
4 columns: brand + one-line description; Services; Company; Legal (Imprint, Privacy, Terms, Cookie settings, **AI use on this site**). Bottom line: © year, "Information on this site is not legal advice."

## Hex-grid texture
A faint hexagon line pattern (SVG, slate at 18% opacity, radial mask), **only** behind the home hero, pillar heroes and the closing CTA. Never behind body text columns or forms.

## HowItWorks (`.flow`)
4 steps with mono numbers 01–04 and a top rule; the first step's rule is `--accent`. No arrows or connectors.

## Faq
Native `<details>/<summary>`, chevron rotates on open. One question per line; answers ≤ 3 sentences.

## Theme
Light is the default. The OS preference is respected on first visit; the toggle stores `ivm-theme`. Implement with `data-theme` on `<html>`, set before first paint (inline script in `app.vue` head) to avoid a flash.
