package config

import (
	"net/netip"
	"time"
)

//go:generate go run gen.go

type Config struct {
	bindAddr             netip.Addr    `flag:"bind-addr" env:"ADDRESS" desc:"bind host address" default:"0.0.0.0"`
	port                 uint16        `flag:"port" env:"PORT" desc:"listen port" default:"8080"`
	jwtSecret            []byte        `flag:"jwt-secret" env:"JWT_SECRET" desc:"base64 encoded JWT secret"`
	jwtExpiration        time.Duration `flag:"jwt-expiration" env:"JWT_EXPIRATION" desc:"JWT access expiration time: 1h20m30s"`
	jwtRefreshExpiration time.Duration `flag:"jwt-refresh-expiration" env:"JWT_REFRESH_EXPIRATION" desc:"JWT refresh expiration time: 1h20m30s"`
	dbHost               netip.Addr    `flag:"db-host" env:"DB_HOST" desc:"database host address" default:"localhost"`
	dbPort               uint16        `flag:"db-port" env:"DB_PORT" desc:"database port" default:"5432"`
	dbUser               string        `flag:"db-user" env:"DB_USER" desc:"database user" default:"postgres"`
	dbPassword           string        `flag:"db-password" env:"DB_PASSWORD" desc:"database user password"`
	dbName               string        `flag:"db-name" env:"DB_NAME" desc:"database name"`
	dbTimeout            time.Duration `flag:"db-timeout" env:"DB_TIMEOUT" desc:"database connection timeout" default:"5s"`
	dbMaxConnections     uint32        `flag:"db-connections" env:"DB_CONNECTIONS" desc:"max database connections"`
	dbMigrate            bool          `flag:"db-migrate" env:"DB_MIGRATE" desc:"Apply database migrations" default:"true"`
	queryComplexity      int           `flag:"query-complexity" env:"GQL_QUERY_COMPLEXITY" desc:"GQL query max complexity" default:"2000"`
	logLevel             string        `flag:"log-level" env:"LOG_LEVEL" desc:"log level: debug|info|warn|error|fatal|trace|disable" default:"info"`
	enableTracing        bool          `flag:"enable-tracing" env:"ENABLE_TRACING" desc:"Enable API request tracing" default:"false" optional:"true"`
}
