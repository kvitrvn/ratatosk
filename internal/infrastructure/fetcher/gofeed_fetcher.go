package fetcher

import (
	"fmt"

	"github.com/mmcdole/gofeed"

	"github.com/kvitrvn/ratatosk/internal/application"
)

type GoFeedFetcher struct {
	parser *gofeed.Parser
}

func NewGoFeedFetcher() *GoFeedFetcher {
	return &GoFeedFetcher{parser: gofeed.NewParser()}
}

func (f *GoFeedFetcher) Fetch(url string) (application.FetchedFeed, error) {
	feed, err := f.parser.ParseURL(url)
	if err != nil {
		return application.FetchedFeed{}, fmt.Errorf("fetch %q: %w", url, err)
	}

	result := application.FetchedFeed{
		Title:    feed.Title,
		Articles: make([]application.FetchedArticle, 0, len(feed.Items)),
	}

	for _, item := range feed.Items {
		a := application.FetchedArticle{
			GUID:        item.GUID,
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
		}
		if item.PublishedParsed != nil {
			t := *item.PublishedParsed
			a.PublishedAt = &t
		}
		result.Articles = append(result.Articles, a)
	}

	return result, nil
}
