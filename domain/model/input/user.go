package input

import (
	"time"

	"graphql-project/domain/model"
)

//go:generate go run gen.go

type User struct {
	Name      model.NullString
	Email     model.NullString
	Address   model.NullString
	City      model.NullString
	State     model.NullString
	Zip       model.NullString
	BirthDate model.NullTime
	Latitude  model.NullDouble
	Longitude model.NullDouble
	Password  model.NullString
	Source    model.NullString
}

type NewUser struct {
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
}
