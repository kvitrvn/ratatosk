package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS feeds (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    url        TEXT     NOT NULL UNIQUE,
    title      TEXT     NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS articles (
    id          INTEGER  PRIMARY KEY AUTOINCREMENT,
    feed_id     INTEGER  NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    guid        TEXT     NOT NULL,
    title       TEXT     NOT NULL DEFAULT '',
    link        TEXT     NOT NULL DEFAULT '',
    description TEXT     NOT NULL DEFAULT '',
    published_at DATETIME,
    read        INTEGER  NOT NULL DEFAULT 0,
    UNIQUE(feed_id, guid)
);
`

func OpenDB(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return db, nil
}

func DefaultDBPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("user config dir: %w", err)
	}
	return filepath.Join(dir, "ratatosk", "ratatosk.db"), nil
}
