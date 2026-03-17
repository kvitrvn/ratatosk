# ratatosk

Terminal RSS feed reader written in Go.

## Commands

```bash
# Build
go build ./...

# Run
go run ./cmd/ratatosk

# Test
go test ./...

# Vet
go vet ./...
```

## Key dependencies

| Package | Role |
|---|---|
| `github.com/charmbracelet/bubbletea` | TUI framework (Elm-architecture event loop) |
| `github.com/charmbracelet/lipgloss` | Terminal styling / layout |
| `github.com/charmbracelet/bubbles` | Reusable TUI components (list, viewport, spinner, …) |
| `github.com/mmcdole/gofeed` | RSS / Atom / JSON feed parsing |
| `modernc.org/sqlite` | SQLite driver (pure Go, no CGo) |
| `github.com/spf13/cobra` | CLI command structure |
| `github.com/spf13/viper` | Configuration (file + env vars) |

## Project layout

```
cmd/
  ratatosk/          # main package — wires Cobra root command and starts the app

internal/
  domain/            # Core business logic (no external dependencies)
                     #   entities, value objects, domain events, repository interfaces

  application/       # Use cases / application services
                     #   orchestrate domain objects, call repository interfaces

  infrastructure/    # Driven adapters
                     #   SQLite repositories, HTTP feed fetcher, gofeed adapter

  interfaces/        # Driving adapters
                     #   Bubble Tea TUI models, views, and Cmd helpers
```

### DDD layer rules

- `domain` imports nothing from the other internal layers.
- `application` imports `domain` only.
- `infrastructure` and `interfaces` import `domain` and `application`.
- Dependency direction always points inward toward `domain`.

## Skills

<!-- @skill go-charm-ui .claude/skills/go-charm-ui.md -->
Use the `go-charm-ui` skill whenever working on TUI code, Bubbletea models, Lipgloss layouts, pane navigation, key bindings, scrolling, overlays, or any terminal rendering in Go.
