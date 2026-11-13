# Insights: index and article (v1.5)

Previews: `blog.html` (index), `post.html` (article). Styles: `.post`, `.post-grid`, `.toolbar`, `.filter`, `.pager`, `.article`, `.prose`, `.next-step`, `.post-nav` in `components/ivmanto.css`.
Routes stay as they are: `/blog` and `/blog/<slug>`. The nav label is "Insights".

## What changes for the reader
| Today | v1.5 |
|---|---|
| About 50 articles on one page | 12 per page, with a pager |
| No way to narrow the list | Topic filter: All, Governance, Risk, Compliance, Engineering |
| First row three columns, the rest two | One grid: 3 columns from 1000px, 2 from 640px, 1 below |
| Article sits in a floating card with a toolbar on top | Article is one text column on the page ground |
| No author, no next step | Byline, a "next step" block, newer and older links, three related articles |

## Data (Clauco, Ivmo)
- New field on every article: `topic`, one of `governance`, `risk`, `compliance`, `engineering`. One topic per article. An article without a topic shows under "All" only and carries no topic in its meta line.
- Reading time is computed from the word count at 200 words a minute, rounded up, never below 1.
- `summary`: 40 words at most, one or two sentences. Longer summaries are cut by the 2-line clamp on cards, and need editing for the article page. Several current summaries run to 60 words and more.
- Dates: day, short month, year ("18 Sep 2026"), the same format everywhere on the site.

## Index (`/blog`)
1. **Page head** (`.page-head.hexgrid`): overline "Insights", h1 "Notes on AI governance", lead "Articles on AI governance, risk, compliance and the data foundations underneath."
2. **Toolbar** (`.toolbar`): topic filter left, search right (stacked below 640px, filter first).
   - Filter: links styled as `.chip`, 44px high, the active one filled `--accent` with `aria-current="true"`. In a `<nav aria-label="Filter by topic">`.
   - Search: the existing field, label "Search articles" (screen-reader only), placeholder "Search articles". It searches title and summary, inside the chosen topic.
   - URL state: `?topic=risk&page=2&q=…`, so a filtered page can be shared and the back button works. Changing the topic or the search goes back to page 1.
3. **Count line** (`.result-count`, `role="status"`): "52 articles · page 1 of 5". With a search: "3 articles for “register”".
4. **Grid** (`.post-grid`) of `PostCard`, 12 per page, newest first.
5. **Pager** (`.pager`): Previous, page numbers, Next. Current page filled, `aria-current="page"`. Previous or Next is greyed out (not a link) at the ends. Above 7 pages show first, last and two around the current one, with "…". Hidden when there is one page.
6. **States:** loading "Loading articles…"; error `.notice--error` "Articles couldn't be loaded. Please try again later."; nothing found: "No articles match “<query>”." plus a link "Show all articles" that clears search and filter.

### PostCard (v1.5)
`article.card.card--link.post`: meta line, title (h3, 18px Montserrat, the only link, covers the card), summary clamped to two lines.
Meta line: `TOPIC · 18 SEP 2026 · 6 MIN` in JetBrains Mono 12px uppercase, topic in `--accent`, the rest muted. On the home page the card keeps date and title only (v1.1).

## Article (`/blog/<slug>`)
One column, 720px at most, centred, on `--bg`. No card around the text.
1. **Back link** "‹ All insights" to `/blog`.
2. **Head** (`.article__head`): meta line (as on the card, with "6 min read"), h1 (`clamp(2rem, 4vw, 2.75rem)`), summary as `.lead`, byline.
3. **Byline** (`.byline`): portrait 40px round, "By Nikolay Tonev · AI & Data Governance Consultant", the name links to `/about`. Job title set by Nick, 2026-11-13.
4. **Body** (`.prose`): Inter 17px, line height 1.7. h2 24px, h3 20px, both Montserrat, sentence case. Links `--accent`, underlined. List markers `--accent`. Quotes: 3px `--accent` rule on the left, 18px, muted. Code: JetBrains Mono on `--surface-sunken`. Tables scroll sideways inside the column on phones; long URLs wrap. The site keeps Tailwind typography; `.prose` in `ivmanto.css` lists the values it must produce.
5. **Sources**: a last paragraph `.sources`, 14px muted, after a rule.
6. **Foot** (`.article__foot`): the existing like button as "Helpful" with the count (`.like`, `aria-pressed`), and "Copy link" (ghost; the label changes to "Link copied" for 2 seconds).
7. **Next step** (`.next-step`, tinted): h2 "Does this apply to your company?", "A 30-minute call finds the first step.", buttons "Book a call" (primary, to `/booking`) and "Ask the advisor" (ghost, to the home advisor).
8. **Newer / Older** (`.post-nav`): two links with the article titles. They replace the arrow buttons and the article dropdown of the old top toolbar.
9. **Keep reading** (sunken section): three `PostCard`s of the same topic, newest first, without the current one; fill up from all articles if the topic has fewer. Heading "More on <topic>"; "More insights" when the article has no topic.
10. **Not found:** h1 "Article not found", "This article may have moved.", button "All insights".

## Removed
The top toolbar on the article (previous, dropdown of all articles, next, share icon), the card and shadow around the article, the centred loading spinner (use "Loading the article…").

## Checks before it ships
One h1 per page. Every chip, pager item and button 44px high. Keyboard: filter, search, cards, pager in that order, visible focus ring. No horizontal scroll at 390px with a long title, a table and a long URL in the text.

## Outside this spec
Many current articles are AI news commentary and sit far from the governance offer. Which stay, which get which topic, and what gets written next is a content decision for Nick (with Ivmo); the design works either way.
