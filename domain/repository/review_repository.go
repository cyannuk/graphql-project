package repository

import (
	"context"

	"graphql-project/domain/model"
	i "graphql-project/interface/model"
)

//go:generate go run gen.go model.Review

type ReviewRepository DataSource

func (r *ReviewRepository) GetProductReviews(ctx context.Context, productId int64, offset int32, limit int32, sort model.Sort) ([]*model.Review, error) {
	reviews := make([]*model.Review, 0, max(int(limit), 128))
	err := FindEntities(ctx, (*DataSource)(r), &model.Review{}, SelectByRefId(ctx, "productId", productId, offset, limit, sort), func(ordinality int64, entity i.Entity) {
		reviews = append(reviews, entity.(*model.Review))
	})
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
