package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kvitrvn/ratatosk/internal/application"
	"github.com/kvitrvn/ratatosk/internal/domain"
)

// --- in-memory stubs ---

type stubFeedRepo struct {
	feeds  map[int64]domain.Feed
	nextID int64
}

func newStubFeedRepo() *stubFeedRepo {
	return &stubFeedRepo{feeds: make(map[int64]domain.Feed), nextID: 1}
}

func (r *stubFeedRepo) Save(_ context.Context, f domain.Feed) (domain.Feed, error) {
	for _, existing := range r.feeds {
		if existing.URL == f.URL {
			existing.Title = f.Title
			r.feeds[existing.ID] = existing
			return existing, nil
		}
	}
	f.ID = r.nextID
	r.nextID++
	r.feeds[f.ID] = f
	return f, nil
}

func (r *stubFeedRepo) FindByID(_ context.Context, id int64) (domain.Feed, error) {
	f, ok := r.feeds[id]
	if !ok {
		return domain.Feed{}, errors.New("not found")
	}
	return f, nil
}

func (r *stubFeedRepo) FindAll(_ context.Context) ([]domain.Feed, error) {
	out := make([]domain.Feed, 0, len(r.feeds))
	for _, f := range r.feeds {
		out = append(out, f)
	}
	return out, nil
}

func (r *stubFeedRepo) Delete(_ context.Context, id int64) error {
	delete(r.feeds, id)
	return nil
}

type stubArticleRepo struct {
	articles map[int64]domain.Article
	nextID   int64
}

func newStubArticleRepo() *stubArticleRepo {
	return &stubArticleRepo{articles: make(map[int64]domain.Article), nextID: 1}
}

func (r *stubArticleRepo) SaveAll(_ context.Context, arts []domain.Article) error {
	for _, a := range arts {
		// deduplicate by feed_id + guid
		for _, existing := range r.articles {
			if existing.FeedID == a.FeedID && existing.GUID == a.GUID {
				goto next
			}
		}
		a.ID = r.nextID
		r.nextID++
		r.articles[a.ID] = a
	next:
	}
	return nil
}

func (r *stubArticleRepo) FindByFeedID(_ context.Context, feedID int64) ([]domain.Article, error) {
	var out []domain.Article
	for _, a := range r.articles {
		if a.FeedID == feedID {
			out = append(out, a)
		}
	}
	return out, nil
}

func (r *stubArticleRepo) MarkRead(_ context.Context, id int64) error {
	a, ok := r.articles[id]
	if !ok {
		return errors.New("article not found")
	}
	a.Read = true
	r.articles[id] = a
	return nil
}

type stubFetcher struct {
	result application.FetchedFeed
	err    error
}

func (f *stubFetcher) Fetch(_ string) (application.FetchedFeed, error) {
	return f.result, f.err
}

// --- tests ---

func newService(fetcher application.Fetcher) (*application.FeedService, *stubFeedRepo, *stubArticleRepo) {
	fr := newStubFeedRepo()
	ar := newStubArticleRepo()
	return application.NewFeedService(fr, ar, fetcher), fr, ar
}

func TestSubscribe_ValidURL(t *testing.T) {
	now := time.Now()
	svc, _, ar := newService(&stubFetcher{result: application.FetchedFeed{
		Title: "Test Feed",
		Articles: []application.FetchedArticle{
			{GUID: "g1", Title: "Article 1", PublishedAt: &now},
		},
	}})

	feed, err := svc.Subscribe("https://example.com/rss")
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	if feed.ID == 0 {
		t.Error("expected non-zero feed ID")
	}

	arts, _ := ar.FindByFeedID(context.Background(), feed.ID)
	if len(arts) != 1 {
		t.Errorf("expected 1 article, got %d", len(arts))
	}
}

func TestSubscribe_InvalidURL(t *testing.T) {
	svc, _, _ := newService(&stubFetcher{})
	if _, err := svc.Subscribe("not-a-url"); err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestSubscribe_FetchError(t *testing.T) {
	svc, fr, _ := newService(&stubFetcher{err: errors.New("network error")})
	_, err := svc.Subscribe("https://example.com/rss")
	// Feed should still be saved despite fetch error.
	if err == nil {
		t.Error("expected fetch error to be propagated")
	}
	feeds, _ := fr.FindAll(context.Background())
	if len(feeds) != 1 {
		t.Error("feed should be persisted even when fetch fails")
	}
}

func TestRefresh(t *testing.T) {
	svc, _, ar := newService(&stubFetcher{result: application.FetchedFeed{
		Articles: []application.FetchedArticle{
			{GUID: "g1", Title: "A1"},
			{GUID: "g2", Title: "A2"},
		},
	}})

	feed, _ := svc.Subscribe("https://example.com/rss")
	// Refresh again — dedup should keep only 2 articles.
	if err := svc.Refresh(feed.ID); err != nil {
		t.Fatalf("Refresh: %v", err)
	}

	arts, _ := ar.FindByFeedID(context.Background(), feed.ID)
	if len(arts) != 2 {
		t.Errorf("expected 2 articles after refresh, got %d", len(arts))
	}
}

func TestRefreshAll(t *testing.T) {
	fetcher := &stubFetcher{err: errors.New("timeout")}
	svc, fr, _ := newService(fetcher)

	fr.Save(context.Background(), domain.Feed{URL: "https://a.com/rss"})
	fr.Save(context.Background(), domain.Feed{URL: "https://b.com/rss"})

	errs := svc.RefreshAll()
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}
}

func TestMarkRead(t *testing.T) {
	svc, _, ar := newService(&stubFetcher{result: application.FetchedFeed{
		Articles: []application.FetchedArticle{{GUID: "g1", Title: "A"}},
	}})

	feed, _ := svc.Subscribe("https://example.com/rss")
	arts, _ := ar.FindByFeedID(context.Background(), feed.ID)
	if len(arts) == 0 {
		t.Fatal("no articles")
	}

	if err := svc.MarkRead(arts[0].ID); err != nil {
		t.Fatalf("MarkRead: %v", err)
	}

	arts, _ = ar.FindByFeedID(context.Background(), feed.ID)
	if !arts[0].Read {
		t.Error("article should be marked as read")
	}
}
