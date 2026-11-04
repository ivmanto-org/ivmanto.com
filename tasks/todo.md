# Blog CamelCase Slug Migration — 10 articles → kebab-case + redirects

## Goal

Bring all 10 remaining CamelCase article slugs into compliance with the `kebab-case` rule in `docs/blog-article-spec.md:7`, and add 301 redirects so the old URLs don't 404 (preserving backlinks, Google search results, and any shared links). Update 8 source-file references that hardcode the old slugs so the cross-link from /services and /about pages resolves to the new URL.

The whole task is shipped as ONE PR (single bundle) so the rename and the redirect can never drift apart. Branch protection requires checks to pass; this is not a doc-only change.

## Constraints (from AGENTS.md + .agents/rules.md)

- Branch: `dev-v0.1.4` (next version after `0.1.3`), never push to main
- Conventional commits, lowercase, scoped
- Pre-PR checklist: `npm run lint`, `npm run generate` green; no secrets in diff
- Zero content regression (the article body, titles, dates, summaries are unchanged)
- Blog system integrity: GCS bucket pipeline preserved; the Go backend's filename-derived slug lookup continues to work because the new files are still `.md` in `ivmanto_com_blog_articles`

## The 10 slugs (rename map)

| Old (CamelCase) | New (kebab-case) | Source references to update |
|---|---|---|
| `AiMarketSentiment` | `ai-market-sentiment` | none in source |
| `DataMeshGovernance` | `data-mesh-governance` | `data/services.ts:167`, `pages/services/index.vue:98`, `pages/about.vue:138` |
| `DataProductFallacy` | `data-product-fallacy` | none in source |
| `DigitalisationSimplicity` | `digitalisation-simplicity` | none in source (British spelling kept) |
| `FromBigDataTo` | `from-big-data-to` | `data/services.ts:49`, `data/services.ts:116` |
| `NavigatingTheDataFrontier` | `navigating-the-data-frontier` | `data/services.ts:99` |
| `OnDataManagement` | `on-data-management` | `data/services.ts:99`, `data/services.ts:167` |
| `SqlVsNoSql2026` | `sql-vs-nosql-2026` | none in source (year kept) |
| `TechnicalDebtAsset` | `technical-debt-asset` | none in source |
| `VisionaryDataArchitecture` | `visionary-data-architecture` | `data/services.ts:49`, `data/services.ts:135`, `pages/services/index.vue:105` |

**Total source-file lines to change: 8** (5 in `data/services.ts`, 2 in `pages/services/index.vue`, 1 in `pages/about.vue`).

## Implementation Plan

### Phase 0 — Branch and prep
- [ ] **0.1** `cd /home/ivmo/projects/ivmanto.com && git checkout main && git pull`
- [ ] **0.2** `git checkout -b dev-v0.1.4`
- [ ] **0.3** Bump `package.json` version 0.1.3 → 0.1.4 (per branch-naming convention)
- [ ] **0.4** Confirm GCS credential path: `/home/ivmo/.hermes/profiles/ivmo/.secrets/google-service-account.json` exists, and `python3 -c "from google.cloud import storage; storage.Client.from_service_account_json('...').bucket('ivmanto_com_blog_articles')"` works (read-only probe)

### Phase 1 — Update source cross-links (safe to do first; can't break anything because the new URLs don't exist yet — they only 404 like the old ones do today)
- [ ] **1.1** `data/services.ts` — 5 `relatedBlogSlugs` array updates
- [ ] **1.2** `pages/services/index.vue` — 2 `to="/blog/..."` attribute updates
- [ ] **1.3** `pages/about.vue` — 1 `to="/blog/..."` attribute update
- [ ] **1.4** Verify: `npm run lint` passes (no errors from the cross-link changes)
- [ ] **1.5** Verify: `git grep -nE 'blog/(AiMarketSentiment|DataMeshGovernance|DataProductFallacy|DigitalisationSimplicity|FromBigDataTo|NavigatingTheDataFrontier|OnDataManagement|SqlVsNoSql2026|TechnicalDebtAsset|VisionaryDataArchitecture)' -- 'data/' 'pages/' 'components/' 'composables/'` returns ZERO matches (other than the 8 we just updated)

### Phase 2 — Add 301 redirects in public/serve.json

**Important: this PR #96 (`cf51162`, dev-v0.1.9, merged 2026-07-28) already shipped redirects for a DIFFERENT set of legacy slugs** (the date-prefixed kebab ones from an earlier cleanup, plus 3 PascalCase slugs that have no successor). The redirect layer on Cloud Run is **`public/serve.json`** (picked up by the `serve` binary in the frontend container), NOT `nuxt.config.ts` routeRules — the SSG output does NOT honor Nuxt runtime routeRules, only the static `serve.json` config. So the right place to add 10 more redirect entries is the bottom of the existing `public/serve.json` array, matching the established shape.

