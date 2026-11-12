# Components

Preview: `../components.html` (both themes). Styles: `ivmanto.css` (reference implementation).

## Rules for implementation (Clauco)
- One Vue component per spec. Keep the class and prop names in the specs so design and code use the same words.
- Use semantic tokens only (`var(--accent)`, Tailwind `bg-surface`), never hex values or primitives.
- Every interactive element: 44px minimum touch target, visible `:focus-visible` ring (`--focus`), keyboard operable.
- Motion: `--dur-fast` / `--dur-base` only, no infinite animation except the typing indicator and the home hero background (`hero-motion.md`, with its pause button), and everything respects `prefers-reduced-motion`. Text never moves after it has appeared.
- Copy in the specs is final unless it's marked `[placeholder]`. Package names, scopes and prices are placeholders until Nick defines the catalogue.

| Spec | Component(s) |
|---|---|
| [button.md](button.md) | `Btn` (primary, secondary, ghost, sm, block) |
| [tag-badge.md](tag-badge.md) | `Tag` (risk, framework, pillar), `BadgeRecommended` |
| [cards.md](cards.md) | `PillarCard`, `PackageCard` (full, compact), `PostCard` |
| [backgrounds.md](backgrounds.md) | Home: dark band (4b + 5), tinted closing CTA |
| [hero-motion.md](hero-motion.md) | Home hero: glow, hex-grid wave, entrance, pause button |
| [advisor.md](advisor.md) | `AdvisorChat`, `AdvisorMessage`, `RecommendationCard`, `CtaSet`, `LeadForm` |
| [order-stepper.md](order-stepper.md) | `OrderStepper`, order summary |
| [credential.md](credential.md) | `CredentialBadge`, `FrameworkStrip` |
| [layout.md](layout.md) | `SiteHeader`, `SiteFooter`, `Section`, `HowItWorks`, `Faq`, hex-grid texture |
| [cookie-banner.md](cookie-banner.md) | `CookieBanner` |
