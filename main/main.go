package main

import (
	"context"
	"flag"
	stdlog "log"
	"os/signal"
	"syscall"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var bindAddr, connectionStr string

func init() {
	// init logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.InterfaceMarshalFunc = json.Marshal
	stdlog.SetFlags(0)
	stdlog.SetOutput(&log.Logger)
	// init flags
	flag.StringVar(&bindAddr, "bind_addr", "0.0.0.0:8080", "bind address")
	flag.StringVar(&connectionStr, "connection_str",
		"host=localhost user=postgres password=postgres dbname=db_gql sslmode=disable connect_timeout=2",
		"connection string")
	flag.Parse()
}

func main() {
	if app, err := NewApplication(connectionStr); err != nil {
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

		if err := app.Listen(bindAddr); err != nil {
			log.Error().Err(err).Msg("start application")
		}
	}
}