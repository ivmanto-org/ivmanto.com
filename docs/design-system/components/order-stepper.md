# Order flow (`OrderStepper`, summary)

Steps: 1 Package → 2 Your details → 3 Review → 4 Confirmation. Preview: `../pages/order.html`.

- Stepper: numbered square markers (radius 6px, echoing the hex grid without drawing hexagons), 2px underline. Current step: `--accent` fill plus `aria-current="step"`. Done: check mark. On mobile only the current step's label is shown.
- Layout: form left, sticky order summary right (stacked under the form on mobile).
- Fields: name, role, work email, phone (optional), company, VAT ID (B2B), billing address, company size, context (optional).
- One checkbox: terms + privacy + "ordering on behalf of a business" (B2B only, so no consumer withdrawal flow).
- Payment wording depends on Nick's decision (#112): **order request → invoice** (current preview) or a Stripe checkout at step 3. The design works for both; at step 3 the summary card hosts the Stripe element.
- Confirmation page: what happens next (3 numbered items), a calendar link for kickoff, and a copy of the order by email.
