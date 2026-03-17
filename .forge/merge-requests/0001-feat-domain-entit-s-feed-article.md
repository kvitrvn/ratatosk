---
id: 1
title: 'feat(domain): entités Feed & Article'
status: merged
source_branch: feat/issue-0001-domain-entities
target_branch: main
author: mcp-agent
linked_issues:
    - 1
approvals:
    - reviewer: Kvitrvn
      approved_at: "2026-03-17"
approvals_required: 1
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Closes #0001

## Changements

- `go.mod` / `go.sum` — toutes les dépendances du projet ajoutées
- `internal/domain/feed.go` — `Feed{ID, URL, Title, CreatedAt}` + `NewFeed` (validation http/https)
- `internal/domain/article.go` — `Article{...}` + `NewArticle`
- `internal/domain/feed_test.go` — 5 cas couvrant URLs valides et invalides

## Vérification

```
go test ./internal/domain/...  ✓
go vet ./internal/domain/...   ✓
```

## Reviews

### Review 1

Reviewer: Kvitrvn
Date: 2026-03-17
Status: commented

c'est ok
