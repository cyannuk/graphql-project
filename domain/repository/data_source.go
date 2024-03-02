package repository

import (
	"database/sql"
	"runtime"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

type DataSource struct {
	*reform.DB
}

func (dataSource DataSource) Close() error {
	return dataSource.DBInterface().(*sql.DB).Close()
}

func NewDataSource(connectionString string) (DataSource, error) {
	connConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return DataSource{}, err
	}
	db := stdlib.OpenDB(*connConfig)
	db.SetMaxOpenConns(runtime.NumCPU() * 4)
	if err = db.Ping(); err != nil {
		return DataSource{}, err
	}
	return DataSource{reform.NewDB(db, postgresql.Dialect, nil)}, nil
}
