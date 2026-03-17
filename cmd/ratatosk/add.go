package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/kvitrvn/ratatosk/internal/application"
)

func addCmd(svc **application.FeedService) *cobra.Command {
	return &cobra.Command{
		Use:   "add <url>",
		Short: "Subscribe to a feed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			feed, err := (*svc).Subscribe(args[0])
			if err != nil {
				return err
			}
			fmt.Printf("Subscribed to %q (id=%d)\n", feed.URL, feed.ID)
			return nil
		},
	}
}
