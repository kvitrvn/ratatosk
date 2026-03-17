---
id: 10
title: '[tui] Vue détail (viewport) & overlay AddFeed (textinput)'
status: open
type: task
priority: medium
author: mcp-agent
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Implémenter la lecture d'article et l'ajout de flux.

**Fichiers :**
- `internal/interfaces/tui/detail.go` — vue détail avec `bubbles/viewport` pour la description
- `internal/interfaces/tui/addfeed.go` — overlay "ajouter flux" avec `bubbles/textinput`

**Comportements :**
- Enter sur un article → FocusDetail, affiche titre + description dans le viewport
- Esc depuis Detail → retour ArticleList
- `a` → affiche l'overlay AddFeed par-dessus le layout
- Enter dans l'overlay → appelle `FeedService.Subscribe`, ferme l'overlay, recharge la feedlist
- Esc dans l'overlay → ferme sans action

**Dépend de :** #0009
