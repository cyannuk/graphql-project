package repository

import (
	"context"

	"graphql-project/domain/model"
)

//go:generate go run gen.go model.Review

type ReviewRepository DataSource

func (r *ReviewRepository) GetProductReviews(ctx context.Context, productId int64, offset int32, limit int32) ([]model.Review, error) {
	var reviews model.Reviews = make([]model.Review, 0, max(int(limit), 128))
	err := FindEntities(ctx, (*DataSource)(r), &reviews, SelectByRefId(ctx, "productId", productId, offset, limit))
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
