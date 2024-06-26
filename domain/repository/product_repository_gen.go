// Code generated by gen; DO NOT EDIT.
package repository

import (
	"context"
	"graphql-project/domain/model"
)

func (r *ProductRepository) GetProductByID(ctx context.Context, id int64) (*model.Product, error) {
	var product model.Product
	err := FindEntity(ctx, (*DataSource)(r), &product, SelectById(ctx, id))
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetProducts(ctx context.Context, offset int32, limit int32) ([]model.Product, error) {
	var products model.Products = make([]model.Product, 0, max(int(limit), 128))
	err := FindEntities(ctx, (*DataSource)(r), &products, SelectMany(ctx, offset, limit))
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetProductByIds(ctx context.Context, ids []int64) ([]*model.Product, []error) {
	var products model.ProductRefs = make([]*model.Product, 0, len(ids))
	err := FindEntities(ctx, (*DataSource)(r), &products, SelectByIds(ctx, ids))
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
	product, err := InsertEntity(ctx, (*DataSource)(r), inputProduct)
	if err != nil {
		return nil, err
	}
	return product.(*model.Product), err
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, id int64, inputProduct *model.ProductInput) (*model.Product, error) {
	product, err := UpdateEntity(ctx, (*DataSource)(r), id, inputProduct)
	if err != nil {
		return nil, err
	}
	return product.(*model.Product), err
}

func NewProductRepository(dataSource *DataSource) *ProductRepository {
	return (*ProductRepository)(dataSource)
}
