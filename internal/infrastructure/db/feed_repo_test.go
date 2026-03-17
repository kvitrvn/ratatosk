package db_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/kvitrvn/ratatosk/internal/domain"
	"github.com/kvitrvn/ratatosk/internal/infrastructure/db"
)

func openTestDB(t *testing.T) *db.SQLiteFeedRepository {
	t.Helper()
	d, err := db.OpenDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { d.Close() })
	return db.NewSQLiteFeedRepository(d)
}

func TestFeedRepo_SaveAndFind(t *testing.T) {
	repo := openTestDB(t)
	ctx := context.Background()

	feed := domain.Feed{URL: "https://example.com/rss", Title: "Example", CreatedAt: time.Now().Truncate(time.Second)}
	saved, err := repo.Save(ctx, feed)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if saved.ID == 0 {
		t.Fatal("expected non-zero ID")
	}

	found, err := repo.FindByID(ctx, saved.ID)
	if err != nil {
		t.Fatalf("find by id: %v", err)
	}
	if found.URL != feed.URL {
		t.Errorf("URL = %q, want %q", found.URL, feed.URL)
	}
}

func TestFeedRepo_FindAll(t *testing.T) {
	repo := openTestDB(t)
	ctx := context.Background()

	for _, u := range []string{"https://a.com/rss", "https://b.com/rss"} {
		if _, err := repo.Save(ctx, domain.Feed{URL: u, CreatedAt: time.Now()}); err != nil {
			t.Fatalf("save: %v", err)
		}
	}

	feeds, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if len(feeds) != 2 {
		t.Errorf("got %d feeds, want 2", len(feeds))
	}
}

func TestFeedRepo_Delete(t *testing.T) {
	repo := openTestDB(t)
	ctx := context.Background()

	saved, err := repo.Save(ctx, domain.Feed{URL: "https://del.com/rss", CreatedAt: time.Now()})
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := repo.Delete(ctx, saved.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	feeds, _ := repo.FindAll(ctx)
	if len(feeds) != 0 {
		t.Errorf("expected 0 feeds after delete, got %d", len(feeds))
	}
}

func TestFeedRepo_SaveDuplicate(t *testing.T) {
	repo := openTestDB(t)
	ctx := context.Background()

	feed := domain.Feed{URL: "https://dup.com/rss", CreatedAt: time.Now()}
	if _, err := repo.Save(ctx, feed); err != nil {
		t.Fatalf("first save: %v", err)
	}
	if _, err := repo.Save(ctx, feed); err != nil {
		t.Fatalf("duplicate save should not error: %v", err)
	}
	feeds, _ := repo.FindAll(ctx)
	if len(feeds) != 1 {
		t.Errorf("expected 1 feed, got %d", len(feeds))
	}
}
