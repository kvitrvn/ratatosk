---
id: 3
title: '[infra] Connexion SQLite & migrations'
status: in-progress
type: task
priority: medium
author: mcp-agent
milestone: M2 – SQLite persistence
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Mettre en place la connexion à la base et les migrations initiales.

**Fichiers :**
- `internal/infrastructure/db/sqlite.go` — `OpenDB(path)`, migrations, `PRAGMA foreign_keys`

**Schéma :**
```sql
CREATE TABLE IF NOT EXISTS feeds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    feed_id INTEGER NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    guid TEXT NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    link TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    published_at DATETIME,
    read INTEGER NOT NULL DEFAULT 0,
    UNIQUE(feed_id, guid)
);
```

Chemin DB via `os.UserConfigDir() + /ratatosk/ratatosk.db`

**Dépend de :** #0002
