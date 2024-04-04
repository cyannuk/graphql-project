package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"graphql-project/core"
	"graphql-project/interface/model"
)

func FindEntity(ctx context.Context, dataSource *DataSource, entity model.Entity, qb selectQuery) (err error) {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return
	}
	defer connection.Release()
	qb.Build(entity)
	rows, err := connection.Query(ctx, qb.Query(), qb.Args()...)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		entity.ScanRow(rows)
		return
	}
	err = rows.Err()
	if err == nil {
		err = core.ErrNotFound
	}
	return
}

func FindEntities(ctx context.Context, dataSource *DataSource, entities model.Entities, qb selectQuery) (err error) {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return
	}
	defer connection.Release()
	entity := entities.NewEntity()
	qb.Build(entity)
	rows, err := connection.Query(ctx, qb.Query(), qb.Args()...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		entity.ScanRow(rows)
		entities.Add(entity)
	}
	err = rows.Err()
	return
}

func InsertEntity(ctx context.Context, dataSource *DataSource, inputEntity model.InputEntity) (entity model.Entity, err error) {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return
	}
	defer connection.Release()
	entity = inputEntity.NewEntity()
	qb := InsertInto(entity.Table())
	inputEntity.EnumerateFields(qb.Value)
	rows, err := connection.Query(ctx, qb.Query(getFields(ctx, entity)), qb.Args()...)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		entity.ScanRow(rows)
		return
	}
	err = rows.Err()
	if err == nil {
		err = core.ErrNotFound
	}
	return
}

func UpdateEntity(ctx context.Context, dataSource *DataSource, id int64, inputEntity model.InputEntity) (entity model.Entity, err error) {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return
	}
	defer connection.Release()
	entity = inputEntity.NewEntity()
	qb := Update(entity.Table())
	inputEntity.EnumerateFields(qb.Value)
	qb.Where(entity.Identity(), id)
	rows, err := connection.Query(ctx, qb.Query(getFields(ctx, entity)), qb.Args()...)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		entity.ScanRow(rows)
		return
	}
	err = rows.Err()
	if err == nil {
		err = core.ErrNotFound
	}
	return
}
