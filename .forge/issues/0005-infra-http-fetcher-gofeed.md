---
id: 5
title: '[infra] HTTP fetcher gofeed'
status: in-progress
type: task
priority: medium
author: mcp-agent
milestone: M3 – Fetcher & application service
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Adapter gofeed derrière l'interface `Fetcher`.

**Fichiers :**
- `internal/infrastructure/fetcher/gofeed_fetcher.go` — implémente `Fetcher`

**Interface `Fetcher` définie dans `application` (pas `infrastructure`) :**
```go
type Fetcher interface {
    Fetch(url string) (FetchedFeed, error)
}
```

**Dépend de :** #0004
