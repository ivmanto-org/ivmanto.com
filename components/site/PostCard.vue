<script setup lang="ts">
// PostCard (DS v1.5, pages/blog.md): one meta line, the title as the only link, summary clamped to two lines.
import type { ArticleMeta } from '~/types/article'
import { cardReadingMinutes, formatArticleDate, topicLabel } from '~/utils/insights'

const props = defineProps<{ article: ArticleMeta }>()

const topic = computed(() => topicLabel(props.article.topic))
const minutes = computed(() => cardReadingMinutes(props.article))
</script>

<template>
  <article class="card card--link post">
    <p class="meta">
      <template v-if="topic"><b>{{ topic }}</b> · </template>{{ formatArticleDate(article.date) }}<template v-if="minutes"> · {{ minutes }} min</template>
    </p>
    <h3><NuxtLink :to="`/blog/${article.slug}`" class="card__link">{{ article.title }}</NuxtLink></h3>
    <p class="post__sum">{{ article.summary }}</p>
  </article>
</template>
