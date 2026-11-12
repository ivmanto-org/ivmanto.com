// Service catalogue: single source in data/catalogue.json (v1, #112). The advisor backend
// validates ids against it. Packages carry a `line`: governance services (the three pillars),
// engineering (Data & AI engineering on Google Cloud) or platform (ivGIL).
import raw from '~/data/catalogue.json'

export type PillarId = 'governance' | 'risk' | 'compliance'
export type Line = 'governance' | 'engineering' | 'platform'

export interface Pillar {
  id: PillarId
  name: string
  /** Home page line, at most seven words (DS v1.1 cards.md). */
  tagline: string
  /** Pillar pages. */
  summary: string
}

export interface Package {
  id: string
  line: Line
  /** A governance pillar id, or 'engineering' for the engineering line. */
  pillar: PillarId | 'engineering'
  name: string
  outcome: string
  duration_weeks: number
  /** Overrides "N weeks" where weeks don't fit (1 day, monthly). */
  duration_label?: string
  deliverable: string
  includes: string[]
  not_included: string[]
  frameworks: string[]
  price_eur: number | null
  /** 'month' for a retainer; absent for a fixed one-off price. */
  price_unit?: 'month'
  price_note?: string
  recommended: boolean
  /** The package featured for its pillar (home, pillar page). */
  pillar_recommended?: boolean
}

interface Catalogue {
  version: number
  pillars: Pillar[]
  packages: Package[]
}

const catalogue = raw as unknown as Catalogue

const priceFormat = new Intl.NumberFormat('en-IE', { style: 'currency', currency: 'EUR', maximumFractionDigits: 0 })

export const PRICE_PLACEHOLDER = '€ [price]'

// Nick, 2026-11-12: prices stay placeholders until the market research and the GTM strategy confirm them.
// The catalogue ships in the bundle, so it carries no figures: price_eur is null until then.
export function formatPrice(pkg: Pick<Package, 'price_eur' | 'price_unit'>): string {
  if (pkg.price_eur == null) return PRICE_PLACEHOLDER
  const amount = priceFormat.format(pkg.price_eur)
  return pkg.price_unit === 'month' ? `${amount}/month` : amount
}

export function priceTerms(pkg: Pick<Package, 'price_unit'>): string {
  return pkg.price_unit === 'month' ? 'per month, excl. VAT' : 'fixed, excl. VAT'
}

export function formatDuration(pkg: Pick<Package, 'duration_weeks' | 'duration_label'>): string {
  return pkg.duration_label ?? `${pkg.duration_weeks} weeks`
}

export const LINE_NAMES: Record<Line, string> = {
  governance: 'AI governance',
  engineering: 'Data & AI engineering',
  platform: 'ivGIL',
}

export function useCatalogue() {
  const pillars = catalogue.pillars
  const packages = catalogue.packages

  const pillar = (id: string) => pillars.find((p) => p.id === id)
  const pkg = (id: string) => packages.find((p) => p.id === id)
  const packagesFor = (pillarId: string) => packages.filter((p) => p.pillar === pillarId)
  const packagesInLine = (line: Line) => packages.filter((p) => p.line === line)
  const featured = () => {
    // Home shows one package per pillar: the site-wide recommended start, else the
    // pillar's recommended package, else its first.
    return pillars
      .map((pl) => {
        const list = packagesFor(pl.id)
        return list.find((p) => p.recommended) ?? list.find((p) => p.pillar_recommended) ?? list[0]
      })
      .filter((p): p is Package => Boolean(p))
  }

  return { pillars, packages, pillar, pkg, packagesFor, packagesInLine, featured }
}
