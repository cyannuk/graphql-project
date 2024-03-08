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

func DataSourceConfig(config *config.Config) (*pgxpool.Config, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=%d",
		config.DbHost(), config.DbPort(), config.DbUser(), config.DbPassword(), config.DbName(), config.DbTimeout())
	if poolConfig, err := pgxpool.ParseConfig(connectionString); err != nil {
		return nil, err
	} else {
		poolConfig.ConnConfig.RuntimeParams["client_encoding"] = "UTF8"
		poolConfig.MaxConns = config.DbMaxConnections()
		return poolConfig, nil
	}
}

func NewDataSource(config *config.Config) (*DataSource, error) {
	poolConfig, err := DataSourceConfig(config)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	if pool, err := pgxpool.NewWithConfig(ctx, poolConfig); err != nil {
		return nil, err
	} else {
		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			return nil, err
		}
		return (*DataSource)(pool), nil
	}
}
