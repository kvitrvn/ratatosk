---
id: 2
title: 'feat(domain): interfaces FeedRepository & ArticleRepository'
status: merged
source_branch: feat/issue-0002-repository-interfaces
target_branch: main
author: mcp-agent
approvals:
    - reviewer: Kvitrvn
      approved_at: "2026-03-17"
approvals_required: 1
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Closes #0002

## Changements

- `internal/domain/repository.go` — interfaces `FeedRepository` et `ArticleRepository`
- `cmd/ratatosk/main.go` — placeholder minimal (`go build ./...` passe)

## Vérification

- `go build ./...` ✓
- `go test ./internal/domain/...` ✓
- `go vet ./internal/domain/...` ✓
