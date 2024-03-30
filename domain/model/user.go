package model

import (
	"time"
)

//go:generate go run gen.go

type User struct {
	ID        int64
	CreatedAt time.Time `auto:"true"`
	Name      string
	Email     string
	Address   string
	City      string
	State     string
	Zip       string
	BirthDate Date
	Latitude  float64
	Longitude float64
	Password  string
	Source    string
	DeletedAt *time.Time `auto:"true"`
	Role      Role
}
