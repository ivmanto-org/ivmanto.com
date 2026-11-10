# Tags and badges

| Component | Use | Look |
|---|---|---|
| `Tag variant="risk" level="low|medium|high"` | Only where a risk level is actually shown (advisor answers, risk register samples) | dot + label; `--success`/`--warning`/`--danger` on their `-subtle` bg |
| `Tag variant="framework"` | Frameworks we work with: EU AI Act, ISO/IEC 42001, ISO/IEC 27001, NIST AI RMF, GDPR | mono 13px, 1px `--border`, `--surface` |
| `Tag variant="pillar"` | Pillar label on package cards | `--accent` on `--accent-subtle` |
| `BadgeRecommended` | One per page at most: the recommended package | amber `--signal` on `--signal-subtle`, uppercase 12px |

Risk labels follow EU AI Act wording: Minimal / Limited / High risk (plus "Prohibited" in text only, never as a coloured tag).
Badge copy: "Recommended" or "Recommended start". **Never "Most popular" or "Best seller"**: we have no data to back it (UWG §5).
Framework tags are never "certified", "partner" or "compliant" claims.
