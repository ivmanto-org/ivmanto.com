export interface ArticleMeta {
  slug: string
  title: string
  summary: string
  date: string
  published: boolean
  /** IPTC DigitalSourceType short name; set by the publishing pipeline. */
  digital_source_type?: string
}

export interface Article extends ArticleMeta {
  content: string
}
