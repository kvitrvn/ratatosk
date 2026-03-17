---
id: 4
title: 'feat(infra): repositories SQLite feed & article'
status: merged
source_branch: feat/issue-0004-sqlite-repositories
target_branch: main
author: mcp-agent
linked_issues:
    - 4
approvals:
    - reviewer: Kvitrvn
      approved_at: "2026-03-17"
approvals_required: 1
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Closes #0004

## Changements

- `feed_repo.go` — `SQLiteFeedRepository` : Save (upsert), FindByID, FindAll, Delete
- `article_repo.go` — `SQLiteArticleRepository` : SaveAll (INSERT OR IGNORE), FindByFeedID, MarkRead
- `feed_repo_test.go` — tests Save, FindAll, Delete, déduplication
- `article_repo_test.go` — tests SaveAll, déduplication par GUID, MarkRead

## Vérification

- `go test ./internal/infrastructure/db/...` ✓
- `go build ./...` ✓
- `go vet ./...` ✓
