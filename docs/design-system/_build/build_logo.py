"""ivmanto logo geometry and the vector logo assets.

Approved by Nick (2026-11-10):
  * symbol   : hexagon with two rising chevrons, the upper one wider (concept C, revision 2)
  * wordmark : "ivmanto" with the tick-v, by Tally. Instrument Sans 650, outlined (SIL OFL 1.1)
  * favicon  : short form "iv" + tick (Tally's design)
  * lockup   : symbol + wordmark
  * short    : "iv" + tick, for square spaces (social avatar, app icon); the favicon uses it

Sources: the symbol is drawn here; the wordmark is read from assets/source/ivmanto-tickv_w650_light.svg
(Tally's file). Replace that file with a polished version and rebuild: every asset and page follows.

Usage:  python3 _build/build_logo.py <design-system dir>   (build.py runs it)
Writes the SVGs into assets/logo/ and assets/favicon/. PNG, JPG and ICO come from export_assets.mjs.
"""
import math, re, sys
from pathlib import Path

DS = Path(__file__).resolve().parent.parent
SRC = DS / "assets/source/ivmanto-tickv_w650_light.svg"

# Colours: the brand tokens (tokens/color.tokens.json). "on-light" = for light backgrounds.
COLOURS = {
    "on-light": {"ink": "#0f172a", "teal": "#00a896", "bg": "#ffffff"},
    "on-dark": {"ink": "#e6edf5", "teal": "#2fd0b8", "bg": "#0a1220"},
    "mono-black": {"ink": "#000000", "teal": "#000000", "bg": "#ffffff"},
    "mono-white": {"ink": "#ffffff", "teal": "#ffffff", "bg": "#000000"},
}

# ---- symbol: hexagon + rising chevrons (64 x 64 grid) ------------------------------------------
SYMBOL_VIEWBOX = (8, 4.75, 48, 54.5)  # the hexagon including its stroke, no extra margin


def _hex(cx, cy, r):
    return " ".join(f"{cx + r * math.cos(math.radians(-90 + 60 * i)):.2f},{cy + r * math.sin(math.radians(-90 + 60 * i)):.2f}" for i in range(6))


def symbol_inner(ink, teal):
    """Upper chevron has longer legs than the lower one (same angle), so the centre reads as rising."""
    slope = 8 / 11
    out = [f'<polygon points="{_hex(32, 32, 25)}" fill="none" stroke="{ink}" stroke-width="4.5" stroke-linejoin="round"/>']
    for apex_y, half, col in ((22.5, 13, teal), (34.5, 9.5, ink)):
        y = apex_y + half * slope
        out.append(f'<polyline points="{32 - half:.2f},{y:.2f} 32,{apex_y} {32 + half:.2f},{y:.2f}" fill="none" stroke="{col}" '
                   'stroke-width="4.5" stroke-linecap="round" stroke-linejoin="round"/>')
    return "".join(out)


# ---- wordmark: Tally's tick-v outlines ----------------------------------------------------------
_paths = re.findall(r'<path fill="(#[0-9a-fA-F]{6})"[^>]*\sd="([^"]+)"', SRC.read_text())
WM_INK = next(d for c, d in _paths if c.lower() == "#0f172a")
WM_TICK = next(d for c, d in _paths if c.lower() == "#00a896")
WM_I = " ".join(re.findall(r"M[^M]*", WM_INK)[:2])  # the i: stem + dot, the first two sub-paths


def _nums(d):
    """Points of an absolute SVG path (M L H V Q Z, as in the wordmark file); control points included."""
    xs, ys, x, y = [], [], 0.0, 0.0
    for cmd, args in re.findall(r"([MLHVQZ])([^MLHVQZ]*)", d):
        n = [float(v) for v in re.findall(r"-?\d+(?:\.\d+)?", args)]
        if cmd == "H":
            for v in n: x = v; xs.append(x); ys.append(y)
        elif cmd == "V":
            for v in n: y = v; xs.append(x); ys.append(y)
        elif cmd in "MLQ":
            for i in range(0, len(n), 2):
                x, y = n[i], n[i + 1]; xs.append(x); ys.append(y)
    return xs, ys


def box(*ds, pad=0):
    xs, ys = [], []
    for d in ds:
        x, y = _nums(d); xs += x; ys += y
    return min(xs) - pad, min(ys) - pad, max(xs) - min(xs) + 2 * pad, max(ys) - min(ys) + 2 * pad


