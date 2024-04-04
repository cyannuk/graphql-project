package repository

import (
	"context"

	"graphql-project/domain/model"
)

//go:generate go run gen.go model.Order

type OrderRepository DataSource

func (r *OrderRepository) GetUserOrders(ctx context.Context, userId int64, offset int32, limit int32) ([]model.Order, error) {
	var orders model.Orders = make([]model.Order, 0, max(int(limit), 128))
	err := FindEntities(ctx, (*DataSource)(r), &orders, SelectByRefId(ctx, "userId", userId, offset, limit))
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrdersByUserIds(ctx context.Context, userIds []int64) ([][]model.Order, []error) {
	var orders model.Orders = make([]model.Order, 0, len(userIds)*128)
	offset, limit := getContextRange(ctx)
	err := FindEntities(ctx, (*DataSource)(r), &orders, SelectByRefIds(ctx, "userId", userIds, offset, limit))
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

func (r *OrderRepository) GetOrdersByProductIds(ctx context.Context, productIds []int64) ([][]model.Order, []error) {
	var orders model.Orders = make([]model.Order, 0, len(productIds)*128)
	offset, limit := getContextRange(ctx)
	err := FindEntities(ctx, (*DataSource)(r), &orders, SelectByRefIds(ctx, "productId", productIds, offset, limit))
	if err != nil {
		return nil, []error{err}
	}
	productOrders := make([][]model.Order, len(productIds))
	var n, k int
	for i, productId := range productIds {
		for k < len(orders) && orders[k].ProductId == productId {
			k++
		}
		if k > n {
			productOrders[i] = orders[n:k]
			n = k
		}
	}
	return productOrders, nil
}
