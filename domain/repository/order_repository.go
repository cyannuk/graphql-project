package repository

import (
	"context"

	"graphql-project/domain/model"
)

type OrderRepository DataSource

func (r *OrderRepository) GetOrderByID(ctx context.Context, id int64) (*model.Order, error) {
	var order model.Order
	err := FindEntity(ctx, (*DataSource)(r), &order, `SELECT {fields} FROM orders WHERE id = $1 AND "deletedAt" IS NULL`, id)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetUserOrders(ctx context.Context, userId int64, offset int32, limit int32) ([]model.Order, error) {
	orders := model.NewOrders(max(int(limit), 128))
	err := FindEntities(ctx, (*DataSource)(r), &orders, `SELECT {fields} FROM orders WHERE "userId" = $1 AND "deletedAt" IS NULL ORDER BY id`, offset, limit, userId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrders(ctx context.Context, offset int32, limit int32) ([]model.Order, error) {
	orders := model.NewOrders(max(int(limit), 128))
	err := FindEntities(ctx, (*DataSource)(r), &orders, `SELECT {fields} FROM orders WHERE "deletedAt" IS NULL ORDER BY id`, offset, limit)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderByIds(ctx context.Context, ids []int64) ([]*model.Order, []error) {
	orders := model.NewPtrOrders(max(len(ids), 128))
	err := FindEntities(ctx, (*DataSource)(r), &orders, `SELECT {fields} FROM orders JOIN UNNEST($1::BIGINT[]) WITH ORDINALITY t(id, n) USING(id) WHERE "deletedAt" IS NULL ORDER BY t.n`, 0, 0, ids)
	if err != nil {
		return nil, []error{err}
	}
	if len(orders) < len(ids) {
		buffer := make([]*model.Order, len(ids))
		n := 0
		for i, id := range ids {
			order := orders[n]
			if order.ID == id {
				buffer[i] = order
				n++
			}
		}
		orders = buffer
	}
	return orders, nil
}

func NewOrderRepository(dataSource *DataSource) *OrderRepository {
	return (*OrderRepository)(dataSource)
}
