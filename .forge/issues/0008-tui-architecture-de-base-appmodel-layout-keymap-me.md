---
id: 8
title: '[tui] Architecture de base : AppModel, layout, KeyMap, messages'
status: open
type: task
priority: medium
author: mcp-agent
created_at: "2026-03-17"
updated_at: "2026-03-17"
---

Poser les fondations de la TUI Bubbletea.

**Fichiers :**
- `internal/interfaces/tui/app.go` — `AppModel` racine, gestion du focus, délégation aux sous-modèles
- `internal/interfaces/tui/layout.go` — `RenderLayout` : header, deux panneaux (30/70), status bar
- `internal/interfaces/tui/keys.go` — `KeyMap` complet
- `internal/interfaces/tui/messages.go` — tous les types `tea.Msg`

**Structure AppModel :**
```go
type AppModel struct {
    svc         *application.FeedService
    focus       Focus   // FocusFeedList | FocusArticleList | FocusDetail
    feedList    FeedListModel
    articleList ArticleListModel
    detail      DetailModel
    overlay     *AddFeedModel
    spinner     spinner.Model
    refreshing  bool
    statusMsg   string
    width, height int
}
```

**Keybindings :** j/k (nav), Tab (focus), Enter (ouvrir), r (refresh), a (overlay), Esc (fermer/retour), q/ctrl+c (quitter)

Refresh async : `tea.Cmd` → spinner → `refreshDoneMsg` → reload liste

**Dépend de :** #0006
