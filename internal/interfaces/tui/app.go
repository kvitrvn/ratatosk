package tui

import (
	"fmt"
	"os/exec"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/kvitrvn/ratatosk/internal/application"
)

// Focus identifies which pane is active.
type Focus int

const (
	FocusFeedList    Focus = iota // 0
	FocusArticleList              // 1
	FocusDetail                   // 2
)

// AppModel is the root Bubble Tea model.
type AppModel struct {
	svc         *application.FeedService
	focus       Focus
	feedList    FeedListModel
	articleList ArticleListModel
	detail      DetailModel
	overlay     *AddFeedModel
	spinner     spinner.Model
	refreshing  bool
	statusMsg   string
	width       int
	height      int
}

func NewApp(svc *application.FeedService) AppModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return AppModel{
		svc:     svc,
		focus:   FocusFeedList,
		detail:  NewDetailModel(),
		spinner: s,
	}
}

func (m AppModel) Init() tea.Cmd {
	return loadFeedsCmd(m.svc)
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Overlay: intercept key messages first; let feedSubscribedMsg fall through.
	if m.overlay != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			return m.handleOverlay(keyMsg)
		}
		if _, ok := msg.(feedSubscribedMsg); !ok {
			overlay, cmd := m.overlay.Update(msg)
			m.overlay = &overlay
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resize()

	case tea.KeyMsg:
		return m.handleKey(msg)

	case spinner.TickMsg:
		if m.refreshing {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case feedsLoadedMsg:
		m.feedList.SetFeeds(msg.feeds)
		if sel := m.feedList.Selected(); sel != nil {
			return m, loadArticlesCmd(m.svc, sel.ID)
		}

	case articlesLoadedMsg:
		m.articleList.SetArticles(msg.articles)
		unread := 0
		for _, a := range msg.articles {
			if !a.Read {
				unread++
			}
		}
		m.feedList.SetUnreadCount(msg.feedID, unread)

	case refreshDoneMsg:
		m.refreshing = false
		if len(msg.errs) > 0 {
			m.statusMsg = fmt.Sprintf("%d error(s) during refresh", len(msg.errs))
		} else {
			m.statusMsg = "Refreshed"
		}
		return m, loadFeedsCmd(m.svc)

	case feedSubscribedMsg:
		m.overlay = nil
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Error: %s", msg.err)
		} else {
			m.statusMsg = fmt.Sprintf("Subscribed to %s", msg.feed.Title)
		}
		return m, loadFeedsCmd(m.svc)

	case markReadMsg:
		if sel := m.feedList.Selected(); sel != nil {
			return m, loadArticlesCmd(m.svc, sel.ID)
		}

	case urlOpenedMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Impossible d'ouvrir : %s", msg.err)
		}
	}

	return m, nil
}

func (m AppModel) handleOverlay(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, keys.Esc) {
		m.overlay = nil
		return m, nil
	}
	overlay, cmd := m.overlay.Update(msg)
	m.overlay = &overlay
	return m, cmd
}

