package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"graphql-project/config"
)

func init() {
	InitLogger()
}

func main() {
	var cfg config.Config
	if err := cfg.Load(); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Println(err.Error())
		} else {
			log.Error().Err(err).Msg("config load")
		}
		return
	}
	SetLogLevel(&cfg)

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
