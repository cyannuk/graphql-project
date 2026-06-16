package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"graphql-project/core"
	i "graphql-project/interface/model"
)

func FindEntity(ctx context.Context, dataSource *DataSource, entity i.Entity, qb selectQuery) (err error) {
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

func FindEntities(ctx context.Context, dataSource *DataSource, entity i.InitEntity, qb selectQuery, yield func(ordinality int64, entity i.Entity)) (err error) {
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

	for rows.Next() {
		e := entity.NewEntity()
		ordinality, empty := e.ScanRow(rows)
		if empty {
			yield(ordinality, nil)
		} else {
			yield(ordinality, e)
		}
	}

	err = rows.Err()
	return
}

func InsertEntity(ctx context.Context, dataSource *DataSource, inputEntity i.InputEntity) (entity i.Entity, err error) {
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

func UpdateEntity(ctx context.Context, dataSource *DataSource, id int64, inputEntity i.InputEntity) (entity i.Entity, err error) {
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
