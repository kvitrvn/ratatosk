package application

import (
	"context"
	"fmt"

	"github.com/kvitrvn/ratatosk/internal/domain"
)

type FeedService struct {
	feeds    domain.FeedRepository
	articles domain.ArticleRepository
	fetcher  Fetcher
}

func NewFeedService(feeds domain.FeedRepository, articles domain.ArticleRepository, fetcher Fetcher) *FeedService {
	return &FeedService{feeds: feeds, articles: articles, fetcher: fetcher}
}

// Subscribe validates the URL, persists the feed, and performs the first fetch.
func (s *FeedService) Subscribe(rawURL string) (domain.Feed, error) {
	feed, err := domain.NewFeed(rawURL)
	if err != nil {
		return domain.Feed{}, err
	}

	saved, err := s.feeds.Save(context.Background(), feed)
	if err != nil {
		return domain.Feed{}, fmt.Errorf("save feed: %w", err)
	}

	if err := s.refresh(saved); err != nil {
		// Non-fatal: feed is saved, fetch failure is reported but not blocking.
		return saved, fmt.Errorf("initial fetch: %w", err)
	}
	return saved, nil
}

// Refresh fetches new articles for the given feed.
func (s *FeedService) Refresh(feedID int64) error {
	feed, err := s.feeds.FindByID(context.Background(), feedID)
	if err != nil {
		return fmt.Errorf("find feed %d: %w", feedID, err)
	}
	return s.refresh(feed)
}

// RefreshAll refreshes every feed and collects all errors.
func (s *FeedService) RefreshAll() []error {
	feeds, err := s.feeds.FindAll(context.Background())
	if err != nil {
		return []error{fmt.Errorf("list feeds: %w", err)}
	}

	var errs []error
	for _, f := range feeds {
		if err := s.refresh(f); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// GetArticles returns all articles for a feed.
func (s *FeedService) GetArticles(feedID int64) ([]domain.Article, error) {
	return s.articles.FindByFeedID(context.Background(), feedID)
}

// MarkRead marks an article as read.
func (s *FeedService) MarkRead(articleID int64) error {
	return s.articles.MarkRead(context.Background(), articleID)
}

// ListFeeds returns all subscribed feeds.
func (s *FeedService) ListFeeds() ([]domain.Feed, error) {
	return s.feeds.FindAll(context.Background())
}

func (s *FeedService) refresh(feed domain.Feed) error {
	fetched, err := s.fetcher.Fetch(feed.URL)
	if err != nil {
		return err
	}

	// Update title if it was empty.
	if feed.Title == "" && fetched.Title != "" {
		feed.Title = fetched.Title
		if _, err := s.feeds.Save(context.Background(), feed); err != nil {
			return fmt.Errorf("update feed title: %w", err)
		}
	}

	articles := make([]domain.Article, 0, len(fetched.Articles))
	for _, a := range fetched.Articles {
		articles = append(articles, domain.NewArticle(
			feed.ID, a.GUID, a.Title, a.Link, a.Description, a.PublishedAt,
		))
	}
	return s.articles.SaveAll(context.Background(), articles)
}
