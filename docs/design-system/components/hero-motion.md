# Hero motion (home)

Approved by Nick 2026-11-11 (variant C). Live preview: `../previews/hero-motion.html` (switch between Today, A, B, C; light and dark).
Applies to the **home hero only**. Pillar heroes and the closing CTA keep the still hex-grid texture.

Three layers, back to front: glow, hex grid with wave, content with a one-time entrance. The headline text never moves after the entrance.

## 1. Glow
- One absolutely placed layer behind the grid (`inset: -20%`, `pointer-events: none`, `aria-hidden`).
- Two radial gradients: `--brand` centred at 74% / 30% (38% × 42%), `--data` centred at 16% / 82% (34% × 40%), each fading to transparent at 70%.
- Strength: colour mixed at 16% (light) and 20% (dark) with transparent.
- Drift: `transform` only, 22 s, ease-in-out, alternate; translate at most 3%, scale 1 to 1.08.

## 2. Hex grid with wave
- Canvas behind the content, same cell geometry as the `.hexgrid` texture (pointy-top, radius 32px, 56px wide, rows 48px apart, odd rows shifted 28px).
- Base lines: slate at 9% (`#64748b` light, `#94a3b8` dark), drawn once per frame, three edges per cell.
- Wave: a straight crest travelling down-right (direction 0.8 / 0.6), 560px between crests, 9 s per crest. Cell intensity = crest³ × a slow ripple (0.7 to 1.0).
- Wave colour = `--accent`. Cell fill at most 7% (light) / 6% (dark); cell outline at most 42% (light) / 50% (dark).
- Mask: `radial-gradient(ellipse 95% 95% at 50% 40%, #000 25%, transparent 78%)`.
- Device pixel ratio capped at 2. The loop stops while the hero is off-screen or the tab is hidden.
- **Below 768px: no wave.** The still grid, the glow and the entrance remain (on a phone the text sits directly over the grid).

## 3. Entrance (once per page load)
CSS only, so the text is in the HTML and readable if scripts fail. Only `opacity` and `transform` change: no layout shift.

| Element | Effect | Duration | Delay |
|---|---|---|---|
| Overline | fade | 500ms | 0 |
| Headline line 1 | rises from below its own line (clipped) | 750ms | 120ms |
| Headline line 2 (teal) | same | 750ms | 320ms |
| Advisor card | fade + 14px lift | 700ms | 380ms |
| Lead | fade + 14px lift | 600ms | 640ms |
| Buttons and trust line | fade + 14px lift | 600ms | 800ms |

Easing `--ease`. Everything is still after about 1.4 s. Each headline line needs `overflow: hidden` with `padding-bottom: .14em; margin-bottom: -.14em` so descenders are not cut.

## Accessibility (required)
- `prefers-reduced-motion: reduce`: no entrance, no drift, no wave. The still grid and the still glow are shown.
- The glow and wave run longer than five seconds, so WCAG 2.2.2 requires a way to stop them: an icon button (`.icon-btn`, 40px) in the hero's bottom-right corner, label "Pause background motion" / "Play background motion", `aria-pressed`. The choice is stored as `ivm-motion` in localStorage (a preference, no consent needed, same as `ivm-theme`).
- Text contrast is unchanged: the wave fill stays under 7%, checked against `--text` and `--text-muted` in both themes.

## Not allowed
- Moving, waving or colour-cycling the headline or any other text.
- The wave or the glow anywhere outside the home hero without a new decision.
