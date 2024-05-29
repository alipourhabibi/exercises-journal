package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alipourhabibi/exercises-journal/echo/config"
	"github.com/alipourhabibi/exercises-journal/echo/internal/handlers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run http server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := config.Load("config/config.yaml")
		if err != nil {
			return err
		}
		level, ok := config.MapLevel[strings.ToUpper(config.Confs.Logger.Level)]
		if !ok {
			level = slog.LevelError
		}
		l := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: config.Confs.Logger.AddSource,
			Level:     level,
		}))
		slog.SetDefault(l)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		server := handlers.New()
		var err error
		go func() {
			err = server.Start(ctx)
			if err != nil {
				return
			}
		}()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

		for {
			sig := <-sigs
			switch sig {
			case syscall.SIGHUP:
				// TODO make function for reloading the configs
				slog.Info("Received SIGHUP, reloading configuration...")

				// reload the config
				err := config.Load("config/config.yaml")
				if err != nil {
					return err
				}

				// reload the logger with new config
				level, ok := config.MapLevel[strings.ToUpper(config.Confs.Logger.Level)]
				if !ok {
					level = slog.LevelError
				}
				l := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
					AddSource: config.Confs.Logger.AddSource,
					Level:     level,
				}))
				slog.SetDefault(l)

				// shutdown the prev server
				// continue with the the last config if failed
				err = server.Shutdown(ctx)
				if err != nil {
					slog.Error("could not start new server", "error", err)
					continue
				}
				server = handlers.New()
				go func() {
					if err := server.Start(ctx); err != nil {
						slog.Error("could not start new server", "error", err)
						// TODO it return an error but it does the job; think about it later
						// os.Exit(1)
					}
				}()

			case syscall.SIGINT, syscall.SIGTERM:
				slog.Info("Received SIGINT/SIGTERM, shutting down...")
				if err := server.Shutdown(ctx); err != nil {
					slog.Error("could not gracefully shut down server", "error", err)
				}
				os.Exit(0)
			}
		}
	},
}
