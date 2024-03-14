package repository

import (
	"context"

	"graphql-project/domain/model"
)

//go:generate go run gen.go model.Order

type OrderRepository DataSource

func (r *OrderRepository) GetUserOrders(ctx context.Context, userId int64, offset int32, limit int32) ([]model.Order, error) {
	orders := model.NewOrders(max(int(limit), 128))
	err := FindEntities(ctx, (*DataSource)(r), &orders, `SELECT {fields} FROM orders WHERE "userId" = $1 AND "deletedAt" IS NULL ORDER BY id`, offset, limit, userId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
