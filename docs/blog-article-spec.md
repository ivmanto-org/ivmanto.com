# Blog Article Specification

This document defines the frontmatter contract for blog articles published to the GCS bucket `ivmanto-com-prod-514215-blog`.

## Filename Convention

- Use **kebab-case**: `the-title-of-the-article.md`
- Must have `.md` extension
- No spaces, no uppercase letters in filenames

## Required Frontmatter Fields

Every article **must** include all of the following fields in the YAML frontmatter block:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | **YES** | The article title displayed on the page |
| `summary` | string | **YES** | A 1-2 sentence summary shown in article lists |
| `date` | string (YYYY-MM-DD) | **YES** | Publication date |
| `published` | boolean | **YES** | Must be `true` for the article to appear on the site |

## Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `author` | string | Author name (defaults to "Ivmanto") |
| `tags` | list | Topic tags for categorization |

## Template

```markdown
---
title: "Your Article Title Here"
summary: "A concise 1-2 sentence summary of what this article covers."
date: 2026-03-02
published: true
author: Ivmanto
tags: [Data, AI, GCP]
---

# Your Article Title Here

Article content in Markdown...
```

## Validation Rules

1. **`published` must be explicitly set to `true`** for the article to appear on the website. If this field is missing or set to `false`, the article is silently skipped.
2. **`summary` is required** for article list cards. Articles without a summary will still publish but display an empty summary.
3. **`date` must be a valid date** in `YYYY-MM-DD` format. Articles are sorted by date (newest first).
4. **`title` is required** for display. Articles without a title will show an empty heading.

## Diagnostic Endpoint

To check why an article isn't appearing, query:

```
GET /api/_internal/articles/status?token=<PUSH_TOKEN>
```

This returns a JSON object with:
- `total_files`: number of `.md` files found in the bucket
- `published`: number of articles currently visible
- `skipped`: list of articles that were skipped, with reasons
