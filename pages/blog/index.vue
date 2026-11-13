<script setup lang="ts">
// Insights index (DS v1.5, pages/blog.md): topic filter, search, 12 per page, state in the URL.
import type { LocationQueryRaw } from 'vue-router'
import {
  clampPage, filterArticles, isTopic, pageCount, pageSlice, pagerItems, topicLabel, topicsInUse,
} from '~/utils/insights'

useSeoMeta({
  title: 'Insights',
  description: 'Articles on AI governance, risk, compliance and the data foundations underneath.',
  ogTitle: 'Insights | ivmanto.com',
  ogDescription: 'Articles on AI governance, risk, compliance and the data foundations underneath.',
  twitterTitle: 'Insights | ivmanto.com',
  twitterDescription: 'Articles on AI governance, risk, compliance and the data foundations underneath.',
})

const { sortedArticles, isLoading, error, fetchArticles } = useArticles()

// Fetch articles at setup time (works during SSG generate and client-side)
await fetchArticles()

const route = useRoute()
const router = useRouter()

// Filter and pager links are plain anchors in a slotted NuxtLink: RouterLink ignores the query when it
// decides what is "current" and would mark every link to /blog with aria-current="page".

// The page is pre-rendered without a query string, so the URL state applies after hydration only.
const isHydrated = ref(false)
onMounted(() => {
  isHydrated.value = true
})

const first = (value: unknown) => (Array.isArray(value) ? value[0] : value)
const topic = computed(() => {
  const value = isHydrated.value ? first(route.query.topic) : ''
  return isTopic(value) ? value : ''
})
const query = computed(() => (isHydrated.value ? String(first(route.query.q) ?? '') : ''))

const topics = computed(() => topicsInUse(sortedArticles.value))
const matches = computed(() => filterArticles(sortedArticles.value, topic.value, query.value))
const pages = computed(() => pageCount(matches.value.length))
const page = computed(() => clampPage(isHydrated.value ? first(route.query.page) : 1, pages.value))
const visible = computed(() => pageSlice(matches.value, page.value))
const pager = computed(() => pagerItems(page.value, pages.value))

function target(next: { topic?: string; page?: number; q?: string }) {
  const state = { topic: topic.value, page: page.value, q: query.value, ...next }
  const out: LocationQueryRaw = {}
  if (state.topic) out.topic = state.topic
  if (state.page > 1) out.page = String(state.page)
  if (state.q.trim()) out.q = state.q
  return { path: '/blog', query: out }
}

// Typing replaces the history entry, so the back button leaves the page instead of undoing letters.
const search = computed({
  get: () => query.value,
  set: (value: string) => {
    router.replace(target({ q: value, page: 1 }))
  },
})

const countLine = computed(() => {
  const n = matches.value.length
  const noun = n === 1 ? 'article' : 'articles'
  if (query.value.trim()) return `${n} ${noun} for “${query.value.trim()}”`
  return pages.value > 1 ? `${n} ${noun} · page ${page.value} of ${pages.value}` : `${n} ${noun}`
})

const emptyLine = computed(() =>
  query.value.trim() ? `No articles match “${query.value.trim()}”.` : `No articles on ${topicLabel(topic.value).toLowerCase()} yet.`,
)

const countEl = ref<HTMLElement | null>(null)
watch(page, () => {
  // A new page starts at its first card, not at the pager the visitor just used.
  nextTick(() => countEl.value?.scrollIntoView({ block: 'start' }))
})
</script>

<template>
  <div>
    <section class="page-head hexgrid">
      <div class="container">
        <div class="section-head">
          <p class="overline">Insights</p>
          <h1>Notes on AI governance</h1>
          <p class="lead">Articles on AI governance, risk, compliance and the data foundations underneath.</p>
        </div>
        <div class="toolbar">
          <!-- Only topics that articles carry are offered; without any, the filter stays away. -->
          <nav v-if="topics.length > 0" aria-label="Filter by topic">
            <ul class="filter">
              <li>
                <NuxtLink v-slot="{ href, navigate }" custom :to="target({ topic: '', page: 1 })"><a class="chip" :href="href" :aria-current="topic === '' ? 'true' : undefined" @click="navigate">All</a></NuxtLink>
              </li>
              <li v-for="t in topics" :key="t">
                <NuxtLink v-slot="{ href, navigate }" custom :to="target({ topic: t, page: 1 })"><a class="chip" :href="href" :aria-current="topic === t ? 'true' : undefined" @click="navigate">{{ topicLabel(t) }}</a></NuxtLink>
              </li>
            </ul>
          </nav>
          <div class="field">
            <label for="blog-search" class="sr-only">Search articles</label>
            <input id="blog-search" v-model="search" type="search" class="input" placeholder="Search articles" />
          </div>
        </div>
      </div>
    </section>

    <section class="section" style="padding-top: 0">
      <div class="container">
        <p v-if="isLoading" class="result-count" role="status">Loading articles…</p>
        <div v-else-if="error" class="notice notice--error" role="alert" style="margin-top: var(--space-6)">
          <p><strong>Articles couldn’t be loaded.</strong> Please try again later.</p>
        </div>
        <template v-else-if="matches.length > 0">
          <p ref="countEl" class="result-count" role="status">{{ countLine }}</p>
          <div class="post-grid">
            <SitePostCard v-for="article in visible" :key="article.slug" :article="article" />
          </div>
          <nav v-if="pages > 1" aria-label="Pages">
            <ul class="pager">
              <li>
                <NuxtLink v-if="page > 1" v-slot="{ href, navigate }" custom :to="target({ page: page - 1 })"><a :href="href" rel="prev" @click="navigate">Previous</a></NuxtLink>
                <span v-else aria-disabled="true">Previous</span>
              </li>
              <li v-for="(item, i) in pager" :key="`${item}-${i}`">
                <span v-if="item === 'gap'" aria-hidden="true">…</span>
                <NuxtLink v-else v-slot="{ href, navigate }" custom :to="target({ page: item })">
                  <a :href="href" :aria-current="item === page ? 'page' : undefined" :aria-label="item === page ? `Page ${item}, current` : `Page ${item}`" @click="navigate">{{ item }}</a>
                </NuxtLink>
              </li>
              <li>
                <NuxtLink v-if="page < pages" v-slot="{ href, navigate }" custom :to="target({ page: page + 1 })"><a :href="href" rel="next" @click="navigate">Next</a></NuxtLink>
                <span v-else aria-disabled="true">Next</span>
              </li>
            </ul>
          </nav>
        </template>
        <p v-else class="result-count" role="status">
          <!-- One line on purpose: Vue drops a line break between two elements, and with it the space. -->
          <strong>{{ emptyLine }}</strong> <NuxtLink v-slot="{ href, navigate }" custom to="/blog"><a :href="href" @click="navigate">Show all articles</a></NuxtLink>
        </p>
      </div>
    </section>
  </div>
</template>

<style scoped>
/* The pager's "…" is not a control: no box around it. */
.pager span[aria-hidden='true'] { border-color: transparent; background: transparent; min-width: 24px; padding: 0; }
.result-count { scroll-margin-top: 96px; }
</style>
