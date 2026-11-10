# Assets: ready to use

Approved by Nick on 2026-11-10. Every logo part comes in four colourways:

| Colourway | Use on | JPG background |
|---|---|---|
| `on-light` | white or light backgrounds | white |
| `on-dark` | navy or dark backgrounds | navy `#0a1220` |
| `mono-black` | single-colour print, fax, embossing | white |
| `mono-white` | single-colour on dark | navy |

PNGs have a transparent background. The number in a file name is the width in pixels.

```
assets/
├── logo/
│   ├── wordmark/   ivmanto with the tick-v; the primary logo     SVG · PNG 600/1200/2400 · JPG 2400
│   ├── symbol/     hexagon, wider upper chevron in teal          SVG · PNG 256/512/1024 · JPG 1024
│   ├── lockup/     symbol + wordmark                             SVG · PNG 1200/2400 · JPG 2400
│   └── short/      "iv" + tick, for square spaces (social avatar) SVG · PNG 512/1024 · JPG 1024
├── favicon/        favicon.svg (follows the OS light/dark), favicon.ico (16/32/48),
│                   favicon-16/32/48.png, apple-touch-icon.png (180), icon-192/512.png,
│                   icon-maskable-512.png, favicon-512.jpg, site.webmanifest
├── social/         og-image.png / .jpg, 1200 × 630, for link previews
└── source/         Tally's original wordmark files (w650). The build reads the light one.
```

## Which file where
- **Website header:** the lockup, `logo/lockup/ivmanto-lockup-on-light.svg`, 30 px high (the site inlines it, so it follows the theme). Footer: the wordmark.
- **Documents, slides, email signature:** the wordmark or lockup PNG at 1200 px.
- **LinkedIn / social profile picture:** `logo/short/ivmanto-short-on-light-1024.jpg`.
- **Browser tab:** the favicon set; the HTML is below.

```html
<link rel="icon" href="/favicon.ico" sizes="48x48">
<link rel="icon" href="/favicon.svg" type="image/svg+xml">
<link rel="apple-touch-icon" href="/apple-touch-icon.png">
<link rel="manifest" href="/site.webmanifest">
<meta property="og:image" content="https://ivmanto.com/og-image.png">
```

## Rules (full version in `../brand.html`)
- Never retype the word or redraw the symbol; use these files.
- The tick and the upper chevron are always brand teal (or the mono colour).
- Minimum size: wordmark 12 px high, lockup 24 px, symbol 16 px.
- ™ is fine now; ® only after registration.

## Updating
All files are generated. To change anything, edit `../_build/build_logo.py` or replace
`source/ivmanto-tickv_w650_light.svg`, then run:

```
python3 _build/build.py                 # SVGs + preview pages
node _build/export_assets.mjs           # PNG, JPG, ICO, OG image (needs puppeteer-core + Chromium, see the script header)
```

Superseded logo concepts, logotype studies and font samples are archived outside this folder
(Brio's workspace, `products/ivmanto/archive/design-system-2026-11-10/`).
