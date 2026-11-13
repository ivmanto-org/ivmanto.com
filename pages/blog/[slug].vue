<template>
  <div>
    <!-- DS v1.5 (pages/blog.md): one text column on the page ground, no card, no top toolbar. -->
    <section class="section" style="padding-top: var(--space-8)">
      <div class="container">
        <p v-if="isLoadingContent" class="muted" role="status">Loading the article…</p>

        <article v-else-if="article" class="article">
          <NuxtLink to="/blog" class="back-link"><SiteChevIcon /> All insights</NuxtLink>
          <header class="article__head">
            <p class="post">
              <span class="meta"><template v-if="topic"><b>{{ topic }}</b> · </template>{{ formatArticleDate(article.date) }} · {{ minutes }} min read</span>
            </p>
            <h1>{{ article.title }}</h1>
            <p class="lead">{{ article.summary }}</p>
            <p class="byline">
              <img src="/images/nick-portrait.jpg" alt="" width="40" height="40" loading="lazy" />
              <span>By <NuxtLink to="/about">Nikolay Tonev</NuxtLink> · AI &amp; Data Governance Consultant</span>
            </p>
          </header>

          <div class="prose max-w-none" v-html="sanitizedContent"></div>

          <footer class="article__foot">
            <button class="like" type="button" :aria-pressed="isLiked" :disabled="isLoadingLikes" @click="handleLike">
              <svg viewBox="0 0 24 24" :fill="isLiked ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M12 20.5l-7.7-7.7a4.5 4.5 0 0 1 6.4-6.4L12 7.7l1.3-1.3a4.5 4.5 0 0 1 6.4 6.4z" /></svg>
              <span>Helpful</span>
              <span v-if="!isLoadingLikes" class="like__count" :aria-label="`${likeCount} ${likeCount === 1 ? 'reader' : 'readers'} found this helpful`">{{ likeCount }}</span>
            </button>
            <button class="btn btn--ghost" type="button" @click="copyLink">
              <span aria-live="polite">{{ copied ? 'Link copied' : 'Copy link' }}</span>
            </button>
          </footer>

          <aside class="next-step" aria-labelledby="next-h">
            <h2 id="next-h">Does this apply to your company?</h2>
            <p>A 30-minute call finds the first step.</p>
            <div class="row">
              <NuxtLink to="/booking" class="btn btn--primary">Book a call</NuxtLink>
              <NuxtLink to="/#main" class="btn btn--ghost">Ask the advisor <SiteChevIcon /></NuxtLink>
            </div>
          </aside>

          <nav v-if="newerArticle || olderArticle" class="post-nav" aria-label="More articles">
            <NuxtLink v-if="newerArticle" :to="`/blog/${newerArticle.slug}`"><small>Newer</small>{{ newerArticle.title }}</NuxtLink>
            <span v-else aria-hidden="true"></span>
            <NuxtLink v-if="olderArticle" :to="`/blog/${olderArticle.slug}`"><small>Older</small>{{ olderArticle.title }}</NuxtLink>
          </nav>
        </article>

        <div v-else class="article">
          <h1>Article not found</h1>
          <p class="lead">This article may have moved.</p>
          <NuxtLink to="/blog" class="btn btn--primary">All insights</NuxtLink>
        </div>
      </div>
    </section>

    <section v-if="article && related.length > 0" class="section section--sunken">
      <div class="container">
        <div class="section-head">
          <p class="overline">Keep reading</p>
          <h2>{{ topic ? `More on ${topic.toLowerCase()}` : 'More insights' }}</h2>
        </div>
        <div class="post-grid">
          <SitePostCard v-for="a in related" :key="a.slug" :article="a" />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { trackEvent } from '~/services/analytics'
import DOMPurify from 'dompurify'
import type { Article } from '~/types/article'
import { articleProvenanceHead } from '~/utils/articleProvenance'
import { cardReadingMinutes, formatArticleDate, jsonLd, markSources, readingMinutes, relatedArticles, topicLabel } from '~/utils/insights'

const route = useRoute()
const slug = computed(() => route.params.slug as string)

const { sortedArticles, fetchArticles, fetchArticle } = useArticles()

const copied = ref(false)
const likeCount = ref(0)
const isLiked = ref(false)
const isLoadingLikes = ref(true)

// Fetch articles list for navigation
await fetchArticles()

// Fetch the current article (works during SSG generate and client-side)
const { data: article, status: articleStatus } = await useAsyncData(
  `article-${slug.value}`,
  () => fetchArticle(slug.value),
  { watch: [slug] },
)

const isLoadingContent = computed(() => articleStatus.value === 'pending')

// Truncate article title so the full "<title> | ivmanto.com" stays ≤ 70 chars
function seoTitle(raw: string | undefined): string {
  if (!raw) return 'Blog | ivmanto.com'
  const suffix = ' | ivmanto.com' // 14 chars
  const maxTitleLen = 70 - suffix.length // 56 chars for the article title
  const truncated = raw.length > maxTitleLen ? raw.slice(0, maxTitleLen - 1).trimEnd() + '…' : raw
  return `${truncated}${suffix}`
}

// Dynamic OG/Twitter meta for shared blog posts
useSeoMeta({
  title: computed(() => seoTitle(article.value?.title)),
  description: computed(() => article.value?.summary ?? ''),
  ogTitle: computed(() => article.value?.title ?? 'Blog | ivmanto.com'),
  ogDescription: computed(() => article.value?.summary ?? ''),
  ogType: 'article',
  ogUrl: computed(() => `https://ivmanto.com/blog/${route.params.slug}`),
  twitterTitle: computed(() => article.value?.title ?? 'Blog | ivmanto.com'),
  twitterDescription: computed(() => article.value?.summary ?? ''),
})

