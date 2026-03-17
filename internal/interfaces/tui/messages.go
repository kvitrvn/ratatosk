package tui

import "github.com/kvitrvn/ratatosk/internal/domain"

type feedsLoadedMsg struct{ feeds []domain.Feed }
type articlesLoadedMsg struct {
	feedID   int64
	articles []domain.Article
}
type refreshDoneMsg struct{ errs []error }
type feedSubscribedMsg struct {
	feed domain.Feed
	err  error
}
type markReadMsg struct{}
