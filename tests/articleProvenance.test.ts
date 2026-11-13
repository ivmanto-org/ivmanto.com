import { test } from 'node:test'
import assert from 'node:assert/strict'
import { articleProvenanceHead } from '../utils/articleProvenance.ts'

const article = {
  slug: 'ai-transparency-rules-now-apply',
  title: 'AI transparency rules now apply: what small companies must label',
  date: '2026-11-12',
}

function jsonLd(value: string) {
  const head = articleProvenanceHead({ ...article, digital_source_type: value })
  assert.equal(head.script.length, 1)
  assert.equal(head.script[0].type, 'application/ld+json')
  return JSON.parse(head.script[0].innerHTML)
}

test('emits the CreativeWork block from the spec for an AI-assisted article', () => {
  assert.deepEqual(jsonLd('algorithmicMedia'), {
    '@context': 'https://schema.org',
    '@type': 'CreativeWork',
    name: article.title,
    url: 'https://ivmanto.com/blog/ai-transparency-rules-now-apply',
    datePublished: '2026-11-12',
    creator: { '@type': 'Organization', name: 'Editorial team ivmanto.com' },
    publisher: { '@type': 'Organization', name: 'ivmanto.com' },
    usageInfo: 'https://iptc.org/standards/digital-source-type/',
    identifier: 'http://cv.iptc.org/codes/digitalSourceType/algorithmicMedia',
    about:
      'DigitalSourceType: algorithmicMedia — AI-assisted, human-reviewed per the IPTC DigitalSourceType vocabulary and the EU Code of Practice on Transparency of AI-generated Content',
  })
})

test('describes a purely generated article as generated with AI', () => {
  const data = jsonLd('trainedAlgorithmicMedia')
  assert.equal(data.identifier, 'http://cv.iptc.org/codes/digitalSourceType/trainedAlgorithmicMedia')
  assert.match(data.about, /^DigitalSourceType: trainedAlgorithmicMedia — generated with AI per the IPTC/)
})

test('emits the meta tag with the short vocabulary name', () => {
  const head = articleProvenanceHead({ ...article, digital_source_type: 'algorithmicMedia' })
  assert.deepEqual(head.meta, [{ name: 'ai-digital-source-type', content: 'algorithmicMedia' }])
})

test('emits nothing for a human-written article', () => {
  assert.deepEqual(articleProvenanceHead({ ...article, digital_source_type: 'human' }), { meta: [], script: [] })
})

test('emits nothing when the payload has no digital source type', () => {
  assert.deepEqual(articleProvenanceHead(article), { meta: [], script: [] })
  assert.deepEqual(articleProvenanceHead(null), { meta: [], script: [] })
})

test('emits nothing for a value outside the IPTC vocabulary we use', () => {
  for (const value of ['', 'AlgorithmicMedia', 'compositeSynthetic', '"></script><script>alert(1)']) {
    assert.deepEqual(articleProvenanceHead({ ...article, digital_source_type: value }), { meta: [], script: [] })
  }
})

test('keeps a title with markup characters from closing the script element', () => {
  const head = articleProvenanceHead({ ...article, title: 'A </script><b> title', digital_source_type: 'algorithmicMedia' })
  assert.ok(!head.script[0].innerHTML.includes('</script'))
  assert.equal(JSON.parse(head.script[0].innerHTML).name, 'A </script><b> title')
})
