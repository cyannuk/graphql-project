package repository

import (
	"context"

	"graphql-project/domain/model"
)

//go:generate go run gen.go model.Review

type ReviewRepository DataSource

func (r *ReviewRepository) GetProductReviews(ctx context.Context, productId int64, offset int32, limit int32) ([]model.Review, error) {
	reviews := model.NewReviews(max(int(limit), 128))
	err := FindEntities(ctx, (*DataSource)(r), &reviews, `SELECT {fields} FROM reviews WHERE "productId" = $1 AND "deletedAt" IS NULL ORDER BY id`, offset, limit, productId)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
