---
id: 9
title: '[tui] Panneaux FeedList & ArticleList'
status: open
type: task
priority: medium
author: mcp-agent
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Implémenter les deux panneaux de liste.

**Fichiers :**
- `internal/interfaces/tui/feedlist.go` — panneau gauche (30% de la largeur)
- `internal/interfaces/tui/articlelist.go` — panneau droit, liste des articles

**Comportements :**
- Navigation j/k dans la liste active
- Tab pour switcher le focus FeedList ↔ ArticleList
- Sélection d'un feed charge ses articles dans le panneau droit
- `MarkRead` fire-and-forget en background à l'ouverture d'un article

Layout : `listW = w * 30 / 100` — toujours soustraire `style.GetHorizontalFrameSize()` avant de passer la largeur au contenu.

**Dépend de :** #0008
