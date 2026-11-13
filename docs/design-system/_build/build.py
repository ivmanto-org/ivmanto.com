"""Build the self-contained HTML previews of the ivmanto design system.

Usage (from docs/design-system):  python3 _build/build.py
Regenerates tokens (build_tokens.py), logo and favicon SVGs (build_logo.py) and every preview page
from _build/templates/. Each preview inlines all CSS so it opens anywhere (Claude Desktop,
a browser, a file share) without a server.
"""
import base64
import re, subprocess, sys
from pathlib import Path

DS = Path(__file__).resolve().parent.parent
BUILD = DS / "_build"
sys.path.insert(0, str(BUILD))

subprocess.run([sys.executable, str(BUILD / "build_tokens.py"), str(DS / "tokens")], check=True, stdout=subprocess.DEVNULL)
import build_logo  # noqa: E402
build_logo.write_all(DS)

CSS = "\n".join((DS / p).read_text() for p in ("tokens/tokens.css", "tokens/foundation.css", "components/ivmanto.css"))

ICONS = {
    "chev": '<svg class="chev" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M6 3l5 5-5 5"/></svg>',
    "cart": '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M3 4h2l2.4 11.2a1 1 0 0 0 1 .8h8.8a1 1 0 0 0 1-.8L20 8H6.2"/><circle cx="9" cy="20" r="1.2"/><circle cx="17" cy="20" r="1.2"/></svg>',
    "phone": '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M4 5c0 8.3 6.7 15 15 15l1-4-4-2-2 2a11 11 0 0 1-6-6l2-2-2-4z"/></svg>',
    "mail": '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="3" y="5" width="18" height="14" rx="2"/><path d="M3 7l9 6 9-6"/></svg>',
    "calendar": '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="3" y="5" width="18" height="16" rx="2"/><path d="M3 10h18M8 3v4M16 3v4"/></svg>',
    "send": '<svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M5 12h14M13 6l6 6-6 6"/></svg>',
    "clock": '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></svg>',
    "doc": '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><path d="M6 3h8l4 4v14H6z"/><path d="M14 3v4h4M9 12h6M9 16h6"/></svg>',
    "theme": '<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><circle cx="12" cy="12" r="8"/><path d="M12 4a8 8 0 0 0 0 16z" fill="currentColor"/></svg>',
    "menu": '<svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" aria-hidden="true"><path d="M4 7h16M4 12h16M4 17h16"/></svg>',
}

ICONS["arrow"] = '<svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M5 12h14M13 6l6 6-6 6"/></svg>'
ICONS["heart"] = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M12 20.5l-7.7-7.7a4.5 4.5 0 0 1 6.4-6.4L12 7.7l1.3-1.3a4.5 4.5 0 0 1 6.4 6.4z"/></svg>'
ICONS["check"] = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M5 12.5l4.5 4.5L19 7.5"/></svg>'
ICONS["video"] = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="3" y="6" width="12" height="12" rx="2"/><path d="M15 10.5l6-3.5v10l-6-3.5z"/></svg>'
ICONS["user"] = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="8" r="4"/><path d="M4 20.5a8 8 0 0 1 16 0"/></svg>'
ICONS["chevleft"] = '<svg viewBox="0 0 16 16" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M10 3L5 8l5 5"/></svg>'
ICONS["chevright"] = '<svg viewBox="0 0 16 16" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M6 3l5 5-5 5"/></svg>'
ICONS["pause"] = '<svg class="i-pause" viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" aria-hidden="true"><path d="M9 6v12M15 6v12"/></svg>'
ICONS["play"] = '<svg class="i-play" viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linejoin="round" aria-hidden="true" hidden><path d="M8 5.5l10 6.5-10 6.5z"/></svg>'
# How-it-works step icons and pillar icons (v1.4): the files in assets/icons are the source
for _f in sorted((DS / "assets/icons").glob("*.svg")):
    ICONS[_f.stem.replace("-", "_")] = _f.read_text().strip().replace("<svg ", '<svg aria-hidden="true" ', 1)

def pillar_icon(letter):
    return ('<svg class="card__icon" viewBox="0 0 44 44" width="44" height="44" aria-hidden="true"><rect x="3" y="3" width="38" height="38" rx="9" fill="var(--accent-subtle)" '
            'stroke="var(--accent)" stroke-width="1.5"/><text x="22" y="27.5" text-anchor="middle" font-family="Montserrat, Inter, sans-serif" '
            f'font-size="15" font-weight="600" fill="var(--accent)">{letter}</text></svg>')

