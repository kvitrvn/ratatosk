package db_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/kvitrvn/ratatosk/internal/domain"
	"github.com/kvitrvn/ratatosk/internal/infrastructure/db"
)

func openTestDBFull(t *testing.T) (*db.SQLiteFeedRepository, *db.SQLiteArticleRepository) {
	t.Helper()
	d, err := db.OpenDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { d.Close() })
	return db.NewSQLiteFeedRepository(d), db.NewSQLiteArticleRepository(d)
}

func TestArticleRepo_SaveAllAndFind(t *testing.T) {
	feedRepo, artRepo := openTestDBFull(t)
	ctx := context.Background()

	feed, err := feedRepo.Save(ctx, domain.Feed{URL: "https://example.com/rss", CreatedAt: time.Now()})
	if err != nil {
		t.Fatalf("save feed: %v", err)
	}

	now := time.Now().Truncate(time.Second)
	articles := []domain.Article{
		{FeedID: feed.ID, GUID: "guid-1", Title: "First", Link: "https://example.com/1", PublishedAt: &now},
		{FeedID: feed.ID, GUID: "guid-2", Title: "Second", Link: "https://example.com/2"},
	}
	if err := artRepo.SaveAll(ctx, articles); err != nil {
		t.Fatalf("save all: %v", err)
	}

	found, err := artRepo.FindByFeedID(ctx, feed.ID)
	if err != nil {
		t.Fatalf("find by feed: %v", err)
	}
	if len(found) != 2 {
		t.Errorf("got %d articles, want 2", len(found))
	}
}

func TestArticleRepo_SaveAllDeduplication(t *testing.T) {
	feedRepo, artRepo := openTestDBFull(t)
	ctx := context.Background()

	feed, _ := feedRepo.Save(ctx, domain.Feed{URL: "https://dup.com/rss", CreatedAt: time.Now()})

	a := domain.Article{FeedID: feed.ID, GUID: "same-guid", Title: "Title"}
	if err := artRepo.SaveAll(ctx, []domain.Article{a, a}); err != nil {
		t.Fatalf("save all: %v", err)
	}
	// Call again — should still deduplicate
	if err := artRepo.SaveAll(ctx, []domain.Article{a}); err != nil {
		t.Fatalf("second save: %v", err)
	}

	found, _ := artRepo.FindByFeedID(ctx, feed.ID)
	if len(found) != 1 {
		t.Errorf("expected 1 article after dedup, got %d", len(found))
	}
}

func TestArticleRepo_MarkRead(t *testing.T) {
	feedRepo, artRepo := openTestDBFull(t)
	ctx := context.Background()

	feed, _ := feedRepo.Save(ctx, domain.Feed{URL: "https://read.com/rss", CreatedAt: time.Now()})
	artRepo.SaveAll(ctx, []domain.Article{{FeedID: feed.ID, GUID: "g1", Title: "T"}})

	articles, _ := artRepo.FindByFeedID(ctx, feed.ID)
	if len(articles) == 0 {
		t.Fatal("no articles found")
	}
	if err := artRepo.MarkRead(ctx, articles[0].ID); err != nil {
		t.Fatalf("mark read: %v", err)
	}

	articles, _ = artRepo.FindByFeedID(ctx, feed.ID)
	if !articles[0].Read {
		t.Error("article should be marked as read")
	}
}
