package gql

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"graphql-pro/domain/model"
	"graphql-pro/interface/repository"
)

type Resolver struct {
	UsersRepo repository.UserRepository
}

// UserByID is the resolver for the userByID field.
func (r Resolver) UserByID(ctx context.Context, id int64) (*model.User, error) {
	fields := graphql.CollectFieldsCtx(ctx, nil)
	for _, field := range fields {
		fmt.Printf("!! field.Name = %s\n", field.Name)
	}
	return r.UsersRepo.GetUserByID(id)
}

// UserByEmail is the resolver for the userByEmail field.
func (r Resolver) UserByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.UsersRepo.GetUserByEmail(email)
}

// Query returns QueryResolver implementation.
func (r Resolver) Query() QueryResolver {
	return r
}
