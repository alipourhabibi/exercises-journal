package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "help",
	Short: "help command",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

func Execute() error {
	return rootCmd.Execute()
}
