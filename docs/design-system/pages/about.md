# About (v1.5)

> **Approved by Nick, 2026-11-13.** Both "Why governance" paragraphs are his own wording: use them verbatim. Any later change to text about Nick needs his yes first.

Preview: `about.html`. Styles: `.about-hero`, `.portrait`, `.profile-links`, `.principles`, `.cred-list` in `components/ivmanto.css`.
Nick, 2026-11-13: the page builds on his 10+ years in data management; AI governance is added to that, it does not replace it. Title: "AI & Data Governance Consultant". h1, lead, meta description and `knowsAbout` name data governance too (Nick approved the wording, 2026-11-13).

## Sections, in order
| # | Section | Ground | Content |
|---|---|---|---|
| 1 | Who | `--bg`, hex grid | Portrait, h1, lead, two buttons, profile links |
| 2 | Why governance | sunken | Heading left, two paragraphs right (Nick's wording) |
| 3 | How I work | `--bg` | Four principles in two columns |
| 4 | Hands-on background | sunken | Heading left, one paragraph and a link right |
| 5 | Credentials | `--bg` | MIT `CredentialBadge`, then a list of certificates |
| 6 | Closing | tinted, hex grid | Same closing as home |

## Copy
**1 · Who**
- Overline: About
- h1: AI and data governance from someone who builds the systems
- Lead: I'm Nikolay Tonev. I help small and mid-sized companies govern their data and their AI: who decides, what can go wrong, and what the law asks for.
- Buttons: "Book a call" (primary, `/booking`), "See the packages" (ghost, `/packages`)
- Profile links (text links, 44px high): LinkedIn, GitHub, Credly. The "DAMA Professionals" link is removed.
- Portrait: the photo already used on home (`/images/nick-portrait.jpg`), 4:5, 360px wide at most, `alt="Nikolay Tonev"`.

**2 · Why governance**
- Overline: Why governance. h2: Rules only work if they fit how the systems are built
- Paragraph 1 (Nick's wording, 2026-11-13, use verbatim): "I've spent more than 20 years in enterprise IT, starting as a business analyst at Wrigley (Mars), and the last ten years designing data platforms and data products in the cloud. At MediaMarktSaturn I built the customer data platform that 100+ marketers relied on, and at Eviden I architected governed data products and generative AI agents on Google Cloud for enterprise clients. Today I bring that experience to SMEs as an independent AI and data governance consultant."
- Paragraph 2 (Nick's wording, 2026-11-13, use verbatim): "At Eviden (Atos Group), I ran a Google Agentspace proof of concept and watched AI agents spread through the enterprise much faster than anyone had planned for. Getting an agent working was the easy part. The harder questions were who controls it, what data it touches, and whether it complies with the rules. That's why I moved into AI governance: to help companies scale AI quickly without losing control of it."
- No third paragraph (Nick, 2026-11-13: closing sentence cut).

**3 · How I work**
- Overline: How I work. h2: Four things you can count on
- One contact: You work with me directly, from the first call to the handover.
- Fixed scope, fixed price: Both are agreed before the work starts.
- Your tools: Registers, policies and reports stay in the tools you already use.
- Not legal advice: Where the law needs interpreting, I work alongside your lawyer.

**4 · Hands-on background**
- Overline: Hands-on background. h2: I still build data platforms and AI agents
- Text: Data platforms and AI agents, built in your own Google Cloud project. That work keeps the governance advice practical.
- Link: "Data & AI engineering" (ghost, `/data-ai`)

**5 · Credentials**
- Overline: Credentials. h2: Training and certificates
- Certificate image (Nick, 2026-11-13): the block shows the certificate Accredible issues, `assets/credentials/mit-pe-certificate.png` (600 x 424, source `https://api.accredible.com/v1/frontend/credential_website_embed_image/certificate/199276375`). Class `.credential--certificate`: image up to 300px wide, left of the text from 640px, above it below; white ground and 1px border in both themes; the image links to the credential; alt text as in `about.html`. **Serve the image from our own site, do not load it from Accredible:** a hotlink sends every visitor's IP address to a third party before any consent (GDPR). Do not crop, recolour or redraw it.
- `CredentialBadge`: "MIT Professional Education certificate", "Applied Agentic AI for Organizational Transformation · 2026", "Verified on Accredible · View credential". Wording rules in `components/credential.md`. Link (Nick, 2026-11-13): `https://www.credential.net/100ff9ab-e82e-41ee-b16e-4ce596b190bb#acc.9ELbmKGO`, the MIT credential itself, new tab. Never add a `?token=` parameter to an Accredible link: that is a personal sign-in token, not part of the public address.
- List (`.cred-list`; name, issuer and year left, one link right):
  - Google Cloud Professional Cloud Architect · Google Cloud · certified 2019, 2022 and 2024 · View on Credly
  - Google Cloud Professional Data Engineer · Google Cloud · certified 2023 · View on Credly
  - Google Cloud Professional Cloud Developer · Google Cloud · certified 2023 · View on Credly
  - Data Architecture Workshop: Practical Exercises in Architecture Design · DATAVERSITY · 2021 · 24 CEUs · Certificate (PDF)
- Line below: "Google Cloud certificates are listed with the years they were earned." and a link "All badges on Credly".
- Credly URL for every Credly link on this page, including the profile link in section 1 (Nick, 2026-11-13): `https://www.credly.com/users/nikolay-tonev.a60be769/badges/credly`. Opens in a new tab.
- Rule (Nick, 2026-11-11): a Google Cloud certificate is never shown as a bare badge, because it reads as current. Always name the years. The badge images are not used on this page; the two skill badges (2024, 2025) move behind the Credly link. `[Clauco: check the exact certificate names and years against Credly before building.]`

**6 · Closing**
- h2: Not sure where to start? Lead: A 30-minute call finds the first step.
- Buttons: "Book a call" (primary), "Ask the advisor" (secondary)

## Removed, and why
| Removed | Why |
|---|---|
| h1 "A Partner in Your Data Journey", the two intro paragraphs | Data-architect positioning; the site now sells AI governance |
| "On the path to becoming a DAMA Certified Data Management Professional" | An intention, not a credential |
| Frank Lloyd Wright quote and the three "Philosophy & Approach" cards | Long, about architecture style; replaced by section 3 |
| "Core Expertise" (three data cards) | Replaced by section 4 with one link |
| Row of five certificate images | They do not load today, and a bare badge reads as current |
| "A testament to my commitment to excellence…" | Says nothing a buyer can check |
| "Featured Insights" with the Data Maturity Assessment card | Insights have their own page; whether the assessment stays is open with Nick |
| Closing "Ready to Unlock Your Data's Potential?" | Replaced by the home closing |

## Page meta (for Clauco; wording to be confirmed by Nick)
- Title: About Nikolay Tonev | ivmanto
- Description: Nikolay Tonev helps small and mid-sized companies govern data and AI: clear owners, known risks, and proof for the EU AI Act and ISO/IEC 42001.
- Person schema: `jobTitle` "AI & Data Governance Consultant" (Nick, 2026-11-13); `knowsAbout`: AI governance, data governance, data management, DAMA-DMBOK, EU AI Act, ISO/IEC 42001, AI risk management, Google Cloud, data architecture (Nick, 2026-11-13).
- The old description claims "20+ years in enterprise data". Do not reuse it: by Nick's own paragraph 1 it is more than 20 years in enterprise IT and the last ten in data platforms. Use durations only as paragraph 1 states them.

## Layout notes
- Section 1: portrait 360px left, text right from 900px; stacked below, portrait first.
- Sections 2 and 4: `1fr / 1.4fr` from 900px, stacked below.
- Section 3: `.principles`, two columns from 768px, each item under a 2px rule in `--text` (the rule of the pillar index).
- No cards on this page except the credential badge.
