package model

import (
	"database/sql"
	"time"
)

//go:generate reform

//reform:users
type User struct {
	ID        int64        `reform:"id,pk" json:"id,omit-dec"`
	CreatedAt time.Time    `reform:"createdAt" json:"created_at"`
	Name      string       `reform:"name" json:"name"`
	Email     string       `reform:"email" json:"email"`
	Address   string       `reform:"address" json:"address"`
	City      string       `reform:"city" json:"city"`
	State     string       `reform:"state" json:"state"`
	Zip       string       `reform:"zip" json:"zip"`
	BirthDate time.Time    `reform:"birthDate" json:"birth_date"`
	Latitude  float64      `reform:"latitude" json:"latitude"`
	Longitude float64      `reform:"longitude" json:"longitude"`
	Password  string       `reform:"password" json:"password,omit-enc"`
	Source    string       `reform:"source" json:"source"`
	DeletedAt sql.NullTime `reform:"deletedAt" json:"deleted_at"`
}
