package repository

import (
	"context"
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DataSource pgxpool.Pool

func (dataSource *DataSource) Close() {
	(*pgxpool.Pool)(dataSource).Close()
}

func NewDataSource(connectionString string) (*DataSource, error) {
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}
	poolConfig.ConnConfig.RuntimeParams["client_encoding"] = "UTF8"
	poolConfig.MaxConns = int32(runtime.NumCPU() * 4)
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
