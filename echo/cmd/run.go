package cmd

import (
	"github.com/alipourhabibi/exercises-journal/echo/internal/handlers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run http server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return handlers.Launch()
	},
}
