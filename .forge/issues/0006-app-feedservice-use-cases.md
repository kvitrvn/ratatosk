---
id: 6
title: '[app] FeedService — use cases'
status: in-progress
type: task
priority: medium
author: mcp-agent
milestone: M3 – Fetcher & application service
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Orchestrer les use cases métier, testables sans TUI.

**Fichiers :**
- `internal/application/service.go` — `FeedService`
- `internal/application/service_test.go` — stubs in-memory (pas de mock lib)

**Méthodes :**
- `Subscribe(rawURL string) (Feed, error)` — valide, persiste, premier fetch
- `Refresh(feedID int64) error`
- `RefreshAll() []error` — collecte les erreurs, ne s'arrête pas à la première
- `GetArticles(feedID int64) ([]Article, error)`
- `MarkRead(articleID int64, read bool) error`

**Dépend de :** #0005
