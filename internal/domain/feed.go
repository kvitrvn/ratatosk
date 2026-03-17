package domain

import (
	"fmt"
	"net/url"
	"time"
)

type Feed struct {
	ID        int64
	URL       string
	Title     string
	CreatedAt time.Time
}

func NewFeed(rawURL string) (Feed, error) {
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return Feed{}, fmt.Errorf("invalid feed URL: %q", rawURL)
	}
	return Feed{
		URL:       rawURL,
		CreatedAt: time.Now(),
	}, nil
}
