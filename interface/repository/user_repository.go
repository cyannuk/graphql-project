package repository

import (
	"graphql-pro/domain/model"
)

type UserRepository interface {
	GetUserByID(id int64) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
}
