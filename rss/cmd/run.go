package cmd

import (
	"time"

	"github.com/alipourhabibi/exercises-journal/rss/config"
	"github.com/alipourhabibi/exercises-journal/rss/internal/core/services/logger"
	"github.com/alipourhabibi/exercises-journal/rss/internal/core/services/rss"
	"github.com/alipourhabibi/exercises-journal/rss/internal/core/services/server"
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
		timeout, err := time.ParseDuration(config.Conf.Http.Timeout)
		if err != nil {
			return err
		}
		serverService, err := server.New(
			server.WithHost(config.Conf.Http.Destination),
			server.WithTimetout(timeout),
			server.WithLogger(logger.Logger()),
		)
		if err != nil {
			return err
		}

		rssService, err := rss.New(
			rss.WithLogger(logger.Logger()),
			rss.WithFeeds("config/rss.yaml"),
			rss.WithInterval(config.Conf.Http.Interval),
			rss.WithRetInterval(config.Conf.Http.RetryInterval),
			rss.WithNewMemDB(),
			rss.WithServerService(serverService),
		)
		if err != nil {
			return err
		}

		err = rssService.Serve()
		if err != nil {
			return err
		}
		return nil
	},
}