- [ ] **2.1** Add 10 entries to the `redirects` array in `public/serve.json` (one per old CamelCase slug, status 301). Each entry follows the exact `{ source, destination, type }` shape of the existing entries:
  ```json
  { "source": "/blog/AiMarketSentiment",     "destination": "/blog/ai-market-sentiment",     "type": 301 },
  { "source": "/blog/DataMeshGovernance",    "destination": "/blog/data-mesh-governance",    "type": 301 },
  { "source": "/blog/DataProductFallacy",    "destination": "/blog/data-product-fallacy",    "type": 301 },
  { "source": "/blog/DigitalisationSimplicity", "destination": "/blog/digitalisation-simplicity", "type": 301 },
  { "source": "/blog/FromBigDataTo",         "destination": "/blog/from-big-data-to",         "type": 301 },
  { "source": "/blog/NavigatingTheDataFrontier", "destination": "/blog/navigating-the-data-frontier", "type": 301 },
  { "source": "/blog/OnDataManagement",      "destination": "/blog/on-data-management",      "type": 301 },
  { "source": "/blog/SqlVsNoSql2026",        "destination": "/blog/sql-vs-nosql-2026",        "type": 301 },
  { "source": "/blog/TechnicalDebtAsset",    "destination": "/blog/technical-debt-asset",    "type": 301 },
  { "source": "/blog/VisionaryDataArchitecture", "destination": "/blog/visionary-data-architecture", "type": 301 }
  ```
  Order: alphabetical by source. No trailing-slash pairs — `serve` matches the exact `source` string, and the existing PR #96 entries follow the same single-source shape.
- [ ] **2.2** Verify: `python3 -c "import json; d = json.load(open('public/serve.json')); print(len(d['redirects']))"` returns 18 (8 existing + 10 new)
- [ ] **2.3** Verify: `python3 -c "import json; d = json.load(open('public/serve.json')); print([r for r in d['redirects'] if r['source'].endswith('AiMarketSentiment')])"` returns the new entry

### Phase 3 — GCS rename (the load-bearing step)
- [ ] **3.1** Write a small one-off Python script `scripts/rename_articles.py` that uses the service-account credential to:
  1. List all 10 old blob names, print their sizes + etags (audit trail)
  2. Copy each old blob to the new kebab name (preserves content + metadata)
  3. Verify the copy (byte-compare) and confirm `published: true` in the new blob's bytes
  4. Delete the old blob
  5. Print a CSV report: `old,new,bytes,etag,old_deleted`
- [ ] **3.2** Run the script with a dry-run flag first (`--dry-run`) and verify the rename map against the bucket list
- [ ] **3.3** Run the script for real (no flag)
- [ ] **3.4** Verify: the script's report shows all 10 old names deleted AND all 10 new names present AND `published: true` preserved
- [ ] **3.5** Spot-check 2 of the 10 new files via `verify_gcs_publish.py <local_path>` if a local copy exists, otherwise via a direct `cat` of the GCS object through the service account (NOT a public `curl` — the bucket is private, anonymous 403 is expected per the blog-publishing skill's "GCS Bucket Is Private" pitfall)

### Phase 4 — Re-build the static site + verify
- [ ] **4.1** `npm run generate` — must complete without errors. The `prerender:routes` hook at `nuxt.config.ts:127-138` fetches `/api/articles` from the live backend; the backend's in-memory cache will pick up the new filenames within seconds via Pub/Sub, so the new slugs will be prerendered. If the cache is stale, wait 30s and re-run.
- [ ] **4.2** Verify `.output/public/blog/` contains all 10 new kebab-case directories AND none of the old CamelCase directories remain. Use:
  ```bash
  ls .output/public/blog/ | grep -E '^[A-Z]'
  ```
  Expect ZERO output.
- [ ] **4.3** Verify `.output/public/blog/` contains the 10 new kebab directories:
  ```bash
  for s in ai-market-sentiment data-mesh-governance data-product-fallacy \
           digitalisation-simplicity from-big-data-to navigating-the-data-frontier \
           on-data-management sql-vs-nosql-2026 technical-debt-asset \
           visionary-data-architecture; do
    test -f ".output/public/blog/$s/index.html" || echo "MISSING: $s"
  done
  ```
  Expect ZERO output.
- [ ] **4.4** Verify `.output/public/sitemap.xml` lists the new slugs and NOT the old ones:
  ```bash
  grep -oE 'blog/[A-Z][a-zA-Z0-9]*' .output/public/sitemap.xml && echo "STALE: old slugs still in sitemap"
  ```
  Expect ZERO output.
- [ ] **4.5** Verify the redirect entries are picked up by Cloud Run's `serve`: `serve.json` is copied into `.output/public/` during the build. Confirm the new entries survive: `python3 -c "import json; print(len(json.load(open('.output/public/serve.json'))['redirects']))"` returns 18. (The Nuxt-side routeRules redirect question I worried about earlier is moot — the live site uses `serve` + `serve.json`, not Nuxt's runtime.)

