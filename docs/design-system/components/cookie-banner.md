# Cookie banner

Bottom-right card, max 560px. v1.4 (2026-11-13): two buttons, compact phone variant.
- First layer: **Reject · Accept**, equal size and visual weight (Reject is not a text link). This follows the German DSK guidance and TTDSG §25. There is no Settings button: the site has one optional category (analytics), so Reject and Accept are the whole choice. If a second optional category is ever added, Settings comes back.
- Text: "We use analytics cookies only if you agree. You can change this any time under “Cookie settings” in the footer. Details".
- **Phones (below 600px), compact variant:** docked edge to edge at the bottom, no radius, 16px padding, no heading, text shortened to "We use analytics cookies only if you agree. Details" (the second sentence is in `.cookie__long`, hidden here). Buttons side by side, equal weight. The dialog keeps `aria-label="Cookie consent"`. Target height at 390px: about 130 to 140px.
- No analytics, GTM or advisor funnel events before Accept; the existing gating stays as Clauco described.
- "Cookie settings" in the footer reopens the banner.
