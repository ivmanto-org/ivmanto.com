# AI service advisor (`AdvisorChat`)

The home-page funnel (epic #112). See `../components.html` for all states and `../pages/home.html` for placement.

## Placement
- Desktop: right column of the hero, next to the headline. Mobile: directly below the headline and lead, **before** the secondary buttons, so it's visible without scrolling on most phones.
- Also reachable from "Ask the advisor" links (pillar pages, closing CTA). These scroll to the home hero, or open the same component in a drawer on other pages.

## Anatomy
1. **Header:** AI avatar (white tick on a teal disc), title ("Which service do you need?"), `AI` label. The label is mandatory.
2. **Log:** AI messages left with avatar; user messages right, on `--surface-sunken`.
3. **Prompt chips** (idle only): 3 example starters. Clicking one fills the textarea; it does not auto-send.
4. **Form:** textarea (max 600 chars, counter), send button (primary, icon, `aria-label="Send"`).
5. **Meta line**, always visible: "AI assistant · guidance, not legal advice · how we handle your message" (link).

## States
| # | State | Behaviour |
|---|---|---|
| 1 | Idle, first visit | Greeting plus chips. **Privacy notice is shown above the input before the first send**; it collapses after the first send. |
| 2 | Thinking | User message appears at once; typing indicator (`role="status"`, label "Advisor is answering"). Send is disabled. |
| 3 | Answer | Text reply (plain paragraphs, max ~120 words), then "Services that fit" with 1–3 `RecommendationCard`s, then `CtaSet`. |
| 4 | Lead form | Opens inline below the answer when "Leave contact" is chosen. |
| 5 | Error / timeout (> 20 s) | Notice: "The advisor can't answer right now. Your question wasn't lost." Keep the user text in the textarea; offer Book a call / Browse packages. |
| 6 | Limit reached (rate limit or 4 turns) | Notice steering to a person: Book a call (primary), Leave contact. |

## RecommendationCard
- Rendered **from the catalogue by `id`**. Name, price and link come from catalogue data, never from model text. Only `why` (one sentence, ≤ 160 chars) comes from the model and is shown muted.
- Layout: pillar icon · title · price (Montserrat) / why underneath. The whole card links to the package page.
- An unknown `id` from the model is dropped silently and logged server-side.

## CtaSet
Always four options, in this order; `next_action` from the model decides which one is **primary**:
| Key | Label | Sub-label | Target |
|---|---|---|---|
| `buy` | Order package | Fixed scope and price | package page → order flow |
| `call` | Book a call | [30 min, free] | booking calendar |
| `schedule` | Schedule later | Pick a time | booking calendar (date picker first) |
| `contact` | Leave contact | [Reply within 1 working day] | inline LeadForm |

## LeadForm
Fields: name, work email, company (required); message prefilled with the chat summary (editable).
- Consent checkbox (required, **unticked**): "I agree that ivmanto may use my details and my advisor conversation to answer this request. Withdraw any time. Privacy notice". Legal basis: Art. 6(1)(b) for answering the request; the checkbox documents the conversation-sharing consent, Art. 6(1)(a).
- Optional second checkbox for insights by email (separate consent, unticked; double opt-in required).
- Success: "Thanks. Nick will reply [within one working day]." No auto-redirect.

## Legal and trust (non-negotiable)
- EU AI Act Art. 50: the visitor must know it's an AI. Hence the `AI` label in the header and the meta line, every state.
- "Not legal advice" is visible next to the input, not only in the terms.
- Privacy notice before the first send: processor, purpose, retention, "don't include personal or confidential data". Draft text: Brio's `advisor-privacy-and-ai-use-DRAFT` (Ivgil reviews, Nick approves; not final). The meta-line link goes to `/ai#advisor`; the privacy section is `/privacy-policy#advisor`.
- Backend facts the notice relies on (Ivmo, #123/#124): no chat text stored; random `session_id`; logs carry no message content (the request logger does record IP, URL and status); rate-limit counters kept 24 h; leads go by email to Nick only, with no hidden fields.
- Funnel analytics events only after cookie consent (existing GTM gating).

## Accessibility
- Log container: `aria-live="polite"`. New AI messages are announced; the typing indicator is announced once.
- Enter sends, Shift+Enter adds a new line. Focus returns to the textarea after the answer arrives.
- Chips are `<button>`s. Recommendation cards and CTAs are links with full-text names.
