# Decision log: ivmanto design system

The single record of what Nick has decided, what is still a placeholder, and what was rejected.
If any other file, chat or agent disagrees with this log, **this log wins**. Newest first within each section.
Owner: Brio. Last updated 2026-11-11.

## 1. Approved (use as is)

| Date | Decision | Where |
|---|---|---|
| 2026-11-10 | **Design system v1.0 approved** for the new site. Clauco builds it (PR #122 docs, #126 / draft PR #129 site, branch `dev-v0.2.0`). | whole folder |
| 2026-11-10 | **Logo: four parts.** (1) Wordmark "ivmanto" with the **tick-v** (Tally, Instrument Sans 650, outlined, SIL OFL). (2) **Symbol**: hexagon with two rising chevrons, the **upper teal chevron wider** (concept C, revision 2). (3) **Lockup**: symbol + wordmark. (4) **Short form** "iv" + tick. | `brand.html`, `assets/logo/` |
| 2026-11-10 | **Website header = lockup, 30 px. Footer = wordmark alone.** | `_build/templates/_header.html`, `_footer.html` |
| 2026-11-10 | **Favicon = short form "iv" + tick** (Tally's design), as SVG (follows OS light/dark), ICO, PNG and JPG. | `assets/favicon/` |
| 2026-11-10 | **Advisor avatar** = white tick on a teal disc. It is never used as the logo. | `components/advisor.md` |
| 2026-11-10 | **Four colourways** for every logo part: `on-light`, `on-dark`, `mono-black`, `mono-white`. | `assets/README.md` |
| 2026-11-10 | **Fonts:** Montserrat for headlines (`--font-display`), Inter for UI and body, JetBrains Mono for code. Self-hosted, no Google Fonts CDN. The logo needs no font file. | `tokens/foundation.*` |
| 2026-11-09 | **Light theme is the default**; dark theme is available via the toggle and the OS preference. | `tokens/tokens.css` |
| 2026-11-09 | **Palette:** brand teal `#00a896`, text teal `#027a89` (AA on white), navy `#0f172a`, dark background `#0a1220`. Amber appears at most once per page ("Recommended"). | `tokens/color.tokens.json` |
| 2026-11-09 | **AI service advisor is the home-page hero** (epic #112). Grounded in the catalogue, AI-labelled in every state, "not legal advice" next to the input. | `components/advisor.md`, `pages/home.html` |
| 2026-11-09 | Calm, trustworthy look: flat surfaces, 1 px borders, no glow, glass, gradients or endless animation. WCAG 2.2 AA in both themes. | `README.md` principles |

## 2. Changed by Nick, details pending

| Date | Change | Status |
|---|---|---|
| 2026-11-10 (late, via Clauco) | **Site scope widened.** The site sells AI governance services **and** is a selling point / tool for **ivGIL**, **and** has a separate section for Nick's **GCP Data & AI Architect / Data Management** skills. This reverses "the site sells services only, not ivGIL". | Brio prepares an information-architecture proposal for Nick's content session (2026-11-11). Until then, page templates show the services part only. |

## 3. Agreed working rules (team, not Nick sign-off)

| Date | Rule | Agreed by |
|---|---|---|
| 2026-11-11 | **Price placeholder is `€ [price]`**, not `€ —`: the repo's dash check (`scripts/check-dashes.mjs`) bans em dashes in site copy. It must never go live; prices come from Nick. | Brio, Clauco |
| 2026-11-10 | **`data/catalogue.json` (site repo) is the single source** for package names, outcomes, durations, includes, prices and the "recommended" flags. Pages and the advisor read it; the previews here only illustrate layout. | Brio, Clauco, Ivmo |
| 2026-11-10 | Catalogue corrections from Brio applied: AI policy & acceptable use **3 weeks**; Governance operating model **6 weeks**; `pillar_recommended` marks each pillar's start (policy, AI risk assessment, EU AI Act readiness); home shows one overall "Recommended start" (EU AI Act readiness); "audit-ready", never "audit-proof". | Clauco (dev-v0.2.0) |
| 2026-11-10 | Advisor backend facts the privacy text relies on: no chat text stored, random session ID, no message content in logs (the request logger keeps IP, URL and status), rate-limit counters 24 h, leads go by email to Nick only, no hidden fields. | Ivmo (#123, #124) |
| 2026-11-11 | **Ivgil verdict on the advisor privacy notice + /ai page: allow-with-obligation.** Before publishing: every `[bracket]` filled from verified fact; the AI provider's DPA checked (training, retention, EU location); Ivgil writes the Art. 6(1)(f) balancing tests; a lawyer reads the final privacy texts; Nick approves. One legal basis per purpose (lead form = Art. 6(1)(b)). | Ivgil |
| 2026-11-11 | Catalogue wording (Ivgil): the incident package says "handle AI incidents" and checks *whether* something is reportable (Art. 73 is a high-risk provider duty), never implies every client must report. The staff briefing *supports* the client's own Art. 4 duty and does not discharge it. | Ivgil |
| 2026-11-10 | Naming: **ivGIL** outward, "Ivgil" for the internal agent. No absolute "no one does X" claims. Comparative or legal-effect claims go to Ivgil before publishing. | Ivgil |

## 4. Placeholders: do not publish

| Item | Waiting for |
|---|---|
| Package names, scope, durations and **prices** (9 packages, 3 per pillar) | Nick (catalogue) |
| Order flow: order request + invoice **or** Stripe checkout | Nick (#112) |
| AI provider: Mistral API **or** Mistral on Vertex AI | Nick (#112) |
| Chat retention, lead retention, lead reply-time promise | Nick |
| **Advisor privacy notice + "AI use on this site" page**: drafts by Brio, reviewed by Ivgil (allow-with-obligation), not approved | DPA, lawyer, Nick |
| Controller legal name and address (the current privacy policy is a generic template) | Nick |
| Disclosure that site texts are drafted with AI help (Brio recommends yes; not legally required for human-reviewed copy) | Nick |
| Accredible credential URL, programme name, portrait photo | Nick |
| FAQ answers, blog teasers, package timelines and FAQs (`[brackets]` in the previews) | Nick / Ivgil |

## 5. Open, not blocking the build

- Logo polish by Tally (optional): tighter gap after the v. If the file is replaced, rebuild and every asset follows.
- Trademark: an attorney should confirm the OFL position; file wordmark + tick as one figurative mark; run a clearance search before filing the symbol. Use ™ now, ® only after registration.

## 6. Rejected or superseded: do not revive

| Date | Item | Replaced by |
|---|---|---|
| 2026-11-10 | Serif headlines (chosen 2026-11-09) | Montserrat |
| 2026-11-10 | Five Montserrat logotype studies (Bookends, Rise, Rise + cell, Cell o, Cell dot): "that is a No" | Tally's tick-v wordmark |
| 2026-11-10 | Tick-v at weights 600, 675, 700 | weight 650 |
| 2026-11-10 | Commercial fonts (TT Norms Pro, Mont, Nexa) | Montserrat (OFL) |
| 2026-11-10 | Wordmark alone in the header | Lockup in the header |
| 2026-11-09 | Logo concepts A, A 1.5×, A 2×, B; concept C revision 1 | Symbol C revision 2 (wider upper chevron) |
| 2026-11-09 | "Site sells services only, not ivGIL" | Widened scope (section 2) |
| 2026-11-09 | Third-party vendor research for the site | Skipped by Nick |

The superseded files are archived (not deleted) in Brio's workspace: `products/ivmanto/archive/design-system-2026-11-10/`.