// AI provenance marking (IPTC DigitalSourceType), empty for unmarked and human articles
const provenanceHead = computed(() => articleProvenanceHead(article.value))

// Head Management for SEO (Fix Soft 404)
useHead({
  meta: computed(() => {
    if (!isLoadingContent.value && !article.value) {
      return [{ name: 'robots', content: 'noindex' }]
    }
    return provenanceHead.value.meta
  }),
  script: computed(() => {
    if (article.value) {
      // Google requires headline ≤ 110 chars for Article rich results
      const headline = article.value.title.length > 110
        ? article.value.title.slice(0, 109) + '…'
        : article.value.title
      return [
        {
          type: 'application/ld+json',
          innerHTML: jsonLd({
            '@context': 'https://schema.org',
            '@type': 'BlogPosting',
            mainEntityOfPage: {
              '@type': 'WebPage',
              '@id': `https://ivmanto.com/blog/${article.value.slug}`,
            },
            headline,
            description: article.value.summary,
            image: 'https://ivmanto.com/social-sharing-card.webp',
            datePublished: article.value.date,
            dateModified: article.value.date,
            author: { '@id': 'https://ivmanto.com/about#person' },
            publisher: { '@id': 'https://ivmanto.com/#organization' },
            url: `https://ivmanto.com/blog/${article.value.slug}`,
          }),
        },
        ...provenanceHead.value.script,
      ]
    }
    return []
  }),
})

const sanitizedContent = computed(() => {
  if (!article.value?.content) return ''
  // Strip the leading <h1> from the content to avoid duplicating the title
  // which is already rendered from article.title in the template above.
  const contentWithoutLeadingH1 = markSources(article.value.content.replace(/^\s*<h1[^>]*>.*?<\/h1>\s*/i, ''))
  // DOMPurify only works client-side; during SSR, return raw content
  if (import.meta.server) return contentWithoutLeadingH1
  return DOMPurify.sanitize(contentWithoutLeadingH1)
})

const topic = computed(() => topicLabel(article.value?.topic))
const minutes = computed(() => {
  if (!article.value) return 1
  return cardReadingMinutes(article.value) ?? readingMinutes(article.value.content)
})
const related = computed(() => (article.value ? relatedArticles(sortedArticles.value, article.value) : []))

const currentArticleIndex = computed(() => {
  if (!article.value) return -1
  return sortedArticles.value.findIndex((a) => a.slug === article.value?.slug)
})

// sortedArticles is newest first: the one before is newer, the one after is older.
const newerArticle = computed(() => (currentArticleIndex.value > 0 ? sortedArticles.value[currentArticleIndex.value - 1] : undefined))
const olderArticle = computed(() =>
  currentArticleIndex.value !== -1 && currentArticleIndex.value < sortedArticles.value.length - 1
    ? sortedArticles.value[currentArticleIndex.value + 1]
    : undefined,
)

const COPIED_LABEL_MS = 2000
let copiedTimer: ReturnType<typeof setTimeout> | undefined

async function copyLink() {
  try {
    await navigator.clipboard.writeText(window.location.href)
  } catch {
    // Clipboard blocked (permissions, insecure origin): the label stays, the address bar still has the link.
    return
  }
  copied.value = true
  clearTimeout(copiedTimer)
  copiedTimer = setTimeout(() => { copied.value = false }, COPIED_LABEL_MS)
}

onBeforeUnmount(() => clearTimeout(copiedTimer))

async function fetchLikes() {
  if (!article.value) return
  isLoadingLikes.value = true
  try {
    const response = await fetch(`/api/articles/${article.value.slug}/likes`)
    if (response.ok) {
      const data = await response.json()
      likeCount.value = data.likes || 0
    } else { likeCount.value = 0 }
  } catch (error) { likeCount.value = 0 } finally { isLoadingLikes.value = false }
}

async function handleLike() {
  if (!article.value) return
  const wasLiked = isLiked.value
  if (wasLiked) {
    isLiked.value = false
    likeCount.value--
    localStorage.removeItem(`liked-${article.value.slug}`)
    trackEvent('like_insight', { insight_id: article.value.slug, insight_title: article.value.title, like_status: 'unliked' })
  } else {
    isLiked.value = true
    likeCount.value++
    localStorage.setItem(`liked-${article.value.slug}`, 'true')
    trackEvent('like_insight', { insight_id: article.value.slug, insight_title: article.value.title, like_status: 'liked' })
  }
  try {
    const method = wasLiked ? 'DELETE' : 'POST'
    const response = await fetch(`/api/articles/${article.value.slug}/like`, { method })
    if (response.ok) { const data = await response.json(); likeCount.value = data.likes }
  } catch (error) { console.error('Network error while syncing like status:', error) }
}

// Client-only: fetch likes and check localStorage on mount
onMounted(() => {
  fetchLikes()
  if (localStorage.getItem(`liked-${slug.value}`)) isLiked.value = true
})

// When slug changes client-side (navigation), reset likes state
watch(slug, (newSlug) => {
  copied.value = false
  likeCount.value = 0
  isLiked.value = false
  isLoadingLikes.value = true
  // article data is refetched automatically by useAsyncData's watch option
  nextTick(() => {
    if (localStorage.getItem(`liked-${newSlug}`)) isLiked.value = true
    fetchLikes()
  })
})
</script>

<style scoped>
.like__count { font-variant-numeric: tabular-nums; }
.like[disabled] { cursor: default; }
.byline img { object-position: 50% 20%; }
/* An empty first cell keeps "Older" in the right column when there is no newer article. */
.post-nav > span { display: none; }
@media (min-width: 640px) { .post-nav > span { display: block; } }
</style>
