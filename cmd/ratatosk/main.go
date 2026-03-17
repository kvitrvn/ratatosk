package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/kvitrvn/ratatosk/internal/application"
	"github.com/kvitrvn/ratatosk/internal/config"
	"github.com/kvitrvn/ratatosk/internal/infrastructure/db"
	"github.com/kvitrvn/ratatosk/internal/infrastructure/fetcher"
	"github.com/kvitrvn/ratatosk/internal/interfaces/tui"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	var (
		svc   *application.FeedService
		sqlDB interface{ Close() error }
	)

	root := &cobra.Command{
		Use:   "ratatosk",
		Short: "Terminal RSS feed reader",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			database, err := db.OpenDB(cfg.DBPath)
			if err != nil {
				return fmt.Errorf("open db: %w", err)
			}
			sqlDB = database

			feedRepo := db.NewSQLiteFeedRepository(database)
			articleRepo := db.NewSQLiteArticleRepository(database)
			f := fetcher.NewGoFeedFetcher()
			svc = application.NewFeedService(feedRepo, articleRepo, f)
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if sqlDB != nil {
				return sqlDB.Close()
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			p := tea.NewProgram(tui.NewApp(svc), tea.WithAltScreen())
			_, err := p.Run()
			return err
		},
	}

	root.AddCommand(addCmd(&svc))
	root.AddCommand(refreshCmd(&svc))
	return root
}