# ---- Logo (Nick 2026-11-10). Geometry lives in build_logo.py; colours come from the theme here. ----
L = build_logo

def _inline(cls, vb, inner, height, label="ivmanto"):
    x, y, w, h = vb
    aria = f'role="img" aria-label="{label}"' if label else 'aria-hidden="true"'
    return f'<svg class="{cls}" viewBox="{x:.1f} {y:.1f} {w:.1f} {h:.1f}" width="{w / h * height:.0f}" height="{height}" {aria}>{inner}</svg>'

INK, TEAL = "var(--text)", "var(--brand)"
def wordmark(h=24): return _inline("wordmark", L.WM_BOX, L.wordmark_inner(INK, TEAL), h)
def symbol(h=32): return _inline("symbol", L.SYMBOL_VIEWBOX, L.symbol_inner(INK, TEAL), h, label=None)
def lockup(h=32): return _inline("lockup", L.lockup_geometry()[1], L.lockup_inner(INK, TEAL), h)
def short(h=32):
    x, y, s = L.square(L.WM_I, L.WM_TICK)
    return _inline("short", (x, y, s, s), L.short_inner(INK, TEAL), h, label=None)

# Advisor avatar: teal disc with the tick in white
AVATAR = (f'<svg class="avatar" viewBox="0 0 28 28" aria-hidden="true"><circle cx="14" cy="14" r="14" fill="var(--accent)"/>'
          f'<path fill="var(--on-accent)" {L.tick_transform(14.9, 14.2, 17.5)} d="{L.WM_TICK}"/></svg>')

PARTIALS = {
    "HEADER": (DS / "_build/templates/_header.html").read_text(),
    "FOOTER": (DS / "_build/templates/_footer.html").read_text(),
    "THEME_JS": (DS / "_build/templates/_theme.js").read_text(),
    "HERO_MOTION_JS": (DS / "_build/templates/_hero_motion.js").read_text(),
    "DISCLOSURE_JS": (DS / "_build/templates/_disclosure.js").read_text(),
}

HEAD = """<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<!-- Preview only: fonts from Google Fonts. The live site self-hosts them (GDPR). -->
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Montserrat:wght@500;600&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
<link href="https://fonts.googleapis.com/css2?family=Fraunces:opsz,wght@9..144,600&text=0123456789&display=swap" rel="stylesheet">
<style>
""" + CSS + "\n</style>"

def render(text):
    for _ in range(3):  # partials may contain placeholders
        text = text.replace("{{HEAD}}", HEAD)
        for k, v in PARTIALS.items():
            text = text.replace("{{" + k + "}}", v)
        text = re.sub(r"\{\{ICON:(\w+)\}\}", lambda m: ICONS[m.group(1)], text)
        text = re.sub(r"\{\{DATAURI:([\w./-]+)\}\}", lambda m: "data:image/png;base64," + base64.b64encode((DS / m.group(1)).read_bytes()).decode(), text)
        text = re.sub(r"\{\{PILLAR:(\w)\}\}", lambda m: pillar_icon(m.group(1)), text)
        text = text.replace("{{AVATAR}}", AVATAR)
        for key, fn, default in (("WORDMARK", wordmark, 24), ("SYMBOL", symbol, 32), ("LOCKUP", lockup, 32), ("SHORT", short, 32)):
            text = re.sub(r"\{\{" + key + r"(?::(\d+))?\}\}", lambda m, fn=fn, d=default: fn(int(m.group(1) or d)), text)
    return text

if __name__ == "__main__":
    for tpl in sorted((BUILD / "templates").rglob("*.html")):
        if tpl.name.startswith("_") or tpl.name == "foundations.html":
            continue
        out = DS / tpl.relative_to(BUILD / "templates")
        out.parent.mkdir(parents=True, exist_ok=True)
        html = render(tpl.read_text())
        left = re.findall(r"\{\{[^}]+\}\}", html)
        if left:
            sys.exit(f"{tpl.name}: unresolved placeholders {left}")
        out.write_text(html)
        print("built", out.relative_to(DS))
    # foundations page keeps its own data-driven template
    f = (BUILD / "templates/foundations.html").read_text()
    css_f = (DS / "tokens/tokens.css").read_text() + "\n" + (DS / "tokens/foundation.css").read_text()
    f = f.replace("/*__TOKENS__*/", css_f).replace("/*__COLOR_JSON__*/", (DS / "tokens/color.tokens.json").read_text()) \
         .replace("/*__CONTRAST_JSON__*/", (DS / "tokens/contrast.json").read_text())
    (DS / "foundations.html").write_text(render(f))
    print("built foundations.html")
