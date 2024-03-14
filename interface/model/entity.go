package model

type Entity interface {
	Table() string
	Field(name string) (string, any)
	Fields() (string, []any)
	Identity() (string, any)
}

type Entities interface {
	New() Entity
	Add(entity Entity)
}

type InputEntity interface {
	InsertFields() (string, string, []any)
}
