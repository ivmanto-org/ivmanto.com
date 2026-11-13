import type { ArticleMeta } from '~/types/article'

// Insights index and article helpers (DS v1.5, pages/blog.md).

export const TOPICS = ['governance', 'risk', 'compliance', 'engineering'] as const
export type Topic = (typeof TOPICS)[number]

export const ARTICLES_PER_PAGE = 12
export const WORDS_PER_MINUTE = 200
export const RELATED_COUNT = 3

const TOPIC_LABELS: Record<Topic, string> = {
  governance: 'Governance',
  risk: 'Risk',
  compliance: 'Compliance',
  engineering: 'Engineering',
}

const MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']

export function isTopic(value: unknown): value is Topic {
  return typeof value === 'string' && (TOPICS as readonly string[]).includes(value)
}

export function topicLabel(value: unknown): string {
  return isTopic(value) ? TOPIC_LABELS[value] : ''
}

/** "18 Sep 2026", the one date format of the site. Reads the calendar date as written. */
export function formatArticleDate(value: string): string {
  const match = /^(\d{4})-(\d{2})-(\d{2})/.exec(value ?? '')
  if (match) {
    const month = MONTHS[Number(match[2]) - 1]
    if (month) return `${Number(match[3])} ${month} ${match[1]}`
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return `${date.getUTCDate()} ${MONTHS[date.getUTCMonth()]} ${date.getUTCFullYear()}`
}

/** Reading time from rendered HTML or plain text: 200 words a minute, rounded up, never below 1. */
export function readingMinutes(content: string): number {
  const text = (content ?? '').replace(/<[^>]*>/g, ' ')
  const words = text.split(/\s+/).filter(Boolean).length
  return Math.max(1, Math.ceil(words / WORDS_PER_MINUTE))
}

/** Minutes for a card: the value the API sends, if it sends a usable one. */
export function cardReadingMinutes(article: ArticleMeta): number | null {
  const minutes = article.reading_minutes
  return typeof minutes === 'number' && Number.isFinite(minutes) && minutes >= 1 ? Math.ceil(minutes) : null
}

export function filterArticles(articles: readonly ArticleMeta[], topic: string, query: string): ArticleMeta[] {
  const needle = query.trim().toLowerCase()
  return articles.filter((article) => {
    if (isTopic(topic) && article.topic !== topic) return false
    if (!needle) return true
    return article.title.toLowerCase().includes(needle) || article.summary.toLowerCase().includes(needle)
  })
}

/** Topics that at least one article carries, in the fixed order of the filter. */
export function topicsInUse(articles: readonly ArticleMeta[]): Topic[] {
  return TOPICS.filter((topic) => articles.some((article) => article.topic === topic))
}

export function pageCount(total: number, perPage = ARTICLES_PER_PAGE): number {
  return Math.max(1, Math.ceil(total / perPage))
}

/** A page number from the URL, held inside 1..pages. */
export function clampPage(value: unknown, pages: number): number {
  const page = Number.parseInt(String(value ?? ''), 10)
  if (!Number.isFinite(page) || page < 1) return 1
  return Math.min(page, pages)
}

export function pageSlice<T>(items: readonly T[], page: number, perPage = ARTICLES_PER_PAGE): T[] {
  return items.slice((page - 1) * perPage, page * perPage)
}

export type PagerItem = number | 'gap'

/** Up to 7 pages: all of them. Above: first, last and two around the current one, with gaps. */
export function pagerItems(current: number, pages: number): PagerItem[] {
  if (pages <= 7) return Array.from({ length: pages }, (_, i) => i + 1)
  const shown = new Set([1, pages])
  for (let p = current - 2; p <= current + 2; p++) if (p >= 1 && p <= pages) shown.add(p)
  const sorted = [...shown].sort((a, b) => a - b)
  return sorted.flatMap((p, i) => (i > 0 && p - sorted[i - 1]! > 1 ? ['gap' as const, p] : [p]))
}

/** Three articles of the same topic, newest first, filled up from all the others. Expects newest first. */
export function relatedArticles(sorted: readonly ArticleMeta[], current: ArticleMeta, count = RELATED_COUNT): ArticleMeta[] {
  const others = sorted.filter((article) => article.slug !== current.slug)
  const sameTopic = isTopic(current.topic) ? others.filter((article) => article.topic === current.topic) : []
  const rest = others.filter((article) => !sameTopic.includes(article))
  return [...sameTopic, ...rest].slice(0, count)
}

/** Marks a closing "Sources" paragraph so it gets the small muted style (.prose .sources). */
export function markSources(html: string): string {
  return html.replace(/<p>(\s*<strong>Sources:?<\/strong>)/gi, '<p class="sources">$1')
}

/** JSON for a <script type="application/ld+json"> block; "<" cannot close the script early. */
export function jsonLd(data: unknown): string {
  return JSON.stringify(data).replace(/</g, '\\u003c')
}
