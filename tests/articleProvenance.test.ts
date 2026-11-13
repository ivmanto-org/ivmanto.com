import { test } from 'node:test'
import assert from 'node:assert/strict'
import { articleProvenanceHead } from '../utils/articleProvenance.ts'

const article = {
  slug: 'ai-transparency-rules-now-apply',
  title: 'AI transparency rules now apply: what small companies must label',
  date: '2026-11-12',
}
const NOTHING = { meta: [], script: [] }

function jsonLd(value: string) {
  const head = articleProvenanceHead({ ...article, digital_source_type: value })
  assert.equal(head.script.length, 1)
  assert.equal(head.script[0].type, 'application/ld+json')
  return JSON.parse(head.script[0].innerHTML)
}

test('emits the CreativeWork block for an article created with generative AI', () => {
  assert.deepEqual(jsonLd('trainedAlgorithmicMedia'), {
    '@context': 'https://schema.org',
    '@type': 'CreativeWork',
    name: article.title,
    url: 'https://ivmanto.com/blog/ai-transparency-rules-now-apply',
    datePublished: '2026-11-12',
    creator: { '@type': 'Organization', name: 'Editorial team ivmanto.com' },
    publisher: { '@type': 'Organization', name: 'ivmanto.com' },
    usageInfo: 'https://iptc.org/standards/digital-source-type/',
    identifier: 'http://cv.iptc.org/newscodes/digitalsourcetype/trainedAlgorithmicMedia',
    about:
      'DigitalSourceType: trainedAlgorithmicMedia — Created using Generative AI, reviewed and edited by the editorial team, per the IPTC DigitalSourceType vocabulary and the EU Code of Practice on Transparency of AI-generated Content',
  })
})

test('uses the IPTC concept name of the other generative AI codes', () => {
  assert.match(jsonLd('compositeSynthetic').about, /— Composite including generative AI elements, reviewed/)
  const edited = jsonLd('compositeWithTrainedAlgorithmicMedia')
  assert.equal(edited.identifier, 'http://cv.iptc.org/newscodes/digitalsourcetype/compositeWithTrainedAlgorithmicMedia')
  assert.match(edited.about, /— Edited using Generative AI, reviewed/)
})

test('emits the meta tag with the vocabulary code', () => {
  const head = articleProvenanceHead({ ...article, digital_source_type: 'trainedAlgorithmicMedia' })
  assert.deepEqual(head.meta, [{ name: 'ai-digital-source-type', content: 'trainedAlgorithmicMedia' }])
})

test('emits nothing when the payload has no digital source type', () => {
  assert.deepEqual(articleProvenanceHead(article), NOTHING)
  assert.deepEqual(articleProvenanceHead(null), NOTHING)
})

test('emits nothing for non-generative codes, including algorithmicMedia (pure algorithm, no trained model)', () => {
  for (const value of ['algorithmicMedia', 'digitalCreation', 'humanEdits', 'human']) {
    assert.deepEqual(articleProvenanceHead({ ...article, digital_source_type: value }), NOTHING)
  }
})

test('emits nothing for a value outside the vocabulary', () => {
  for (const value of ['', 'TrainedAlgorithmicMedia', 'toString', '"></script><script>alert(1)']) {
    assert.deepEqual(articleProvenanceHead({ ...article, digital_source_type: value }), NOTHING)
  }
})

test('keeps a title with markup characters from closing the script element', () => {
  const head = articleProvenanceHead({ ...article, title: 'A </script><b> title', digital_source_type: 'trainedAlgorithmicMedia' })
  assert.ok(!head.script[0].innerHTML.includes('</script'))
  assert.equal(JSON.parse(head.script[0].innerHTML).name, 'A </script><b> title')
})
