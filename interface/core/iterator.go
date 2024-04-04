package core

type Iterator interface {
	Get() string
	Next() bool
}