def square(*ds, pad=60):
    x, y, w, h = box(*ds)
    s = max(w, h) + 2 * pad
    return x + w / 2 - s / 2, y + h / 2 - s / 2, s


WM_BOX = box(WM_INK, WM_TICK, pad=8)
SQ = square(WM_I, WM_TICK, pad=0)[2]


def wordmark_inner(ink, teal):
    return f'<path fill="{ink}" d="{WM_INK}"/><path fill="{teal}" d="{WM_TICK}"/>'


def short_inner(ink, teal):
    return f'<path fill="{ink}" d="{WM_I}"/><path fill="{teal}" d="{WM_TICK}"/>'


def tick_transform(cx, cy, height):
    tx, ty, tw, th = box(WM_TICK)
    k = height / th
    return f'transform="translate({cx - (tx + tw / 2) * k:.2f} {cy - (ty + th / 2) * k:.2f}) scale({k:.5f})"'


# ---- lockup: symbol left of the wordmark, in wordmark units --------------------------------------
_, TICK_TOP, _, TICK_H = box(WM_TICK)
LOCKUP_SYMBOL_H = TICK_H * 1.45            # symbol a little taller than the tick
LOCKUP_K = LOCKUP_SYMBOL_H / SYMBOL_VIEWBOX[3]
LOCKUP_GAP = LOCKUP_SYMBOL_H * 0.30


def lockup_geometry():
    sx, sy, sw, sh = SYMBOL_VIEWBOX
    wx, wy, ww, wh = WM_BOX
    cy = TICK_TOP + TICK_H / 2 + 20                 # optical centre: middle of the tick, nudged to the x-height
    sym_w = sw * LOCKUP_K
    left = wx - LOCKUP_GAP - sym_w
    top = cy - LOCKUP_SYMBOL_H / 2
    tf = f'transform="translate({left - sx * LOCKUP_K:.2f} {top - sy * LOCKUP_K:.2f}) scale({LOCKUP_K:.4f})"'
    vb = (left, min(top, wy), wx + ww - left, max(top + LOCKUP_SYMBOL_H, wy + wh) - min(top, wy))
    return tf, vb


def lockup_inner(ink, teal):
    tf, _ = lockup_geometry()
    return f'<g {tf}>{symbol_inner(ink, teal)}</g>{wordmark_inner(ink, teal)}'


# ---- files -----------------------------------------------------------------------------------------
def _svg(vb, inner, title="ivmanto", pad=0.0):
    x, y, w, h = vb
    x, y, w, h = x - pad * h, y - pad * h, w + 2 * pad * h, h + 2 * pad * h
    return (f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="{x:.2f} {y:.2f} {w:.2f} {h:.2f}" '
            f'width="{w / h * 256:.0f}" height="256" role="img"><title>{title}</title>{inner}</svg>\n')


def favicon_svg():
    x, y, s = square(WM_I, WM_TICK)
    return (f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="{x:.0f} {y:.0f} {s:.0f} {s:.0f}"><title>ivmanto</title>'
            '<style>.i{fill:#0f172a}.t{fill:#00a896}@media (prefers-color-scheme:dark){.i{fill:#e6edf5}.t{fill:#2fd0b8}}</style>'
            f'<path class="i" d="{WM_I}"/><path class="t" d="{WM_TICK}"/></svg>\n')


def write_all(ds=DS):
    logo, fav = ds / "assets/logo", ds / "assets/favicon"
    for sub in ("symbol", "wordmark", "lockup", "short"):
        (logo / sub).mkdir(parents=True, exist_ok=True)
    fav.mkdir(parents=True, exist_ok=True)
    _, lvb = lockup_geometry()
    for name, c in COLOURS.items():
        ink, teal = c["ink"], c["teal"]
        (logo / f"symbol/ivmanto-symbol-{name}.svg").write_text(_svg(SYMBOL_VIEWBOX, symbol_inner(ink, teal)))
        (logo / f"wordmark/ivmanto-wordmark-{name}.svg").write_text(_svg(WM_BOX, wordmark_inner(ink, teal)))
        (logo / f"lockup/ivmanto-lockup-{name}.svg").write_text(_svg(lvb, lockup_inner(ink, teal)))
        (logo / f"short/ivmanto-short-{name}.svg").write_text(_svg(square(WM_I, WM_TICK, pad=0)[:2] + (SQ, SQ), short_inner(ink, teal)))
    (fav / "favicon.svg").write_text(favicon_svg())


if __name__ == "__main__":
    write_all(Path(sys.argv[1]) if len(sys.argv) > 1 else DS)
