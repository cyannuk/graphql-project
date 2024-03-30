package repository

import (
	"context"

	"graphql-project/domain/model"
)

//go:generate go run gen.go model.Order

type OrderRepository DataSource

func (r *OrderRepository) GetUserOrders(ctx context.Context, userId int64, offset int32, limit int32) ([]model.Order, error) {
	var orders model.Orders = make([]model.Order, 0, max(int(limit), 128))
	err := FindEntities(ctx, (*DataSource)(r), &orders, `SELECT {fields} FROM orders WHERE "userId" = $1 AND "deletedAt" IS NULL ORDER BY id`, offset, limit, userId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrdersByUsersIds(ctx context.Context, userIds []int64) ([][]model.Order, []error) {
	var from, to int64
	offset, limit := getContextRange(ctx)
	if offset > 0 {
		from = int64(offset)
	} else {
		from = 1
	}
	if limit > 0 {
		to = from + int64(limit)
	} else {
		to = 9_223_372_036_854_775_807
	}
	var orders model.Orders = make([]model.Order, 0, len(userIds) * 128)
	err := FindEntities(ctx, (*DataSource)(r), &orders,
		`WITH o AS (`+
			`  SELECT *, ROW_NUMBER() OVER (PARTITION BY "userId" ORDER BY "id") AS r `+
			`  FROM orders `+
			`  WHERE "deletedAt" IS NULL AND "userId" = ANY($1::BIGINT[])`+
			`) `+
			`SELECT {fields} FROM o `+
			`JOIN UNNEST($1::BIGINT[]) WITH ORDINALITY t("userId", n) USING("userId") `+
			`WHERE r >= $2 AND r < $3 `+
			`ORDER BY t.n, r`, 0, 0, userIds, from, to)
	if err != nil {
		return nil, []error{err}
	}
	userOrders := make([][]model.Order, len(userIds))
	var n, k int
	for i, userId := range userIds {
		for k < len(orders) && orders[k].UserId == userId {
			k++
		}
		if k > n {
			userOrders[i] = orders[n:k]
			n = k
		}
	}
	return userOrders, nil
}
