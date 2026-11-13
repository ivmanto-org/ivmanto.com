// Machine-readable AI provenance marking for article pages.
// Contract: spec-ai-provenance-marking.md (Ivmo, 2026-11-13), section 3.
// The value comes from the article payload (`digital_source_type`, a code of the
// IPTC Digital Source Type vocabulary, https://cv.iptc.org/newscodes/digitalsourcetype/).
// Only the generative AI codes below are marked; unmarked articles and any other
// value emit nothing.

const SITE_URL = 'https://ivmanto.com'
const IPTC_CODE_BASE = 'http://cv.iptc.org/newscodes/digitalsourcetype/'
const IPTC_VOCABULARY_INFO = 'https://iptc.org/standards/digital-source-type/'
const CREATOR_NAME = 'Editorial team ivmanto.com'
const PUBLISHER_NAME = 'ivmanto.com'
export const PROVENANCE_META_NAME = 'ai-digital-source-type'

// Code -> IPTC concept name, quoted from the vocabulary.
const MARKED_SOURCE_TYPES: Record<string, string> = {
  trainedAlgorithmicMedia: 'Created using Generative AI',
  compositeSynthetic: 'Composite including generative AI elements',
  compositeWithTrainedAlgorithmicMedia: 'Edited using Generative AI',
}

export interface ProvenanceSource {
  slug: string
  title: string
  date: string
  digital_source_type?: string | null
}

export interface ProvenanceHead {
  meta: { name: string; content: string }[]
  script: { type: string; innerHTML: string }[]
}

export function articleProvenanceHead(article: ProvenanceSource | null | undefined): ProvenanceHead {
  const value = article?.digital_source_type
  if (!article || !value || !Object.hasOwn(MARKED_SOURCE_TYPES, value)) return { meta: [], script: [] }

  const jsonLd = {
    '@context': 'https://schema.org',
    '@type': 'CreativeWork',
    name: article.title,
    url: `${SITE_URL}/blog/${article.slug}`,
    datePublished: article.date,
    creator: { '@type': 'Organization', name: CREATOR_NAME },
    publisher: { '@type': 'Organization', name: PUBLISHER_NAME },
    usageInfo: IPTC_VOCABULARY_INFO,
    identifier: `${IPTC_CODE_BASE}${value}`,
    about: `DigitalSourceType: ${value} — ${MARKED_SOURCE_TYPES[value]}, reviewed and edited by the editorial team, per the IPTC DigitalSourceType vocabulary and the EU Code of Practice on Transparency of AI-generated Content`,
  }

  return {
    meta: [{ name: PROVENANCE_META_NAME, content: value }],
    // "<" is escaped so no field can close the script element.
    script: [{ type: 'application/ld+json', innerHTML: JSON.stringify(jsonLd).replace(/</g, '\\u003c') }],
  }
}