### Phase 5 — Commit, push, PR
- [ ] **5.1** `git add -A && git commit -m "fix(blog): rename 10 CamelCase slugs to kebab-case + add 301 redirects"`
- [ ] **5.2** `git push -u origin dev-v0.1.4`
- [ ] **5.3** Open PR against `main` using `.github/pull_request_template.md`. PR body must include:
  - The rename map table from this doc
  - Output of `npm run lint` and `npm run generate` (paste verbatim)
  - The script's audit CSV (10 lines)
  - The Phase 4 verification grep results (all 3 must be empty)
- [ ] **5.4** Address review feedback (reviewer = Claude, per the AGENTS.md contract)

### Phase 6 — Post-merge smoke (this is the "is it actually live" check)
- [ ] **6.1** After Nick merges to main and Cloud Build runs, run the standard `verify_live_publish.py ai-market-sentiment` (and the other 9). Each must exit 0. The skill says these are conditional on exit 0, not on GCS upload success.
- [ ] **6.2** Manually verify one of the 10 old URLs with a HEAD probe (per the ivmanto.com post-merge pattern from the memory entry on GCS =/= live): `curl -I https://ivmanto.com/blog/TechnicalDebtAsset` must return 301 with `Location: /blog/technical-debt-asset`. Repeat for one more (DataMeshGovernance) for sanity.
- [ ] **6.3** Tell Nick in the main session: "CamelCase → kebab migration merged + live. 10/10 slugs return exit 0 from `verify_live_publish.py`, 2/10 old URLs 301 to new URLs on the live site. Old bookmarks preserved. PR review close the loop on the next turn."

## Risks & Mitigations

1. **Race with the Go backend's in-memory cache.** The backend reads articles from GCS and caches them by slug (`backend/internal/blog/cache.go:24`). When we rename, the cache will still hold the old slugs until the next Pub/Sub refresh (a few seconds). The backend's `/api/articles` will return both old + new slugs during that window. Mitigation: the Nuxt `prerender:routes` hook waits 0s; we add a 60s sleep between Phase 3.3 and Phase 4.1 to let the cache settle.
2. **Cloud Run's static serving doesn't honor `routeRules` redirects for the prerendered output.** Updated after live probe: Cloud Run uses `serve` (the static-file server), which honors `public/serve.json`'s `redirects` array, NOT Nuxt `routeRules`. The right place for the redirects is `public/serve.json` (the established pattern, per PR #96). Phase 4.5 verifies the entries survive the Nuxt build.
3. **Google Search Console indexing the old slugs.** After the rename + 301s ship, Google's crawler will see the 301s and re-index under the new URLs. This takes days to weeks; during that window, search results may still show the old URLs and the 301 will pass PageRank correctly. Not a correctness issue, just a slow burn. Documented in the PR body.
4. **Old CamelCase URLs in the SEO audit HTMLs (`docs/SEO-reports/SEO Audit Ivmanto_2.html`, `…RankClaw.html`).** These are historical artifacts, not linked from anywhere live. They reference the old slugs because the audit was run when the slugs were still CamelCase. We will NOT touch these in this PR (they're the audit's evidence of the prior state); a follow-up chore PR can update them or annotate them as historical.

## Out of scope (explicit)

- Historical doc updates (`docs/SEO-reports/*`, `docs/features/dev-v0.0.17/*`) — separate chore PR
- Plagiarism re-check on the 10 articles (no content changed)
- Any change to the Go backend (slug derivation in `cache.go:116` continues to work for the new filenames)
- Any change to `cloudbuild.yaml`, the `dist` symlink, or Secret Manager

---

# Open tasks

- [ ] **Share your Applied Agentic AI for Organizational Transformation credential**
      Accredible URL: https://www.credential.net/100ff9ab-e82e-41ee-b16e-4ce596b190bb#acc.lWqCGLPj