func (m AppModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, keys.Esc):
		if m.focus == FocusDetail {
			m.focus = FocusArticleList
		}

	case key.Matches(msg, keys.Tab):
		m.focus = (m.focus + 1) % 3
		if m.focus == FocusDetail {
			m.detail.SetArticle(m.articleList.Selected())
		}

	case key.Matches(msg, keys.AddFeed):
		overlay := NewAddFeedModel(m.svc)
		m.overlay = &overlay
		return m, m.overlay.Init()

	case key.Matches(msg, keys.Refresh):
		if !m.refreshing {
			m.refreshing = true
			m.statusMsg = "Refreshing…"
			return m, tea.Batch(refreshAllCmd(m.svc), m.spinner.Tick)
		}

	case key.Matches(msg, keys.Enter):
		switch m.focus {
		case FocusFeedList:
			if sel := m.feedList.Selected(); sel != nil {
				m.focus = FocusArticleList
				return m, loadArticlesCmd(m.svc, sel.ID)
			}
		case FocusArticleList:
			if sel := m.articleList.Selected(); sel != nil {
				m.focus = FocusDetail
				m.detail.SetArticle(sel)
				if !sel.Read {
					return m, markReadCmd(m.svc, sel.ID)
				}
			}
		}

	case key.Matches(msg, keys.Open):
		var link string
		switch m.focus {
		case FocusArticleList:
			if sel := m.articleList.Selected(); sel != nil {
				link = sel.Link
			}
		case FocusDetail:
			if m.detail.Article() != nil {
				link = m.detail.Article().Link
			}
		}
		if link != "" {
			return m, openURLCmd(link)
		}

	case key.Matches(msg, keys.Up), key.Matches(msg, keys.Down):
		switch m.focus {
		case FocusFeedList:
			m.feedList, _ = m.feedList.Update(msg)
			if sel := m.feedList.Selected(); sel != nil {
				return m, loadArticlesCmd(m.svc, sel.ID)
			}
		case FocusArticleList:
			m.articleList, _ = m.articleList.Update(msg)
			m.detail.SetArticle(m.articleList.Selected())
		case FocusDetail:
			var cmd tea.Cmd
			m.detail, cmd = m.detail.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m *AppModel) resize() {
	listW, articleW := dimensions(m.width)
	// Inner content height: total height − status bar (1) − top+bottom border (2).
	innerH := m.height - 3
	if innerH < 0 {
		innerH = 0
	}
	m.feedList.SetSize(listW, innerH)
	m.articleList.SetSize(articleW, innerH)
	m.detail.SetSize(articleW, innerH)
}

func (m AppModel) View() string {
	if m.width == 0 {
		return ""
	}

	feedView := m.feedList.View(m.focus == FocusFeedList)

	var rightView string
	if m.focus == FocusDetail {
		rightView = m.detail.View(true)
	} else {
		rightView = m.articleList.View(m.focus == FocusArticleList)
	}

	main := lipgloss.JoinHorizontal(lipgloss.Top, feedView, rightView)

	status := m.statusMsg
	if status == "" {
		status = "a:add  r:refresh  tab:pane  enter:ouvrir  o:navigateur  j/k:nav  q:quitter"
	}
	if m.refreshing {
		status = m.spinner.View() + " Rafraîchissement…"
	}
	statusBar := lipgloss.NewStyle().
		Width(m.width).
		Foreground(colorSubtle).
		Render(status)

	view := lipgloss.JoinVertical(lipgloss.Left, main, statusBar)

	if m.overlay != nil {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			m.overlay.View())
	}

	return view
}

// --- async commands ---

func loadFeedsCmd(svc *application.FeedService) tea.Cmd {
	return func() tea.Msg {
		feeds, _ := svc.ListFeeds()
		return feedsLoadedMsg{feeds: feeds}
	}
}

func loadArticlesCmd(svc *application.FeedService, feedID int64) tea.Cmd {
	return func() tea.Msg {
		articles, _ := svc.GetArticles(feedID)
		return articlesLoadedMsg{feedID: feedID, articles: articles}
	}
}

func refreshAllCmd(svc *application.FeedService) tea.Cmd {
	return func() tea.Msg {
		errs := svc.RefreshAll()
		return refreshDoneMsg{errs: errs}
	}
}

func openURLCmd(url string) tea.Cmd {
	return func() tea.Msg {
		err := exec.Command("xdg-open", url).Start()
		return urlOpenedMsg{err: err}
	}
}

func markReadCmd(svc *application.FeedService, articleID int64) tea.Cmd {
	return func() tea.Msg {
		_ = svc.MarkRead(articleID)
		return markReadMsg{}
	}
}

// truncate shortens s to at most maxLen runes, appending "…" if truncated.
func truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-1]) + "…"
}
