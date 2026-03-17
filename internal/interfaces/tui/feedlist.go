package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/kvitrvn/ratatosk/internal/domain"
)

// FeedListModel renders the left-hand feed list pane.
type FeedListModel struct {
	feeds        []domain.Feed
	unreadCounts map[int64]int
	cursor       int
	offset       int
	width        int
	height       int
}

func (m *FeedListModel) SetFeeds(feeds []domain.Feed) {
	m.feeds = feeds
	if m.cursor >= len(feeds) {
		m.cursor = max(0, len(feeds)-1)
	}
	m.clampOffset()
}

func (m *FeedListModel) SetUnreadCount(feedID int64, count int) {
	if m.unreadCounts == nil {
		m.unreadCounts = make(map[int64]int)
	}
	m.unreadCounts[feedID] = count
}

func (m *FeedListModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.clampOffset()
}

func (m *FeedListModel) Selected() *domain.Feed {
	if len(m.feeds) == 0 || m.cursor >= len(m.feeds) {
		return nil
	}
	f := m.feeds[m.cursor]
	return &f
}

func (m FeedListModel) Update(msg tea.KeyMsg) (FeedListModel, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		if m.cursor > 0 {
			m.cursor--
			m.clampOffset()
		}
	case key.Matches(msg, keys.Down):
		if m.cursor < len(m.feeds)-1 {
			m.cursor++
			m.clampOffset()
		}
	}
	return m, nil
}

func (m FeedListModel) View(focused bool) string {
	style := panelStyle
	if focused {
		style = focusedPanelStyle
	}
	style = style.Width(m.width).Height(m.height)

	if len(m.feeds) == 0 {
		hint := lipgloss.NewStyle().Foreground(colorSubtle).Render("Aucun flux — appuie sur 'a'")
		return style.Render(hint)
	}

	subtleStyle := lipgloss.NewStyle().Foreground(colorSubtle)

	end := m.offset + m.height
	if end > len(m.feeds) {
		end = len(m.feeds)
	}

	var content string
	for i := m.offset; i < end; i++ {
		f := m.feeds[i]
		title := f.Title
		if title == "" {
			title = f.URL
		}

		badge := ""
		if n := m.unreadCounts[f.ID]; n > 0 {
			badge = subtleStyle.Render(fmt.Sprintf(" (%d)", n))
		}

		maxTitle := m.width - 1 - len([]rune(badge)) // rough visual width of badge
		if maxTitle < 0 {
			maxTitle = 0
		}
		line := " " + truncate(title, maxTitle) + badge
		if i == m.cursor {
			line = selectedItemStyle.Width(m.width).Render(line)
		}
		content += line + "\n"
	}
	return style.Render(content)
}

func (m *FeedListModel) clampOffset() {
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.height > 0 && m.cursor >= m.offset+m.height {
		m.offset = m.cursor - m.height + 1
	}
	if m.offset < 0 {
		m.offset = 0
	}
}
