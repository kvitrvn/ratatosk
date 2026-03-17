package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/kvitrvn/ratatosk/internal/domain"
)

const (
	dotWidth  = 2
	dateWidth = 10
)

// ArticleListModel renders the article list pane.
type ArticleListModel struct {
	articles []domain.Article
	cursor   int
	offset   int
	width    int
	height   int
}

func (m *ArticleListModel) SetArticles(articles []domain.Article) {
	m.articles = articles
	m.cursor = 0
	m.offset = 0
}

func (m *ArticleListModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *ArticleListModel) Selected() *domain.Article {
	if len(m.articles) == 0 || m.cursor >= len(m.articles) {
		return nil
	}
	a := m.articles[m.cursor]
	return &a
}

func (m *ArticleListModel) UnreadCount() int {
	count := 0
	for _, a := range m.articles {
		if !a.Read {
			count++
		}
	}
	return count
}

func (m ArticleListModel) Update(msg tea.KeyMsg) (ArticleListModel, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		if m.cursor > 0 {
			m.cursor--
			m.clampOffset()
		}
	case key.Matches(msg, keys.Down):
		if m.cursor < len(m.articles)-1 {
			m.cursor++
			m.clampOffset()
		}
	}
	return m, nil
}

func (m ArticleListModel) View(focused bool) string {
	style := panelStyle
	if focused {
		style = focusedPanelStyle
	}
	style = style.Width(m.width).Height(m.height)

	titleWidth := m.width - dotWidth - dateWidth - 1
	if titleWidth < 0 {
		titleWidth = 0
	}

	end := m.offset + m.height
	if end > len(m.articles) {
		end = len(m.articles)
	}

	if len(m.articles) == 0 {
		hint := lipgloss.NewStyle().Foreground(colorSubtle).
			Render("Sélectionne un flux pour voir ses articles")
		return style.Render(hint)
	}

	var content string
	for i := m.offset; i < end; i++ {
		a := m.articles[i]
		dot := readDot
		if !a.Read {
			dot = unreadDot
		}

		title := truncate(a.Title, titleWidth)
		title = fmt.Sprintf("%-*s", titleWidth, title)

		date := ""
		if a.PublishedAt != nil {
			date = a.PublishedAt.Format("2006-01-02")
		}

		line := fmt.Sprintf("%s %s %s", dot, title, date)
		if i == m.cursor {
			line = selectedItemStyle.Width(m.width).Render(line)
		}
		content += line + "\n"
	}
	return style.Render(content)
}

func (m *ArticleListModel) clampOffset() {
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
