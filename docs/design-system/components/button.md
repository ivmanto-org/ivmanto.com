# Button (`Btn`)

| Variant | Use | Tokens |
|---|---|---|
| `primary` | The one main action per view (Order, Send, Book a call in the header) | bg `--accent`, text `--on-accent`, hover `--accent-hover` |
| `secondary` | Alternatives next to a primary | bg `--surface`, border `--border-input`, text `--text` |
| `ghost` | Tertiary and inline actions ("How it works ›") | text `--accent`, hover bg `--accent-subtle` |

- Size: default min-height 44px, padding 10/20, `--text-md` 500, radius `--radius-md`. `sm` is 36px (header only, desktop).
- Optional trailing chevron icon (`.chev`); it moves 2px on hover. No leading icons except in the advisor send button.
- At most **one primary per section**. In card grids, only the recommended card gets primary; the others get secondary.
- Disabled: 50% opacity plus `aria-disabled`. Loading: the label changes to "Sending…" and the button stays the same width.
- Never: gradients, glows, uppercase labels, pill radius.
