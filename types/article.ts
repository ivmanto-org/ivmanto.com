export interface ArticleMeta {
  slug: string
  title: string
  summary: string
  date: string
  published: boolean
  /** One of governance, risk, compliance, engineering (DS v1.5); an article without one shows under "All" only. */
  topic?: string
  /** Reading time in minutes at 200 words a minute, computed by the API from the article text. */
  reading_minutes?: number
  /** IPTC Digital Source Type code (e.g. trainedAlgorithmicMedia); set by the publishing pipeline. */
  digital_source_type?: string
}

export interface Article extends ArticleMeta {
  content: string
}
