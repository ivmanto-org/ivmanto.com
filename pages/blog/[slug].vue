<template>
  <div class="bg-light-gray py-12 sm:py-16" @click="closeDropdown">
    <!-- Loading State -->
    <div v-if="isLoadingContent" class="text-center py-16">
      <div
        class="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-current border-r-transparent align-[-0.125em] motion-reduce:animate-[spin_1.5s_linear_infinite]"
        role="status"
      >
        <span class="!absolute !-m-px !h-px !w-px !overflow-hidden !whitespace-nowrap !border-0 !p-0 ![clip:rect(0,0,0,0)]"
          >Loading...</span
        >
      </div>
    </div>

    <!-- Article Content -->
    <div v-else-if="article" class="mx-auto max-w-4xl px-4 sm:px-6 lg:px-8">
      <!-- Top Navigation Bar -->
      <div
        class="mb-8 p-2 bg-surface rounded-lg shadow-sm flex items-center justify-between gap-2 text-sm font-medium"
      >
        <!-- Previous Article Button -->
        <NuxtLink
          v-if="previousArticle"
          :to="`/blog/${previousArticle.slug}`"
          class="nav-button"
          :title="`Previous: ${previousArticle.title}`"
        >
          <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
          </svg>
        </NuxtLink>
        <span v-else class="nav-button-disabled">
          <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
          </svg>
        </span>

        <!-- Custom Dropdown -->
        <div class="relative flex-grow flex justify-center">
          <button @click.stop="toggleDropdown" class="nav-button-dropdown">
            <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h6a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
            </svg>
            <span class="hidden sm:inline ml-2">Articles</span>
          </button>
          <transition
            enter-active-class="transition ease-out duration-100"
            enter-from-class="transform opacity-0 scale-95"
            enter-to-class="transform opacity-100 scale-100"
            leave-active-class="transition ease-in duration-75"
            leave-from-class="transform opacity-100 scale-100"
            leave-to-class="transform opacity-0 scale-95"
          >
            <div v-if="isDropdownOpen" class="absolute top-full mt-2 w-72 max-h-80 overflow-y-auto bg-surface rounded-lg shadow-xl ring-1 ring-black ring-opacity-5 z-50">
              <div class="py-1">
                <button v-for="a in otherArticles" :key="a.slug" @click="navigateToArticle(a.slug)" class="dropdown-item">
                  {{ a.title }}
                </button>
              </div>
            </div>
          </transition>
        </div>

        <!-- Next Article Button -->
        <NuxtLink
          v-if="nextArticle"
          :to="`/blog/${nextArticle.slug}`"
          class="nav-button"
          :title="`Next: ${nextArticle.title}`"
        >
          <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
          </svg>
        </NuxtLink>
        <span v-else class="nav-button-disabled">
          <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
          </svg>
        </span>

        <!-- Share Button -->
        <div class="border-l border-border pl-2">
          <button @click="shareArticle" class="nav-button" title="Share article">
            <svg v-if="!copied" class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path d="M15 8a3 3 0 10-2.977-2.63l-4.94 2.47a3 3 0 100 4.319l4.94 2.47a3 3 0 10.895-1.789l-4.94-2.47a3.027 3.027 0 000-.74l4.94-2.47C13.456 7.68 14.19 8 15 8z" />
            </svg>
            <svg v-else class="h-5 w-5 text-green-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          </button>
        </div>
      </div>

      <!-- Article Content -->
      <div class="bg-surface rounded-lg shadow-lg p-8 md:p-12">
        <div class="mb-8 border-b pb-8 border-border">
          <p class="text-base text-muted">Published on {{ formatDate(article.date) }}</p>
          <h1 class="mt-2 text-3xl font-bold tracking-tight text-text sm:text-4xl">
            {{ article.title }}
          </h1>
          <p class="mt-6 text-xl leading-8 text-muted">{{ article.summary }}</p>
        </div>
        <div class="prose prose-lg max-w-none" v-html="sanitizedContent"></div>

        <!-- Like button section -->
        <div class="mt-12 pt-8 border-t border-border flex justify-end items-center">
          <button
            @click="handleLike"
            :disabled="isLoadingLikes"
            class="flex items-center gap-2 text-muted hover:text-primary disabled:text-primary disabled:cursor-not-allowed transition-colors"
            title="Like this article"
          >
            <svg
              class="h-6 w-6"
              :class="{ 'text-red-500 fill-current': isLiked }"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4.318 6.318a4.5 4.5 0 016.364 0L12 7.5l1.318-1.182a4.5 4.5 0 116.364 6.364L12 20.364l-7.682-7.682a4.5 4.5 0 010-6.364z"
              />
            </svg>
            <span v-if="!isLoadingLikes" class="font-semibold text-lg tabular-nums">{{ likeCount }}</span>
            <span v-else class="font-semibold text-lg">...</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Article Not Found State -->
    <div v-else class="text-center py-16">
      <h1 class="text-3xl font-bold">Article not found</h1>
      <p class="mt-4">Sorry, we couldn't find the article you're looking for.</p>
      <NuxtLink to="/blog" class="mt-6 inline-block text-primary hover:underline">
        &larr; Back to all articles
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { trackEvent } from '~/services/analytics'
import DOMPurify from 'dompurify'
import type { Article } from '~/types/article'

