---
id: 1
title: '[domain] Entités Feed & Article'
status: closed
type: task
priority: medium
author: mcp-agent
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Créer les entités du domaine et mettre à jour `go.mod` avec toutes les dépendances.

**Fichiers :**
- `go.mod` — ajouter bubbletea, lipgloss, bubbles, gofeed, sqlite, cobra, viper
- `internal/domain/feed.go` — `Feed{ID, URL, Title, CreatedAt}` + `NewFeed(rawURL)` (validation http/https)
- `internal/domain/article.go` — `Article{ID, FeedID, GUID, Title, Link, Description, PublishedAt, Read}` + `NewArticle(...)`
- `internal/domain/feed_test.go` — tests NewFeed (URL invalide, URL valide)
