package domain

import "time"

type Article struct {
	ID          int64
	FeedID      int64
	GUID        string
	Title       string
	Link        string
	Description string
	PublishedAt *time.Time
	Read        bool
}

func NewArticle(feedID int64, guid, title, link, description string, publishedAt *time.Time) Article {
	return Article{
		FeedID:      feedID,
		GUID:        guid,
		Title:       title,
		Link:        link,
		Description: description,
		PublishedAt: publishedAt,
	}
}
