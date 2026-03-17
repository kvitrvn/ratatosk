---
name: go-charm-tui
description: >
  Expert guide for building terminal UIs in Go using the Charm stack
  (Bubbletea, Lipgloss, Glamour). Use this skill whenever working on TUI
  code, Bubbletea models, Lipgloss layouts, pane navigation, key bindings,
  scrolling, overlays, or any terminal rendering in Go — even if the user
  just says "bubbletea", "TUI", "lipgloss", "terminal interface", "lazygit-style",
  or asks to build/fix/extend a view, model, or layout in the forge project.
---

# go-charm-tui

Expert patterns for the Charm TUI stack in the context of the Forge project:
a two-pane TUI (list + detail) with tabs, filters, overlays, and Markdown rendering.

## Architecture: Elm in Go

Bubbletea follows the Elm architecture. Every component is a `tea.Model`:

```go
type Model interface {
    Init() tea.Cmd          // side-effects on startup
    Update(tea.Msg) (tea.Model, tea.Cmd)  // pure state transition
    View() string           // pure render
}
```

**Critical rule**: `Update()` must be pure. Never do file I/O, git calls, or
blocking operations directly inside `Update()`. Wrap them in `tea.Cmd`.

---

## Project Model Structure

The forge TUI uses a root model that owns sub-models:

```go
// tui/app.go
type AppModel struct {
    activeTab  Tab
    issues     views.IssuesModel
    mrs        views.MRsModel
    width      int
    height     int
    statusMsg  string
}

type Tab int
const (
    TabIssues Tab = iota
    TabMRs
)
```

Delegate to sub-models in `Update()`:

```go
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            m.activeTab = (m.activeTab + 1) % 2
            return m, nil
        case "q", "ctrl+c":
            return m, tea.Quit
        }
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        // propagate to sub-models
    }
    // delegate to active sub-model
    switch m.activeTab {
    case TabIssues:
        updated, cmd := m.issues.Update(msg)
        m.issues = updated.(views.IssuesModel)
        return m, cmd
    }
    return m, nil
}
```

---

## Two-Pane Layout with Lipgloss

The forge layout: 30% list | 70% detail, fixed header and footer.

```go
// tui/layout.go
var (
    listPaneStyle = lipgloss.NewStyle().
        Width(30).
        BorderRight(true).
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("240"))

    detailPaneStyle = lipgloss.NewStyle().
        PaddingLeft(2)

    headerStyle = lipgloss.NewStyle().
        Bold(true).
        Background(lipgloss.Color("62")).
        Foreground(lipgloss.Color("230")).
        Padding(0, 1)

    statusBarStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("240")).
        Padding(0, 1)
)

func RenderLayout(header, list, detail, status string, w, h int) string {
    listW := w * 30 / 100
    detailW := w - listW - 1 // -1 for border

    listPane := listPaneStyle.Width(listW).Height(h - 3).Render(list)
    detailPane := detailPaneStyle.Width(detailW).Height(h - 3).Render(detail)

    body := lipgloss.JoinHorizontal(lipgloss.Top, listPane, detailPane)
    return lipgloss.JoinVertical(lipgloss.Left,
        headerStyle.Width(w).Render(header),
        body,
        statusBarStyle.Width(w).Render(status),
    )
}
```

**Width propagation**: always pass `tea.WindowSizeMsg` down to sub-models and
re-compute widths. Never hardcode terminal dimensions.

---

## List Navigation

```go
type IssuesModel struct {
    items    []forge.Issue
    cursor   int
    offset   int       // for scrolling
    height   int       // visible rows
    focused  bool      // list vs detail focus
}

func (m IssuesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "j", "down":
            if m.cursor < len(m.items)-1 {
                m.cursor++
                // scroll if cursor goes below visible area
                if m.cursor >= m.offset+m.height {
                    m.offset++
                }
            }
        case "k", "up":
            if m.cursor > 0 {
                m.cursor--
                if m.cursor < m.offset {
                    m.offset--
                }
            }
        }
    }
    return m, nil
}

func (m IssuesModel) View() string {
    var b strings.Builder
    visible := m.items[m.offset:min(m.offset+m.height, len(m.items))]
    for i, item := range visible {
        abs := i + m.offset
        row := renderIssueRow(item)
        if abs == m.cursor {
            row = selectedRowStyle.Render(row)
        }
        b.WriteString(row + "\n")
    }
    return b.String()
}
```

---

## Async File I/O with tea.Cmd

Never read files in `Update()`. Use commands:

```go
// Define the message type
type issuesLoadedMsg struct {
    issues []forge.Issue
    err    error
}

// Define the command (returns a tea.Cmd)
func loadIssuesCmd(path string) tea.Cmd {
    return func() tea.Msg {
        issues, err := forge.LoadIssues(path)
        return issuesLoadedMsg{issues: issues, err: err}
    }
}

// Handle in Update()
case issuesLoadedMsg:
    if msg.err != nil {
        m.statusMsg = "Error: " + msg.err.Error()
        return m, nil
    }
    m.items = msg.issues
    return m, nil
```

