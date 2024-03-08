package main

import (
	"context"
	stdlog "log"
	"os/signal"
	"syscall"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"graphql-project/config"
)

func init() {
	// init logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.InterfaceMarshalFunc = json.Marshal
	stdlog.SetFlags(0)
	stdlog.SetOutput(&log.Logger)
}

func main() {
	// TODO mutations, dbmate, integration tests, go-bin/go-embedd
	var cfg config.Config
	if err := cfg.Load(); err != nil {
		log.Error().Err(err).Msg("config load")
		return
	}
	if app, err := NewApplication(&cfg); err != nil {
		log.Error().Err(err).Msg("init application")
	} else {
		go func() {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
			_ = <-ctx.Done()
			stop()
			if err := app.Shutdown(); err != nil {
				log.Error().Err(err).Msg("shutdown application")
			}
		}()

		if err := app.Start(); err != nil {
			log.Error().Err(err).Msg("start application")
		}
	}
}
