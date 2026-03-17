---
id: 11
title: '[polish] États vides, badge non-lus, README & build'
status: open
type: task
priority: medium
author: mcp-agent
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Finaliser l'expérience et préparer la distribution.

**Fichiers :**
- `internal/interfaces/tui/app.go` — propagation correcte de `tea.WindowSizeMsg`
- `internal/interfaces/tui/feedlist.go` — badge `(N non-lus)` en style atténué
- `README.md` — usage, install, table des keybindings
- `Makefile` ou `.goreleaser.yaml` — build cross-platform

**Empty states :**
- Panneau gauche vide : `"Aucun flux. Appuie sur 'a' pour en ajouter un."`
- Panneau droit vide : `"Sélectionne un flux pour voir ses articles."`

**Vérification end-to-end :**
```bash
go build ./...
go test ./...
go run ./cmd/ratatosk add https://news.ycombinator.com/rss
go run ./cmd/ratatosk refresh
go run ./cmd/ratatosk
```

**Dépend de :** #0010
