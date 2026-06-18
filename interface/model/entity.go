package model

import (
	"github.com/jackc/pgx/v5"
)

type Entity interface {
	ScanRow(rows pgx.Rows) int64
	Table() string
	Field(name string) string
	Fields() []string
	Identity() string
}

type InputEntity interface {
	NewEntity() Entity
	EnumerateFields(func(name string, value any))
}

type InitEntity interface {
	InputEntity
	Entity
}
