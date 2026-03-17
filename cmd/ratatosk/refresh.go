package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/kvitrvn/ratatosk/internal/application"
)

func refreshCmd(svc **application.FeedService) *cobra.Command {
	return &cobra.Command{
		Use:   "refresh",
		Short: "Refresh all feeds",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			errs := (*svc).RefreshAll()
			if len(errs) > 0 {
				for _, err := range errs {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: %v\n", err)
				}
			}
			fmt.Println("Refresh done.")
			return nil
		},
	}
}
