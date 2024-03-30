package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"graphql-project/core"
	"graphql-project/interface/model"
)

func FindEntity(ctx context.Context, dataSource *DataSource, entity model.Entity, query string, args ...any) (err error) {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return
	}
	defer connection.Release()
	rows, err := connection.Query(ctx, core.Replace(query, "{fields}", getFields(ctx, entity), 1), args...)
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

func FindEntities(ctx context.Context, dataSource *DataSource, entities model.Entities, query string, offset int32, limit int32, args ...any) (err error) {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return
	}
	defer connection.Release()

	entity := entities.New()
	query = core.Replace(query, "{fields}", getFields(ctx, entity), 1)
	if offset > 0 {
		query = core.Join(query, " OFFSET ", core.Int32ToStr(offset))
	}
	if limit > 0 {
		query = core.Join(query, " LIMIT ", core.Int32ToStr(limit))
	}

	rows, err := connection.Query(ctx, query, args...)
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

func InsertEntity(ctx context.Context, dataSource *DataSource, entity model.Entity, inputEntity model.InputEntity) (err error) {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return
	}
	defer connection.Release()
	fields, placeholders, args := inputEntity.InsertFields()
	query := core.Join("INSERT INTO ", entity.Table(), "(", fields, ") VALUES(", placeholders, ") RETURNING ", getFields(ctx, entity))
	rows, err := connection.Query(ctx, query, args...)
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

func UpdateEntity(ctx context.Context, dataSource *DataSource, id int64, entity model.Entity, inputEntity model.InputEntity) (err error) {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return
	}
	defer connection.Release()
	fields, placeholders, args := inputEntity.InsertFields()
	args = append(args, id)
	query := core.Join("UPDATE ", entity.Table(), " SET (", fields, ") = (", placeholders, ") WHERE ", entity.Identity(), " = $", core.IntToStr(len(args)), " RETURNING ", getFields(ctx, entity))
	rows, err := connection.Query(ctx, query, args...)
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
