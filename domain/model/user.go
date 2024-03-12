package model

import (
	"time"
)

//go:generate go run gen.go

type Role int32

const (
	RoleAnon  Role = -1
	RoleUser  Role = 0
	RoleAdmin Role = 10
)

type User struct {
	ID        int64
	CreatedAt time.Time
	Name      string
	Email     string
	Address   string
	City      string
	State     string
	Zip       string
	BirthDate time.Time
	Latitude  float64
	Longitude float64
	Password  string
	Source    string
	DeletedAt *time.Time
	Role      Role
}
