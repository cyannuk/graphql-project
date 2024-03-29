package main

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"graphql-project/config"
	"graphql-project/core"
)

func zeroLogLevel(cfg *config.Config) zerolog.Level {
	if level, err := zerolog.ParseLevel(cfg.LogLevel()); err == nil {
		return level
	} else {
		return zerolog.InfoLevel
	}
}

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.InterfaceMarshalFunc = json.Marshal
	log.SetFlags(0)
	log.SetOutput(core.LoggerWriter())
}

func SetLogLevel(cfg *config.Config) {
	zerolog.SetGlobalLevel(zeroLogLevel(cfg))
}

func FiberLogConfig() fiberzerolog.Config {
	return fiberzerolog.Config{
		Logger: &zlog.Logger,
		Levels: []zerolog.Level{zerolog.ErrorLevel, zerolog.WarnLevel, zerolog.DebugLevel},
	}
}
