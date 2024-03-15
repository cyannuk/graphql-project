package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"graphql-project/config"
)

type DataSource pgxpool.Pool

func (dataSource *DataSource) Close() {
	(*pgxpool.Pool)(dataSource).Close()
}

func DataSourceConfig(cfg *config.Config) (poolConfig *pgxpool.Config, err error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=%d",
		cfg.DbHost(), cfg.DbPort(), cfg.DbUser(), cfg.DbPassword(), cfg.DbName(), cfg.DbTimeout())
	poolConfig, err = pgxpool.ParseConfig(connectionString)
	if err == nil {
		poolConfig.ConnConfig.RuntimeParams["client_encoding"] = "UTF8"
		poolConfig.MaxConns = cfg.DbMaxConnections()
	}
	return
}

func NewDataSource(cfg *config.Config) (dataSource *DataSource, err error) {
	poolConfig, err := DataSourceConfig(cfg)
	if err != nil {
		return
	}
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return
	}
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
	} else {
		dataSource = (*DataSource)(pool)
	}
	return
}
