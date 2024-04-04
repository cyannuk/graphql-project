package repository

import (
	"context"

	"graphql-project/domain/model"
)

//go:generate go run gen.go model.User

type UserRepository DataSource

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := FindEntity(ctx, (*DataSource)(r), &user, SelectBy(ctx, "email", email))
	if err != nil {
		return nil, err
	}
	return &user, nil
}
