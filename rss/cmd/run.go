package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alipourhabibi/exercises-journal/rss/config"
	zlogger "github.com/alipourhabibi/exercises-journal/rss/internal/core/services/logger"
	"github.com/alipourhabibi/exercises-journal/rss/internal/core/services/rss"
	"github.com/alipourhabibi/exercises-journal/rss/internal/core/services/server"
	"github.com/spf13/cobra"
)

var configFile string = "config/config.yaml"

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the rss server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := config.Conf.Load(configFile)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, err := zlogger.New(
			zlogger.WithConfig(),
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
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		rssService, err := rss.New(
			rss.WithLogger(logger.Logger()),
			rss.WithFeeds(ctx, "config/rss.yaml"),
			rss.WithInterval(config.Conf.Http.Interval),
			rss.WithRetInterval(config.Conf.Http.RetryInterval),
			rss.WithNewRetryMemDB(ctx),
			rss.WithServerService(serverService),
		)
		if err != nil {
			cancel()
			return err
		}

		go rssService.Serve(ctx)
		var state byte
		const (
			waitForSignal byte = iota
			reconfigureStatus
		)
		signalCh := make(chan os.Signal, 1)
		for {
			switch state {
			case reconfigureStatus:
				cancel()

				ctx = context.Background()
				ctx, cancel = context.WithCancel(ctx)

				err := config.Conf.Load(configFile)
				if err != nil {
					logger.Logger().Sugar().Info("can't read config file for reload; will use the previous settings", "error", err.Error())
					state = waitForSignal
					continue
				}
				logger, err := zlogger.New(
					zlogger.WithConfig(),
				)
				if err != nil {
					logger.Logger().Sugar().Errorw("setting up new logger", "error", err, "status", "using the previous logger")
				} else {
					rssService.SetLogger(logger.Logger())
				}

				rssService.SetRetDB(ctx)
				serverService, err = server.New(
					server.WithHost(config.Conf.Http.Destination),
					server.WithTimetout(timeout),
					server.WithLogger(logger.Logger()),
				)
				if err != nil {
					logger.Logger().Sugar().Errorw("setting up new server", "error", err, "status", "using the previous server")
				} else {
					rssService.SetServer(serverService)
				}
				rssService.SetInterval(config.Conf.Http.Interval)
				rssService.SetRetInterval(config.Conf.Http.RetryInterval)

				go rssService.Serve(ctx)

				state = waitForSignal
			case waitForSignal:
				signal.Notify(signalCh,
					syscall.SIGHUP,
					syscall.SIGINT,
					syscall.SIGTERM,
					syscall.SIGQUIT)

				sig := <-signalCh
				logger.Logger().Sugar().Info("msg", "signal recieved", "signal", sig)
				switch sig {
				case syscall.SIGHUP:
					state = reconfigureStatus
				case syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM:
					cancel()
					return nil
				}

			}
		}
	},
}
