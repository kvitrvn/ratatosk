package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kvitrvn/ratatosk/internal/domain"
)

type SQLiteArticleRepository struct {
	db *sql.DB
}

func NewSQLiteArticleRepository(db *sql.DB) *SQLiteArticleRepository {
	return &SQLiteArticleRepository{db: db}
}

func (r *SQLiteArticleRepository) SaveAll(ctx context.Context, articles []domain.Article) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT OR IGNORE INTO articles
			(feed_id, guid, title, link, description, published_at, read)
		VALUES (?, ?, ?, ?, ?, ?, 0)`)
	if err != nil {
		return fmt.Errorf("prepare insert article: %w", err)
	}
	defer stmt.Close()

	for _, a := range articles {
		var publishedAt *string
		if a.PublishedAt != nil {
			s := a.PublishedAt.UTC().Format(time.RFC3339)
			publishedAt = &s
		}
		if _, err := stmt.ExecContext(ctx,
			a.FeedID, a.GUID, a.Title, a.Link, a.Description, publishedAt,
		); err != nil {
			return fmt.Errorf("insert article %q: %w", a.GUID, err)
		}
	}
	return tx.Commit()
}

func (r *SQLiteArticleRepository) FindByFeedID(ctx context.Context, feedID int64) ([]domain.Article, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, feed_id, guid, title, link, description, published_at, read
		FROM articles WHERE feed_id = ? ORDER BY published_at DESC, id DESC`, feedID)
	if err != nil {
		return nil, fmt.Errorf("find articles by feed: %w", err)
	}
	defer rows.Close()

	var articles []domain.Article
	for rows.Next() {
		var a domain.Article
		var publishedAt *string
		var read int
		if err := rows.Scan(
			&a.ID, &a.FeedID, &a.GUID, &a.Title, &a.Link,
			&a.Description, &publishedAt, &read,
		); err != nil {
			return nil, err
		}
		if publishedAt != nil {
			t, _ := time.Parse(time.RFC3339, *publishedAt)
			a.PublishedAt = &t
		}
		a.Read = read == 1
		articles = append(articles, a)
	}
	return articles, rows.Err()
}

func (r *SQLiteArticleRepository) MarkRead(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE articles SET read = 1 WHERE id = ?`, id)
	return err
}
