// Service catalogue: single source in data/catalogue.json (placeholder v0 until Nick
// defines names, scope and prices, #112). The advisor backend validates ids against it.
import raw from '~/data/catalogue.json'

export type PillarId = 'governance' | 'risk' | 'compliance'

export interface Pillar {
  id: PillarId
  name: string
  summary: string
}

export interface Package {
  id: string
  pillar: PillarId
  name: string
  outcome: string
  duration_weeks: number
  deliverable: string
  includes: string[]
  not_included: string[]
  frameworks: string[]
  price_eur: number | null
  recommended: boolean
  /** The package featured for its pillar (home, pillar page). */
  pillar_recommended?: boolean
}

interface Catalogue {
  version: number
  currency: string
  pillars: Pillar[]
  packages: Package[]
}

const catalogue = raw as unknown as Catalogue

const priceFormat = new Intl.NumberFormat('en-IE', { style: 'currency', currency: 'EUR', maximumFractionDigits: 0 })

export function formatPrice(pkg: Pick<Package, 'price_eur'>): string {
  return pkg.price_eur == null ? '€ [price]' : priceFormat.format(pkg.price_eur)
}

export function useCatalogue() {
  const pillars = catalogue.pillars
  const packages = catalogue.packages

  const pillar = (id: string) => pillars.find((p) => p.id === id)
  const pkg = (id: string) => packages.find((p) => p.id === id)
  const packagesFor = (pillarId: string) => packages.filter((p) => p.pillar === pillarId)
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

  return { pillars, packages, pillar, pkg, packagesFor, featured }
}
