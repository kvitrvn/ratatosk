---
id: 7
title: 'feat(cmd): config Viper & commandes Cobra'
status: merged
source_branch: feat/issue-0007-cli-config
target_branch: main
author: mcp-agent
linked_issues:
    - 7
approvals:
    - reviewer: kvitrvn
      approved_at: "2026-03-17"
approvals_required: 1
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Closes #0007

## Changements

- `internal/config/config.go` — `Config{DBPath, HTTPTimeout}`, defaults + `~/.config/ratatosk/config.yaml`
- `cmd/ratatosk/main.go` — root command avec DI dans `PersistentPreRunE`
- `cmd/ratatosk/add.go` — `ratatosk add <url>`
- `cmd/ratatosk/refresh.go` — `ratatosk refresh`

## Vérification

- `go build ./...` ✓
- `go vet ./...` ✓
