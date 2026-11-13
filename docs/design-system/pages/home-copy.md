# Home page copy (v1.1, approved by Nick 2026-11-11)

The "text cut": same sections and links, about half the visible words. Detail lives one click away.
Only the copy that changed is listed; everything else stays as built on `dev-v0.2.0`.

| # | Section | Element | Copy / rule |
|---|---|---|---|
| 1 | Hero | Lead | Fixed-scope services for the EU AI Act and ISO/IEC 42001. |
| 3 | Pillars | Lead | removed |
| 3 | Pillars | Governance tagline | Who decides what about AI. |
| 3 | Pillars | Risk management tagline | Find and reduce your AI risks. |
| 3 | Pillars | Compliance tagline | Know your obligations. Prove you meet them. |
| 3 | Pillars | Package lists in the cards | removed (home only) |
| 3 | Pillars | Layout | numbered index 01 to 03 (v1.3) that opens in place (v1.4): click shows the pillar `summary`, its packages and "More about <pillar>"; label "Packages and details" / "Close" |
| 4 | Packages | Lead | removed |
| 4 | Packages | Cards | compact ticket (home only); "What’s included" opens the `includes` list and a "Not included" line in place (v1.4); see `components/cards.md` |
| 4b | Engineering | Data & AI card | Data platforms and AI agents, built in your own Google Cloud project. |
| 4b | Engineering | ivGIL card | Our AI governance platform. Early access is open. |
| 5 | How it works | Describe | Ask the advisor or book a call. |
| 5 | How it works | Choose | Fixed scope and price, up front. |
| 5 | How it works | Deliver | Workshops and documents, built with you. |
| 5 | How it works | Hand over | Everything stays in your tools. |
| 6 | About | Text | I'm Nikolay Tonev. You work with me directly, from the first call to the handover. |
| 6 | About | Credentials | the sentence is replaced by badges: `CredentialBadge` (MIT Professional Education, wording rules in `components/credential.md`) and the Google Cloud Professional Cloud Architect badge already used on /about, with the line "Certified 2019, 2022 and 2024" (Nick, 2026-11-11: the last certificate expired 2026-09-01, so the badge states history; never show it bare). Then a link "More about me" to /about. |
| 7 | Insights | Cards | date and title only, no summary |
| 8 | FAQ | Questions | "Do you also build data platforms?" removed (section 4b covers it) |
| 9 | Closing | Lead | A 30-minute call finds the first step. |

Catalogue: pillars get a new `tagline` field (the three lines above). `summary` stays and is used on the pillar pages.
The background text about Nick (still a placeholder) is written once, for /about.

## Hero layout on desktop (v1.3.1, decided by Nick 2026-11-12)
From 1000px up the headline must sit on **two lines**: "Use AI with confidence." and "Govern it with clarity." (option C in `../previews/hero-headline.html`).
- Columns: text `1.15fr`, advisor `1fr` (was `1fr` / `1.05fr`).
- Headline size on the home hero only: `clamp(2.25rem, calc(4.4vw - 7px), 3rem)`, so 48px at most (was 52px). Every other `h1` keeps the scale in `ivmanto.css`.
- Below 1000px nothing changes: one column, headline wraps as before.
