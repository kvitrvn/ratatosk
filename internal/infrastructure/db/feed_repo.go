package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kvitrvn/ratatosk/internal/domain"
)

type SQLiteFeedRepository struct {
	db *sql.DB
}

func NewSQLiteFeedRepository(db *sql.DB) *SQLiteFeedRepository {
	return &SQLiteFeedRepository{db: db}
}

func (r *SQLiteFeedRepository) Save(ctx context.Context, feed domain.Feed) (domain.Feed, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO feeds (url, title, created_at) VALUES (?, ?, ?)
		 ON CONFLICT(url) DO UPDATE SET title = excluded.title`,
		feed.URL, feed.Title, feed.CreatedAt.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return domain.Feed{}, fmt.Errorf("save feed: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Feed{}, fmt.Errorf("last insert id: %w", err)
	}
	if id == 0 {
		// ON CONFLICT path — fetch existing row
		return r.findByURL(ctx, feed.URL)
	}
	feed.ID = id
	return feed, nil
}

func (r *SQLiteFeedRepository) FindByID(ctx context.Context, id int64) (domain.Feed, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, url, title, created_at FROM feeds WHERE id = ?`, id)
	return scanFeed(row)
}

func (r *SQLiteFeedRepository) FindAll(ctx context.Context) ([]domain.Feed, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, url, title, created_at FROM feeds ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("find all feeds: %w", err)
	}
	defer rows.Close()

	var feeds []domain.Feed
	for rows.Next() {
		var f domain.Feed
		var createdAt string
		if err := rows.Scan(&f.ID, &f.URL, &f.Title, &createdAt); err != nil {
			return nil, err
		}
		f.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		feeds = append(feeds, f)
	}
	return feeds, rows.Err()
}

func (r *SQLiteFeedRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM feeds WHERE id = ?`, id)
	return err
}

func (r *SQLiteFeedRepository) findByURL(ctx context.Context, url string) (domain.Feed, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, url, title, created_at FROM feeds WHERE url = ?`, url)
	return scanFeed(row)
}

func scanFeed(row *sql.Row) (domain.Feed, error) {
	var f domain.Feed
	var createdAt string
	if err := row.Scan(&f.ID, &f.URL, &f.Title, &createdAt); err != nil {
		return domain.Feed{}, fmt.Errorf("scan feed: %w", err)
	}
	f.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return f, nil
}
