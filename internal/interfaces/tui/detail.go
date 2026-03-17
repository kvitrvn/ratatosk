package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/kvitrvn/ratatosk/internal/domain"
)

// DetailModel renders a scrollable article detail pane.
type DetailModel struct {
	article  *domain.Article
	viewport viewport.Model
	width    int
	height   int
}

func NewDetailModel() DetailModel {
	return DetailModel{viewport: viewport.New(0, 0)}
}

func (m *DetailModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.viewport.Width = w - frameSize
	m.viewport.Height = h - 3 // title + date + separator
	if m.article != nil {
		m.viewport.SetContent(m.renderContent())
	}
}

func (m *DetailModel) SetArticle(a *domain.Article) {
	m.article = a
	if a == nil {
		m.viewport.SetContent("")
		return
	}
	m.viewport.SetContent(m.renderContent())
	m.viewport.GotoTop()
}

func (m *DetailModel) renderContent() string {
	if m.article == nil || m.article.Description == "" {
		return "(no content)"
	}
	return htmlToText(m.article.Description, m.viewport.Width)
}

func (m DetailModel) Article() *domain.Article { return m.article }

func (m DetailModel) Update(msg tea.KeyMsg) (DetailModel, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m DetailModel) View(focused bool) string {
	style := panelStyle
	if focused {
		style = focusedPanelStyle
	}
	style = style.Width(m.width).Height(m.height)

	if m.article == nil {
		return style.Render("")
	}

	boldStyle := lipgloss.NewStyle().Bold(true)
	metaStyle := lipgloss.NewStyle().Foreground(colorSubtle)
	sepStyle := lipgloss.NewStyle().Foreground(colorSubtle)

	title := boldStyle.Render(truncate(m.article.Title, m.width))
	date := ""
	if m.article.PublishedAt != nil {
		date = m.article.PublishedAt.Format("2006-01-02")
	}
	meta := metaStyle.Render(date)
	sep := sepStyle.Render("───")

	content := title + "\n" + meta + "\n" + sep + "\n" + m.viewport.View()
	return style.Render(content)
}
