package model

import (
	"github.com/jackc/pgx/v5"
)

type Entity interface {
	pgx.RowScanner
	Table() string
	Field(name string) string
	Fields() string
	Identity() string
}

type Entities interface {
	New() Entity
	Add(entity Entity)
}

type InputEntity interface {
	InsertFields() (string, string, []any)
}
