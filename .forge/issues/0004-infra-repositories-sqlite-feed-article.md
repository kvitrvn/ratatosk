---
id: 4
title: '[infra] Repositories SQLite (feed + article)'
status: open
type: task
priority: medium
author: mcp-agent
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Implémenter les interfaces repository sur SQLite avec tests.

**Fichiers :**
- `internal/infrastructure/db/feed_repo.go` — implémente `FeedRepository`
- `internal/infrastructure/db/article_repo.go` — implémente `ArticleRepository`
- `internal/infrastructure/db/feed_repo_test.go`
- `internal/infrastructure/db/article_repo_test.go`

**Notes :**
- `INSERT OR IGNORE` pour déduplication par GUID
- Tests avec `t.TempDir()` sur vrai fichier SQLite (pas de mock)

**Dépend de :** #0003
