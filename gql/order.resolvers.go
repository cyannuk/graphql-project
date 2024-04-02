package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	"graphql-project/domain/model"
	"graphql-project/gql/dataloader"
	"graphql-project/tracing"
)

// User is the resolver for the user field.
func (r *orderResolver) User(ctx context.Context, obj *model.Order) (*model.User, error) {
	ctx, span := tracing.InitSpan(ctx, "/query/Order.User")
	defer span.End()
	loaders := dataloader.FromContext(ctx)
	if user, err := loaders.UserLoader.Load(ctx, obj.UserId); err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

// Product is the resolver for the user product.
func (r *orderResolver) Product(ctx context.Context, obj *model.Order) (*model.Product, error) {
	ctx, span := tracing.InitSpan(ctx, "/query/Order.Product")
	defer span.End()
	loaders := dataloader.FromContext(ctx)
	if product, err := loaders.ProductLoader.Load(ctx, obj.ProductId); err != nil {
		return nil, err
	} else {
		return product, nil
	}
}

// Order returns OrderResolver implementation.
func (r *Resolver) Order() OrderResolver { return &orderResolver{r} }

type orderResolver struct{ *Resolver }