const route = useRoute()
const slug = computed(() => route.params.slug as string)

const { sortedArticles, fetchArticles, fetchArticle } = useArticles()

const isDropdownOpen = ref(false)
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

// Head Management for SEO (Fix Soft 404)
useHead({
  meta: computed(() => {
    if (!isLoadingContent.value && !article.value) {
      return [{ name: 'robots', content: 'noindex' }]
    }
    return []
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
          innerHTML: JSON.stringify({
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
      ]
    }
    return []
  }),
})

const sanitizedContent = computed(() => {
  if (!article.value?.content) return ''
  // Strip the leading <h1> from the content to avoid duplicating the title
  // which is already rendered from article.title in the template above.
  const contentWithoutLeadingH1 = article.value.content.replace(/^\s*<h1[^>]*>.*?<\/h1>\s*/i, '')
  // DOMPurify only works client-side; during SSR, return raw content
  if (import.meta.server) return contentWithoutLeadingH1
  return DOMPurify.sanitize(contentWithoutLeadingH1)
})

const currentArticleIndex = computed(() => {
  if (!article.value) return -1
  return sortedArticles.value.findIndex((a) => a.slug === article.value?.slug)
})

const previousArticle = computed(() => {
  if (currentArticleIndex.value > 0) return sortedArticles.value[currentArticleIndex.value - 1]
  return undefined
})

const nextArticle = computed(() => {
  if (currentArticleIndex.value < sortedArticles.value.length - 1 && currentArticleIndex.value !== -1)
    return sortedArticles.value[currentArticleIndex.value + 1]
  return undefined
})

const otherArticles = computed(() => {
  if (!article.value) return []
  return sortedArticles.value.filter((a) => a.slug !== article.value?.slug)
})

function navigateToArticle(targetSlug: string) {
  navigateTo(`/blog/${targetSlug}`)
}

function toggleDropdown() {
  isDropdownOpen.value = !isDropdownOpen.value
}

function closeDropdown() {
  isDropdownOpen.value = false
}

async function shareArticle() {
  if (!article.value) return
  const shareData = {
    title: article.value.title,
    text: article.value.summary,
    url: window.location.href,
  }
  if (navigator.share) {
    try { await navigator.share(shareData) } catch (err) { console.error('Error sharing:', err) }
  } else {
    await navigator.clipboard.writeText(window.location.href)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  }
}

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

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })
}

// Client-only: fetch likes and check localStorage on mount
onMounted(() => {
  fetchLikes()
  if (localStorage.getItem(`liked-${slug.value}`)) isLiked.value = true
})

// When slug changes client-side (navigation), reset likes state
watch(slug, (newSlug) => {
  closeDropdown()
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

<style lang="postcss" scoped>
.nav-button {
  @apply h-10 w-10 flex items-center justify-center bg-surface-sunken hover:bg-border text-muted hover:text-primary rounded-md transition-all duration-200 hover:scale-110;
}
.nav-button-dropdown {
  @apply h-10 px-4 flex items-center justify-center bg-surface-sunken hover:bg-border text-muted hover:text-primary rounded-md transition-all duration-200;
}
.nav-button-disabled {
  @apply h-10 w-10 flex items-center justify-center bg-surface-sunken text-muted rounded-md cursor-not-allowed opacity-75;
}
.dropdown-item {
  @apply block w-full text-left px-4 py-2 text-sm text-muted hover:bg-surface-sunken hover:text-text;
}
.tabular-nums {
  font-variant-numeric: tabular-nums;
}
</style>
