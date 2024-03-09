package config

import (
	"net/netip"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

//go:generate go run gen.go

type Config struct {
	bindAddr         netip.Addr    `cmd:"bind-addr" env:"ADDRESS" desc:"bind host address"`
	port             uint16        `cmd:"port" env:"PORT" desc:"listen port"`
	jwtSecret        []byte        `cmd:"jwt-secret" env:"JWT_SECRET" desc:"base64 encoded JWT secret"`
	jwtExpiration    time.Duration `cmd:"jwt-expiration" env:"JWT_EXPIRATION" desc:"JWT expiration time in hours"`
	queryComplexity  int           `cmd:"query-complexity" env:"GQL_QUERY_COMPLEXITY" desc:"GQL query max complexity"`
	dbHost           netip.Addr    `cmd:"db-host" env:"DB_HOST" desc:"database host address"`
	dbPort           uint16        `cmd:"db-port" env:"DB_PORT" desc:"database port"`
	dbUser           string        `cmd:"db-user" env:"DB_USER" desc:"database user"`
	dbPassword       string        `cmd:"db-password" env:"DB_PASSWORD" desc:"database user password"`
	dbName           string        `cmd:"db-name" env:"DB_NAME" desc:"database name"`
	dbTimeout        int32         `cmd:"db-timeout" env:"DB_TIMEOUT" desc:"database connection timeout"`
	dbMaxConnections int32         `cmd:"db-connections" env:"DB_CONNECTIONS" desc:"max database connections"`
	logLevel         string        `cmd:"log-level" env:"LOG_LEVEL" desc:"log level: debug, info, warn, error, fatal, panic, trace, disable"`
}

func (config *Config) ZeroLogLevel() zerolog.Level {
	switch strings.ToLower(config.logLevel) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "disable":
		return zerolog.Disabled
	case "trace":
		return zerolog.TraceLevel
	default:
		return zerolog.WarnLevel
	}
}