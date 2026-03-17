# Ratatosk

<p align="center">
  <img src="docs/logo.png" alt="Ratatosk logo" width="180" />
</p>

**Ratatosk — terminal RSS/Atom feed reader written in Go.**

```
┌────────────────┐┌──────────────────────────────────────────┐
│ Feeds          ││ ● Linus steps down             2026-03-01│
│ HN (3)         ││ ○ Rust 2.0 released            2026-02-14│
│▶ Go Blog       ││▶● Go 1.27 release notes        2026-01-10│
│ LWN            ││                                          │
└────────────────┘└──────────────────────────────────────────┘
 a:add  r:refresh  tab:pane  enter:ouvrir  j/k:nav  q:quitter
```

## Install

```bash
go install github.com/kvitrvn/ratatosk/cmd/ratatosk@latest
```

Or build from source:

```bash
git clone https://github.com/kvitrvn/ratatosk
cd ratatosk
make build
```

## Usage

```bash
# Launch the TUI
ratatosk

# Add a feed from the command line
ratatosk add https://news.ycombinator.com/rss

# Refresh all feeds
ratatosk refresh
```

## Keyboard shortcuts

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `tab` | Switch pane (feeds → articles → detail) |
| `enter` | Open selected feed / article |
| `esc` | Go back to article list |
| `a` | Add a feed (overlay) |
| `r` | Refresh all feeds |
| `q` / `ctrl+c` | Quit |

## Configuration

The config file is created automatically on first run.

| Platform | Path |
|----------|------|
| Linux    | `~/.config/ratatosk/config.yaml` |
| macOS    | `~/Library/Application Support/ratatosk/config.yaml` |

## Dependencies

| Package | Role |
|---------|------|
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) | TUI event loop (Elm architecture) |
| [Lip Gloss](https://github.com/charmbracelet/lipgloss) | Terminal styling & layout |
| [Bubbles](https://github.com/charmbracelet/bubbles) | Viewport, textinput, spinner |
| [gofeed](https://github.com/mmcdole/gofeed) | RSS / Atom / JSON feed parsing |
| [modernc sqlite](https://gitlab.com/cznic/sqlite) | Pure-Go SQLite (no CGo) |
| [Cobra](https://github.com/spf13/cobra) + [Viper](https://github.com/spf13/viper) | CLI & config |