Start it from `Init()` or in response to a key:

```go
func (m AppModel) Init() tea.Cmd {
    return loadIssuesCmd(".forge/issues")
}
```

---

## Markdown Rendering with Glamour

```go
import "github.com/charmbracelet/glamour"

func renderMarkdown(content string, width int) string {
    r, err := glamour.NewTermRenderer(
        glamour.WithAutoStyle(),
        glamour.WithWordWrap(width),
    )
    if err != nil {
        return content // fallback to raw
    }
    out, err := r.Render(content)
    if err != nil {
        return content
    }
    return out
}
```

Cache the rendered output — re-rendering on every `View()` call is expensive.
Store rendered string in the model and invalidate when selection changes.

---

## Overlay / Modal Pattern

For filter prompts, help screens, confirmations:

```go
type AppModel struct {
    // ...
    overlay     *OverlayModel  // nil = no overlay
}

type OverlayModel struct {
    kind    OverlayKind  // Help, Filter, Confirm
    input   textinput.Model
}

// In View(): render overlay on top of base layout
func (m AppModel) View() string {
    base := m.renderBase()
    if m.overlay != nil {
        return renderOverlay(base, m.overlay.View(), m.width, m.height)
    }
    return base
}

func renderOverlay(base, content string, w, h int) string {
    box := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1, 2).
        Width(w/2).
        Render(content)
    // center the box over the base
    return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, box,
        lipgloss.WithWhitespaceBackground(lipgloss.Color("0")))
}
```

---

## Status Bar & Error Display

Never crash on errors — surface them in the status bar:

```go
// Timed status messages
type clearStatusMsg struct{}

func clearStatusAfter(d time.Duration) tea.Cmd {
    return tea.Tick(d, func(t time.Time) tea.Msg {
        return clearStatusMsg{}
    })
}

// In Update():
case clearStatusMsg:
    m.statusMsg = ""

// When an error occurs:
m.statusMsg = "✗ Could not load issues: " + err.Error()
return m, clearStatusAfter(4 * time.Second)
```

---

## Key Bindings with help.KeyMap

```go
// tui/keys.go
import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
    Up     key.Binding
    Down   key.Binding
    New    key.Binding
    Edit   key.Binding
    Close  key.Binding
    Filter key.Binding
    Search key.Binding
    Tab    key.Binding
    Help   key.Binding
    Quit   key.Binding
}

var DefaultKeys = KeyMap{
    Up:     key.NewBinding(key.WithKeys("k", "up"),     key.WithHelp("k/↑", "up")),
    Down:   key.NewBinding(key.WithKeys("j", "down"),   key.WithHelp("j/↓", "down")),
    New:    key.NewBinding(key.WithKeys("n"),            key.WithHelp("n", "new")),
    Edit:   key.NewBinding(key.WithKeys("e"),            key.WithHelp("e", "edit")),
    Close:  key.NewBinding(key.WithKeys("c"),            key.WithHelp("c", "close")),
    Filter: key.NewBinding(key.WithKeys("f"),            key.WithHelp("f", "filter")),
    Search: key.NewBinding(key.WithKeys("/"),            key.WithHelp("/", "search")),
    Tab:    key.NewBinding(key.WithKeys("tab"),          key.WithHelp("tab", "switch")),
    Help:   key.NewBinding(key.WithKeys("?"),            key.WithHelp("?", "help")),
    Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}
```

---

## Opening $EDITOR for new/edit

```go
type editorFinishedMsg struct{ file string; err error }

func openEditorCmd(path string) tea.Cmd {
    editor := os.Getenv("EDITOR")
    if editor == "" {
        editor = "vi"
    }
    c := exec.Command(editor, path)
    return tea.ExecProcess(c, func(err error) tea.Msg {
        return editorFinishedMsg{file: path, err: err}
    })
}
```

`tea.ExecProcess` suspends the TUI, hands control to the editor, then resumes.
Always use this — never `exec.Command(...).Run()` directly.

---

## Common Pitfalls

- **Lipgloss width includes padding/border**: account for this when computing
  content width: `contentW := paneW - style.GetHorizontalFrameSize()`
- **View() called every frame**: avoid heavy computation. Pre-render and cache.
- **tea.Cmd returns nil**: returning `nil` as a Cmd is valid (no-op), but
  `tea.Batch(nil, someCmd)` panics — always filter nils before batching.
- **String truncation for list rows**: use `lipgloss.NewStyle().MaxWidth(n).Render(s)`
  or `runewidth.Truncate()` for multi-byte safe truncation.

---

## Reference Files

- `references/lipgloss-layout.md` — advanced layout patterns, responsive sizing
- `references/bubbles.md` — pre-built components: textinput, viewport, spinner, table
