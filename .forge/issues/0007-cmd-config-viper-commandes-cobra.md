---
id: 7
title: '[cmd] Config Viper & commandes Cobra'
status: closed
type: task
priority: medium
author: mcp-agent
milestone: M4 – CLI & config
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Binaire exécutable avec config fichier et sous-commandes.

**Fichiers :**
- `internal/config/config.go` — `Config{DBPath, HTTPTimeout}`, defaults viper
- `cmd/ratatosk/main.go` — root command (lance TUI placeholder)
- `cmd/ratatosk/add.go` — `ratatosk add <url>`
- `cmd/ratatosk/refresh.go` — `ratatosk refresh`

**Notes :**
- Config dans `~/.config/ratatosk/config.yaml`
- DI explicite dans `PersistentPreRunE`, pas de framework DI

**Peut démarrer dès que #0006 est mergé (le package config n'a pas de dépendance sur la TUI)**
