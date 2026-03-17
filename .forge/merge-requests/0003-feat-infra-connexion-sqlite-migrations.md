---
id: 3
title: 'feat(infra): connexion SQLite & migrations'
status: merged
source_branch: feat/issue-0003-sqlite-connection
target_branch: main
author: mcp-agent
linked_issues:
    - 3
approvals:
    - reviewer: Kvitrvn
      approved_at: "2026-03-17"
approvals_required: 1
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Closes #0003

## Changements

- `internal/infrastructure/db/sqlite.go` — `OpenDB(path)`, migrations, `PRAGMA foreign_keys`
- `DefaultDBPath()` — chemin `~/.config/ratatosk/ratatosk.db`

## Vérification

- `go build ./...` ✓
- `go vet ./...` ✓
