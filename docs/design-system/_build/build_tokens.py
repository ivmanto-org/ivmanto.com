"""Generate ivmanto design tokens (W3C JSON + CSS variables) and a contrast report from one source."""
import json, sys
from pathlib import Path

OUT = Path(sys.argv[1])

PRIMITIVES = {
    "teal": {"50": "#ecfbf8", "100": "#cff4ec", "200": "#a0e8db", "300": "#62d5c3", "400": "#02c39a",
             "500": "#00a896", "600": "#01948f", "700": "#027a89", "800": "#05636f", "900": "#0a4a55", "950": "#062f38"},
    "slate": {"50": "#f8fafc", "100": "#f1f5f9", "200": "#e2e8f0", "300": "#cbd5e1", "400": "#94a3b8",
              "500": "#64748b", "600": "#475569", "700": "#334155", "800": "#1e293b", "900": "#0f172a", "950": "#020617"},
    "navy": {"700": "#152238", "800": "#0f1a2b", "900": "#0a1220"},
    "dusk": {"100": "#e3eaf5", "300": "#9fb4d6", "500": "#4f6d9a", "700": "#33496b"},
    "amber": {"100": "#fef3c7", "400": "#fbbf24", "500": "#f59e0b", "700": "#b45309"},
    "green": {"100": "#e3f4ea", "300": "#7fd1a0", "700": "#2f7a52"},
    "red": {"100": "#fbe9e8", "300": "#f0a19c", "700": "#b4443f"},
}

# semantic name -> (light ref, dark ref, description)
SEMANTIC = {
    "bg":             ("slate.50",  "navy.900",  "Page background"),
    "surface":        ("#ffffff",   "navy.800",  "Cards, panels"),
    "surface-sunken": ("slate.100", "navy.700",  "Alternate sections, inputs"),
    "border":         ("slate.200", "#233149",   "Hairlines, card borders"),
    "border-input":   ("slate.500", "#6b7a93",   "Form control borders (WCAG 1.4.11, 3:1)"),
    "text":           ("slate.900", "#e6edf5",   "Headlines and body"),
    "text-muted":     ("slate.600", "#9aa8bc",   "Secondary text, captions"),
    "accent":         ("teal.700",  "teal.300",  "Links, primary buttons, focus, brand marks"),
    "accent-hover":   ("teal.800",  "teal.200",  "Hover/pressed of accent"),
    "on-accent":      ("#ffffff",   "teal.950",  "Text on accent fill"),
    "accent-subtle":  ("teal.50",   "#0d2b33",   "Tinted backgrounds, selected states"),
    "brand":          ("teal.500",  "teal.400",  "Logo cell, decorative only (not text on light)"),
    "data":           ("dusk.500",  "dusk.300",  "Diagrams, charts, secondary icons"),
    "signal":         ("amber.700", "amber.400", "Rare highlight: 'recommended' only"),
    "signal-subtle":  ("amber.100", "#2e2410",   "Background behind signal"),
    "success":        ("green.700", "green.300", "Low risk, done"),
    "success-subtle": ("green.100", "#12291f",   ""),
    "warning":        ("amber.700", "amber.400", "Medium risk"),
    "danger":         ("red.700",   "red.300",   "High risk, errors"),
    "danger-subtle":  ("red.100",   "#2e1615",   ""),
    "focus":          ("teal.500",  "teal.300",  "Focus ring (2px + 2px offset)"),
}

def resolve(ref):
    if ref.startswith("#"):
        return ref
    hue, step = ref.split(".")
    return PRIMITIVES[hue][step]

def as_token_value(ref):
    return ref if ref.startswith("#") else "{color.%s}" % ref

def lum(hex_):
    def ch(c):
        c = c / 255
        return c / 12.92 if c <= 0.03928 else ((c + 0.055) / 1.055) ** 2.4
    r, g, b = (int(hex_[i:i + 2], 16) for i in (1, 3, 5))
    return 0.2126 * ch(r) + 0.7152 * ch(g) + 0.0722 * ch(b)

def contrast(a, b):
    la, lb = sorted((lum(a), lum(b)), reverse=True)
    return round((la + 0.05) / (lb + 0.05), 2)

def build_color_json():
    color = {hue: {s: {"$type": "color", "$value": v} for s, v in steps.items()} for hue, steps in PRIMITIVES.items()}
    themes = {}
    for idx, theme in enumerate(("light", "dark")):
        themes[theme] = {name: {"$type": "color", "$value": as_token_value(refs[idx]),
                                **({"$description": refs[2]} if refs[2] else {})}
                         for name, refs in SEMANTIC.items()}
    return {"$description": "ivmanto colour tokens. Primitives under color.*, semantic roles under theme.light / theme.dark. Use semantic roles in UI code.",
            "color": color, "theme": themes}

def build_css():
    lines = ["/* Generated from color.tokens.json. Do not edit by hand. */", ":root {"]
    for hue, steps in PRIMITIVES.items():
        lines += [f"  --color-{hue}-{s}: {v};" for s, v in steps.items()]
    lines.append("")
    lines += [f"  --{n}: {resolve(r[0])};" for n, r in SEMANTIC.items()]
    lines += ["  color-scheme: light;", "}", "", ":root[data-theme=\"dark\"] {"]
    lines += [f"  --{n}: {resolve(r[1])};" for n, r in SEMANTIC.items()]
    lines += ["  color-scheme: dark;", "}", "", "@media (prefers-color-scheme: dark) {", "  :root:not([data-theme=\"light\"]) {"]
    lines += [f"    --{n}: {resolve(r[1])};" for n, r in SEMANTIC.items()]
    lines += ["    color-scheme: dark;", "  }", "}", ""]
    return "\n".join(lines)

PAIRS = [("text", "bg"), ("text", "surface"), ("text-muted", "bg"), ("text-muted", "surface"),
         ("accent", "bg"), ("accent", "surface"), ("on-accent", "accent"), ("data", "surface"),
         ("signal", "surface"), ("success", "surface"), ("danger", "surface"), ("border-input", "surface")]

def contrast_report():
    rows = []
    for fg, bg in PAIRS:
        for idx, theme in enumerate(("light", "dark")):
            c = contrast(resolve(SEMANTIC[fg][idx]), resolve(SEMANTIC[bg][idx]))
            rows.append((theme, fg, bg, c))
    return rows

if __name__ == "__main__":
    (OUT / "color.tokens.json").write_text(json.dumps(build_color_json(), indent=2) + "\n")
    (OUT / "tokens.css").write_text(build_css())
    (OUT / "contrast.json").write_text(json.dumps(
        [{"theme": t, "fg": f, "bg": b, "ratio": c} for t, f, b, c in contrast_report()], indent=2) + "\n")
    for t, f, b, c in contrast_report():
        need = 3.0 if f in ("border-input", "data") else 4.5
        print(f"{t:5} {f:14} on {b:8} {c:5}  {'OK' if c >= need else 'FAIL'} (need {need})")
