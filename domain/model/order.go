package model

import (
	"time"
)

//go:generate go run gen.go

type Order struct {
	ID        int64
	CreatedAt time.Time
	UserId    int64 `gql:"user"`
	ProductId int64 `gql:"product"`
	Discount  float64
	Quantity  int32
	Subtotal  float64
	Tax       float64
	Total     float64
	DeletedAt *time.Time
}
