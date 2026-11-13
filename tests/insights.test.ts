import { test } from 'node:test'
import assert from 'node:assert/strict'
import {
  cardReadingMinutes, clampPage, filterArticles, formatArticleDate, jsonLd, markSources, pageCount, pageSlice,
  pagerItems, readingMinutes, relatedArticles, topicLabel, topicsInUse,
} from '../utils/insights.ts'

const make = (slug: string, topic?: string, title = slug, summary = '') =>
  ({ slug, title, summary, date: '2026-09-18', published: true, topic })

test('formats dates as day, short month, year', () => {
  assert.equal(formatArticleDate('2026-09-18'), '18 Sep 2026')
  assert.equal(formatArticleDate('2026-03-09T00:00:00Z'), '9 Mar 2026')
  assert.equal(formatArticleDate('not a date'), '')
})

test('reading time is 200 words a minute, rounded up, never below 1', () => {
  assert.equal(readingMinutes(''), 1)
  assert.equal(readingMinutes('<p>' + 'word '.repeat(200) + '</p>'), 1)
  assert.equal(readingMinutes('<p>' + 'word '.repeat(201) + '</p>'), 2)
  assert.equal(readingMinutes('<h2 class="x">one</h2><p>two</p>'), 1)
})

test('a card shows minutes only when the API sends a usable number', () => {
  assert.equal(cardReadingMinutes({ ...make('a'), reading_minutes: 6 }), 6)
  assert.equal(cardReadingMinutes(make('a')), null)
  assert.equal(cardReadingMinutes({ ...make('a'), reading_minutes: 0 }), null)
})

test('filters by topic and searches title and summary inside it', () => {
  const list = [make('a', 'risk', 'Risk register'), make('b', 'risk', 'Vendors', 'A register of suppliers'), make('c', 'governance', 'Register of owners'), make('d')]
  assert.deepEqual(filterArticles(list, '', '').map((a) => a.slug), ['a', 'b', 'c', 'd'])
  assert.deepEqual(filterArticles(list, 'risk', '').map((a) => a.slug), ['a', 'b'])
  assert.deepEqual(filterArticles(list, 'risk', ' REGISTER ').map((a) => a.slug), ['a', 'b'])
  assert.deepEqual(filterArticles(list, 'governance', 'register').map((a) => a.slug), ['c'])
  assert.deepEqual(filterArticles(list, 'unknown', '').map((a) => a.slug), ['a', 'b', 'c', 'd'])
})

test('topic labels and the topics in use', () => {
  assert.equal(topicLabel('risk'), 'Risk')
  assert.equal(topicLabel('toString'), '')
  assert.equal(topicLabel(undefined), '')
  assert.deepEqual(topicsInUse([make('a', 'risk'), make('b', 'governance'), make('c', 'other')]), ['governance', 'risk'])
  assert.deepEqual(topicsInUse([make('a')]), [])
})

test('pages: count, clamp and slice', () => {
  assert.equal(pageCount(0), 1)
  assert.equal(pageCount(12), 1)
  assert.equal(pageCount(13), 2)
  assert.equal(clampPage('2', 5), 2)
  assert.equal(clampPage('99', 5), 5)
  assert.equal(clampPage('-1', 5), 1)
  assert.equal(clampPage(undefined, 5), 1)
  assert.equal(clampPage('abc', 5), 1)
  assert.deepEqual(pageSlice([1, 2, 3, 4, 5], 2, 2), [3, 4])
})

test('pager shows every page up to 7, then first, last and two around the current one', () => {
  assert.deepEqual(pagerItems(1, 5), [1, 2, 3, 4, 5])
  assert.deepEqual(pagerItems(4, 7), [1, 2, 3, 4, 5, 6, 7])
  assert.deepEqual(pagerItems(1, 10), [1, 2, 3, 'gap', 10])
  assert.deepEqual(pagerItems(5, 10), [1, 'gap', 3, 4, 5, 6, 7, 'gap', 10])
  assert.deepEqual(pagerItems(10, 10), [1, 'gap', 8, 9, 10])
  assert.deepEqual(pagerItems(3, 10), [1, 2, 3, 4, 5, 'gap', 10])
})

test('related: same topic first, filled from the rest, never the current article', () => {
  const list = [make('a', 'risk'), make('b', 'governance'), make('c', 'risk'), make('d'), make('e', 'risk')]
  assert.deepEqual(relatedArticles(list, list[0]!).map((a) => a.slug), ['c', 'e', 'b'])
  assert.deepEqual(relatedArticles(list, list[3]!).map((a) => a.slug), ['a', 'b', 'c'])
  assert.deepEqual(relatedArticles([list[0]!], list[0]!), [])
})

test('marks a Sources paragraph and escapes JSON-LD', () => {
  assert.equal(markSources('<p>Text</p><p><strong>Sources:</strong> A</p>'), '<p>Text</p><p class="sources"><strong>Sources:</strong> A</p>')
  assert.equal(markSources('<p>Sources of risk</p>'), '<p>Sources of risk</p>')
  assert.ok(!jsonLd({ t: '</script><script>' }).includes('<'))
  assert.deepEqual(JSON.parse(jsonLd({ t: '</script>' })), { t: '</script>' })
})
