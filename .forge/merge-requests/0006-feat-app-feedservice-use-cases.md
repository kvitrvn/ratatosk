---
id: 6
title: 'feat(app): FeedService — use cases'
status: merged
source_branch: feat/issue-0006-feed-service
target_branch: main
author: mcp-agent
linked_issues:
    - 6
approvals:
    - reviewer: kvitrvn
      approved_at: "2026-03-17"
approvals_required: 1
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Closes #0006

## Changements

- `internal/application/service.go` — `FeedService` : Subscribe, Refresh, RefreshAll, GetArticles, MarkRead
- `internal/application/service_test.go` — stubs in-memory, 6 tests

## Vérification

- `go test ./internal/application/...` ✓
- `go build ./...` ✓
- `go vet ./...` ✓
