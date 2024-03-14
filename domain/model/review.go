package model

import (
	"time"
)

//go:generate go run gen.go

type Review struct {
	ID        int64
	CreatedAt time.Time `auto:"true"`
	Reviewer  string
	ProductId int64 `gql:"product"`
	Rating    int32
	Body      string
	DeletedAt *time.Time `auto:"true"`
}
