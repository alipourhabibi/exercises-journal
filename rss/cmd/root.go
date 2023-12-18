package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rss",
	Short: "rss cil",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)

}
