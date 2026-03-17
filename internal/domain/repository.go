package domain

import "context"

type FeedRepository interface {
	Save(ctx context.Context, feed Feed) (Feed, error)
	FindByID(ctx context.Context, id int64) (Feed, error)
	FindAll(ctx context.Context) ([]Feed, error)
	Delete(ctx context.Context, id int64) error
}

type ArticleRepository interface {
	SaveAll(ctx context.Context, articles []Article) error
	FindByFeedID(ctx context.Context, feedID int64) ([]Article, error)
	MarkRead(ctx context.Context, id int64) error
}
