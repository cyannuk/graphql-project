package model

import (
	"time"
)

//go:generate go run gen.go

type Product struct {
	ID        int64
	CreatedAt time.Time `auto:"true"`
	Category  string
	Ean       string
	Price     float64
	Quantity  int32
	Rating    float64
	Name      string
	Vendor    string
	DeletedAt *time.Time `auto:"true"`
}
