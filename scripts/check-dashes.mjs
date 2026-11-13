// Fails when em or en dashes appear in the services section copy.
// Usage: node scripts/check-dashes.mjs [extra files or dirs...]
import { readFileSync, readdirSync, statSync } from 'node:fs'
import { join } from 'node:path'

const targets = [
  'data/services.ts',
  'pages/services',
  'components/services',
  'components/services-content',
  'components/layout/AppHeader.vue',
  'composables/usePageMetadata.ts',
  // Relaunch (v0.2.0): funnel, catalogue and advisor copy.
  'pages/index.vue',
  'pages/packages',
  'pages/order',
  'pages/ai-use.vue',
  'components/site',
  'components/advisor',
  'composables/useCatalogue.ts',
  'data/catalogue.json',
  ...process.argv.slice(2),
]

const pattern = /[–—]|&mdash;|&ndash;|&#821[12];/

function files(path) {
  if (statSync(path).isFile()) return [path]
  return readdirSync(path).flatMap((name) => files(join(path, name)))
}

const hits = []
for (const file of targets.flatMap(files)) {
  readFileSync(file, 'utf8')
    .split('\n')
    .forEach((line, i) => {
      if (pattern.test(line)) hits.push(`${file}:${i + 1}: ${line.trim()}`)
    })
}

if (hits.length) {
  console.error(`Found ${hits.length} em/en dash(es):\n${hits.join('\n')}`)
  process.exit(1)
}
console.log('No em/en dashes found.')
