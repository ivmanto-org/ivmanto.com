# Layout and structure

## Grid
- Container 1200px; gutters 20px (mobile) / 32px (≥ 768px). Sections: 64px vertical (mobile) / 96px (desktop).
- Alternate sections between `--bg` and `--surface-sunken` instead of using dividers. Max line length 68ch.
- Home (v1.2): sections 4b and 5 form one dark band and the closing CTA is tinted, see `backgrounds.md`.

## SiteHeader
Sticky, 68px, surface at 92% opacity with a light blur, hairline bottom border.
Left: the lockup (symbol + wordmark), 30px high, inline SVG so it follows the theme (`assets/logo/lockup/`). Footer: the wordmark alone, 22px. Nav (≥ 960px): Services, Packages, Insights, About; active item gets a 2px `--accent` underline.
Right: theme toggle (icon button), "Book a call" (primary sm; moves into the menu below 480px), menu button (< 960px).

## SiteFooter
4 columns: brand + one-line description; Services; Company; Legal (Imprint, Privacy, Terms, Cookie settings, **AI use on this site**). Bottom line: © year, "Information on this site is not legal advice."

## Hex-grid texture
A faint hexagon line pattern (SVG, slate at 18% opacity, radial mask), **only** behind the home hero, pillar heroes and the closing CTA. Never behind body text columns or forms. On the home hero it is drawn on a canvas and carries the wave (`hero-motion.md`); everywhere else it stays still.

## HowItWorks (`.flow`)
4 steps as a diagram (v1.1): each step is an icon in a 48px circle (1px `--border-input`, icon 24px in `--accent`, files `assets/icons/step-*.svg`), joined by a 1px `--border-input` line through the circle centres; the first circle is filled `--accent-subtle`. Under each: title (h4) and one line of at most seven words (`pages/home-copy.md`). Below 768px the steps stack and the line runs vertically on the left. The mono numbers 01–04 are gone. Static, no animation.

## Faq
Native `<details>/<summary>`, chevron rotates on open. One question per line; answers ≤ 3 sentences.

## Theme
Light is the default. The OS preference is respected on first visit; the toggle stores `ivm-theme`. Implement with `data-theme` on `<html>`, set before first paint (inline script in `app.vue` head) to avoid a flash.
