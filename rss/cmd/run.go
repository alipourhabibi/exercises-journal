package cmd

import (
	"github.com/alipourhabibi/exercises-journal/rss/config"
	"github.com/alipourhabibi/exercises-journal/rss/internal/core/services/logger"
	"github.com/alipourhabibi/exercises-journal/rss/internal/core/services/rss"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the rss server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := config.Conf.Load("config/config.yaml")
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, err := logger.New(
			logger.WithConfig(),
		)
		if err != nil {
			return err
		}
		rssService, err := rss.New(
			rss.WithLogger(logger.Logger()),
		)

		err = rssService.Serve()
		if err != nil {
			return err
		}
		return nil
	},
}
