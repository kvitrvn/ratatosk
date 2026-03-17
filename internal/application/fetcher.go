package application

import "time"

// FetchedFeed is the result of fetching a remote feed.
type FetchedFeed struct {
	Title    string
	Articles []FetchedArticle
}

// FetchedArticle represents a single item from a remote feed.
type FetchedArticle struct {
	GUID        string
	Title       string
	Link        string
	Description string
	PublishedAt *time.Time
}

// Fetcher fetches and parses a remote feed URL.
type Fetcher interface {
	Fetch(url string) (FetchedFeed, error)
}
