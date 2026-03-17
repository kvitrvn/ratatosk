package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/kvitrvn/ratatosk/internal/application"
)

// AddFeedModel is the overlay form for subscribing to a new feed URL.
type AddFeedModel struct {
	svc       *application.FeedService
	textInput textinput.Model
}

func NewAddFeedModel(svc *application.FeedService) AddFeedModel {
	ti := textinput.New()
	ti.Placeholder = "https://example.com/feed.xml"
	ti.Focus()
	return AddFeedModel{svc: svc, textInput: ti}
}

func (m AddFeedModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AddFeedModel) Update(msg tea.Msg) (AddFeedModel, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && key.Matches(keyMsg, keys.Enter) {
		return m, subscribeCmd(m.svc, m.textInput.Value())
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m AddFeedModel) View() string {
	hint := "\n\nPress Enter to subscribe, Esc to cancel"
	return overlayStyle.Render("Add Feed\n\n" + m.textInput.View() + hint)
}

func subscribeCmd(svc *application.FeedService, rawURL string) tea.Cmd {
	return func() tea.Msg {
		feed, err := svc.Subscribe(rawURL)
		return feedSubscribedMsg{feed: feed, err: err}
	}
}
