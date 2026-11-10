<script setup lang="ts">
useSeoMeta({
  title: 'Insights',
  description:
    'Articles on AI governance, risk, compliance and the data foundations underneath.',
  ogTitle: 'Insights & Articles | ivmanto.com',
  ogDescription:
    'Expert articles on data strategy, cloud architecture, AI/ML, and software engineering.',
  twitterTitle: 'Insights & Articles | ivmanto.com',
  twitterDescription:
    'Expert articles on data strategy, cloud architecture, AI/ML, and software engineering.',
})

const { sortedArticles, isLoading, error, fetchArticles } = useArticles()

// Fetch articles at setup time (works during SSG generate and client-side)
await fetchArticles()

const searchQuery = ref('')

const filteredArticles = computed(() => {
  if (!searchQuery.value) {
    return sortedArticles.value
  }
  return sortedArticles.value.filter(
    (article) =>
      article.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      article.summary.toLowerCase().includes(searchQuery.value.toLowerCase()),
  )
})

const featuredArticles = computed(() => filteredArticles.value.slice(0, 3))
const olderArticles = computed(() => filteredArticles.value.slice(3))

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}
</script>

<template>
  <div>
    <section class="section hexgrid" style="padding-bottom: var(--space-8)">
      <div class="container">
        <div class="section-head">
          <p class="overline">Insights</p>
          <h1>Notes on AI governance</h1>
          <p class="lead">Articles on AI governance, risk, compliance and the data foundations underneath.</p>
        </div>
        <div class="field" style="max-width: 560px; margin: 0">
          <label for="blog-search" class="sr-only">Search articles</label>
          <input id="blog-search" v-model="searchQuery" type="search" class="input" placeholder="Search articles by title or keyword" />
        </div>
      </div>
    </section>

    <section class="section" style="padding-top: var(--space-8)">
      <div class="container">
        <div v-if="isLoading" class="muted">Loading articles…</div>
        <div v-else-if="error" class="notice notice--error" role="alert">
          <p><strong>Articles couldn’t be loaded.</strong> Please try again later.</p>
        </div>
        <template v-else-if="filteredArticles.length > 0">
          <div v-if="featuredArticles.length > 0" class="grid grid-3">
            <article v-for="article in featuredArticles" :key="article.slug" class="card card--link post">
              <p class="meta">{{ formatDate(article.date) }}</p>
              <h3>
                <NuxtLink :to="`/blog/${article.slug}`" class="card__link" style="color: inherit; text-decoration: none">{{ article.title }}</NuxtLink>
              </h3>
              <p class="muted small clamp">{{ article.summary }}</p>
            </article>
          </div>
          <div v-if="olderArticles.length > 0" class="grid grid-2" style="margin-top: var(--space-12)">
            <article v-for="article in olderArticles" :key="article.slug" class="card card--link post">
              <p class="meta">{{ formatDate(article.date) }}</p>
              <h3 style="font-size: var(--text-lg)">
                <NuxtLink :to="`/blog/${article.slug}`" class="card__link" style="color: inherit; text-decoration: none">{{ article.title }}</NuxtLink>
              </h3>
              <p class="muted small clamp">{{ article.summary }}</p>
            </article>
          </div>
        </template>
        <div v-else class="muted">No articles match “{{ searchQuery }}”.</div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.clamp { display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; margin: 0; }
</style>
