package cmd

import (
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

		rssService, err := rss.New(
			rss.WithLogger(logger.Logger()),
			rss.WithFeeds("config/rss.yaml"),
			rss.WithInterval(config.Conf.Http.Interval),
			rss.WithRetInterval(config.Conf.Http.RetryInterval),
			rss.WithNewRetryMemDB(),
			rss.WithServerService(serverService),
		)
		if err != nil {
			return err
		}

		ch := make(chan struct{})
		go rssService.Serve(ch)
		var state byte
		const (
			waitForSignal byte = iota
			reconfigureStatus
		)
		signalCh := make(chan os.Signal, 1)
		for {
			switch state {
			case reconfigureStatus:
				// makes new logger
				// set it to http
				err := config.Conf.Load(configFile)
				if err != nil {
					logger.Logger().Sugar().Info("can't read config file for reload; will use the previous logger", "error", err.Error())
					state = waitForSignal
					continue
				}
				logger, err := zlogger.New(
					zlogger.WithConfig(),
				)
				rssService.SetLogger(logger.Logger())
				/*
					rssService.SetRetDB()
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
					rssService.SetRetInterval(config.Conf.DB.RetryDBPath)

					state = waitForSignal
					ch <- struct{}{}
					ch = make(chan struct{})
					rssService.Serve(ch)
				*/
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
					return nil
				}

			}
		}
	},
}
