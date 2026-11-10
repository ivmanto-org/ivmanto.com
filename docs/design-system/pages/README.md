# Pages

Self-contained HTML previews (open in Claude Desktop or any browser; the header toggle switches the theme). Content in `[brackets]` and `€ —` prices are placeholders.

| File | Page | Purpose |
|---|---|---|
| `home.html` | Home | Funnel. The AI advisor sits in the hero, next to the value proposition. |
| `pillar.html` | Pillar page (Governance; Risk and Compliance use the same template) | Explain the pillar in executive terms, then route to its packages |
| `package.html` | Package detail | Everything needed to decide: scope, timeline, exclusions, fixed price, order |
| `order.html` | Order flow, step 2 | B2B order request with a sticky summary |

## Home: section order and why
1. **Hero + advisor.** Answer "what is this and is it for me" in one line. Ask "what's your situation" right away, so the visitor gets an individual answer instead of reading a menu.
2. **Frameworks strip.** Instant orientation for a CAIO or CDO: the vocabulary they search for.
3. **Three pillars.** For visitors who prefer to browse; each pillar card lists its services.
4. **Packages.** Fixed scope and fixed price are the differentiator against "request a demo" vendors.
5. **How it works.** Lowers perceived risk (four steps, no lock-in).
6. **About Nick.** A person, not a platform: portrait, two to three sentences, Accredible credential.
7. **Insights.** The existing blog, re-skinned.
8. **FAQ.** Objections: legal advice? does the AI Act apply? what happens after ordering? what happens to my chat?
9. **Closing CTA.** Book a call / Ask the advisor.

## Pillar template
Hero (pillar icon, outcome headline, "you likely need this if…" checklist) → what you get (3 deliverables + framework references, **clause numbers checked by Ivgil**) → the pillar's packages (one marked recommended start) → links to the other two pillars.

## Package template
Breadcrumb → tags → title + outcome lead → fact row (duration, your effort, format, company size) → what you get → week-by-week timeline → **not included** (it builds trust and sets scope) → FAQ. Right column (sticky): price, three reassurance ticks, Order (primary), Book a call first (secondary), credential.

## Existing features carried over
Booking calendar (restyle with tokens; target of "Book a call"), blog (PostCard + prose styles: Montserrat headings / Inter body), contact form (same field and consent styles as LeadForm), cookie consent (see `components/cookie-banner.md`).
The ideas generator is replaced by the advisor (#112). Whether the assessment page stays is open with Nick; if it stays, it uses the order stepper as its question progress bar.
