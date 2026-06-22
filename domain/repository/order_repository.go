package repository

import (
	"context"

	"graphql-project/domain/model"
	i "graphql-project/interface/model"
)

//go:generate go run gen.go model.Order

type OrderRepository DataSource

func (r *OrderRepository) GetUserOrders(ctx context.Context, userId int64, offset int32, limit int32, sort model.Sort) ([]*model.Order, error) {
	orders := make([]*model.Order, 0, max(int(limit), 128))
	err := FindEntities(ctx, (*DataSource)(r), &model.Order{}, SelectByRefId(ctx, "userId", userId, offset, limit, sort), func(ordinality int64, entity i.Entity) {
		orders = append(orders, entity.(*model.Order))
	})
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrdersByUserIds(ctx context.Context, userIds []int64) ([][]*model.Order, []error) {
	offset, limit, sort := getContextRange(ctx)
	userOrders := make([][]*model.Order, len(userIds))

	err := FindEntities(ctx, (*DataSource)(r), &model.Order{}, SelectByRefIds(ctx, "userId", userIds, offset, limit, sort), func(ordinality int64, entity i.Entity) {
		if !entity.IsEmpty() {
			orders := userOrders[ordinality]
			if len(orders) == 0 {
				orders = make([]*model.Order, 0, max(int(limit), 128))
			}
			orders = append(orders, entity.(*model.Order))
			userOrders[ordinality] = orders
		}
	})
	if err != nil {
		return nil, []error{err}
	}

	return userOrders, nil
}

func (r *OrderRepository) GetOrdersByProductIds(ctx context.Context, productIds []int64) ([][]*model.Order, []error) {
	offset, limit, sort := getContextRange(ctx)
	productOrders := make([][]*model.Order, len(productIds))

	err := FindEntities(ctx, (*DataSource)(r), &model.Order{}, SelectByRefIds(ctx, "productId", productIds, offset, limit, sort), func(ordinality int64, entity i.Entity) {
		if !entity.IsEmpty() {
			orders := productOrders[ordinality]
			if len(orders) == 0 {
				orders = make([]*model.Order, 0, max(int(limit), 128))
			}
			orders = append(orders, entity.(*model.Order))
			productOrders[ordinality] = orders
		}
	})

	if err != nil {
		return nil, []error{err}
	}

	return productOrders, nil
}
