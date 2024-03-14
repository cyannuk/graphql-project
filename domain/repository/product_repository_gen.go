// Code generated by gen; DO NOT EDIT.
package repository

import (
	"context"

	"graphql-project/domain/model"
)

func (r *ProductRepository) GetProductByID(ctx context.Context, id int64) (*model.Product, error) {
	var product model.Product
	err := FindEntity(ctx, (*DataSource)(r), &product, `SELECT {fields} FROM products WHERE id = $1 AND "deletedAt" IS NULL`, id)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetProducts(ctx context.Context, offset int32, limit int32) ([]model.Product, error) {
	products := model.NewProducts(max(int(limit), 128))
	err := FindEntities(ctx, (*DataSource)(r), &products, `SELECT {fields} FROM products WHERE "deletedAt" IS NULL ORDER BY id`, offset, limit)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetProductByIds(ctx context.Context, ids []int64) ([]*model.Product, []error) {
	products := model.NewPtrProducts(max(len(ids), 128))
	err := FindEntities(ctx, (*DataSource)(r), &products, `SELECT {fields} FROM products JOIN UNNEST($1::BIGINT[]) WITH ORDINALITY t(id, n) USING(id) WHERE "deletedAt" IS NULL ORDER BY t.n`, 0, 0, ids)
	if err != nil {
		return nil, []error{err}
	}
	if len(products) < len(ids) {
		buffer := make([]*model.Product, len(ids))
		n := 0
		for i, id := range ids {
			product := products[n]
			if product.ID == id {
				buffer[i] = product
				n++
			}
		}
		products = buffer
	}
	return products, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, inputProduct *model.Product) (*model.Product, error) {
	var product model.Product
	err := InsertEntity(ctx, (*DataSource)(r), &product, inputProduct)
	return &product, err
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, id int64, inputProduct *model.ProductInput) (*model.Product, error) {
	var product model.Product
	err := UpdateEntity(ctx, (*DataSource)(r), id, &product, inputProduct)
	return &product, err
}

func NewProductRepository(dataSource *DataSource) *ProductRepository {
	return (*ProductRepository)(dataSource)
}
