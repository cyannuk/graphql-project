package model

import (
	"database/sql"
	"time"
)

//go:generate reform

//reform:users
type User struct {
	ID        int64        `reform:"id,pk"`
	CreatedAt time.Time    `reform:"createdAt"`
	Name      string       `reform:"name"`
	Email     string       `reform:"email"`
	Address   string       `reform:"address"`
	City      string       `reform:"city"`
	State     string       `reform:"state"`
	Zip       string       `reform:"zip"`
	BirthDate time.Time    `reform:"birthDate"`
	Latitude  float64      `reform:"latitude"`
	Longitude float64      `reform:"longitude"`
	Password  string       `reform:"password"`
	Source    string       `reform:"source"`
	DeletedAt sql.NullTime `reform:"deletedAt"`
}
