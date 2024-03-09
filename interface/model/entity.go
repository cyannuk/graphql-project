package model

type Entity interface {
	Field(name string) (string, any)
	Fields() (string, []any)
	Identity() (string, any)
}

type Entities interface {
	New() Entity
	Add(entity Entity)
}
