package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"graphql-project/core"
	"graphql-project/interface/model"
)

func FindEntity(ctx context.Context, dataSource *DataSource, entity model.Entity, query string, args ...any) error {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()
	fieldList, fields := getFields(ctx, entity)
	if rows, err := connection.Query(ctx, core.Replace(query, "{fields}", fieldList, 1), args...); err != nil {
		return err
	} else {
		defer rows.Close()
		if rows.Next() {
			if err := rows.Scan(fields...); err != nil {
				return err
			}
			return nil
		}
		return core.ErrNotFound
	}
}

func FindEntities(ctx context.Context, dataSource *DataSource, entities model.Entities, query string, offset int32, limit int32, args ...any) error {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()

	entity := entities.New()
	fieldList, fields := getFields(ctx, entity)
	query = core.Replace(query, "{fields}", fieldList, 1)
	if offset > 0 {
		query = core.Join(query, " OFFSET ", core.Int32ToStr(offset))
	}
	if limit > 0 {
		query = core.Join(query, " LIMIT ", core.Int32ToStr(limit))
	}

	if rows, err := connection.Query(ctx, query, args...); err != nil {
		return err
	} else {
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(fields...); err != nil {
				return err
			}
			entities.Add(entity)
		}
		return nil
	}
}

func InsertEntity(ctx context.Context, dataSource *DataSource, entity model.Entity, inputEntity model.InputEntity) error {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()
	fieldList, fields := getFields(ctx, entity)
	insertFieldList, valueList, args := inputEntity.Fields()
	query := core.Join("INSERT INTO ", entity.Table(), "(", insertFieldList, ") VALUES(", valueList, ") RETURNING ", fieldList)
	if rows, err := connection.Query(ctx, query, args...); err != nil {
		return err
	} else {
		defer rows.Close()
		if rows.Next() {
			if err := rows.Scan(fields...); err != nil {
				return err
			}
			return nil
		}
		return core.ErrNotFound
	}
}

func UpdateEntity(ctx context.Context, dataSource *DataSource, id int64, entity model.Entity, inputEntity model.InputEntity) error {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()
	key, _ := entity.Identity()
	fieldList, fields := getFields(ctx, entity)
	updateFieldList, valueList, args := inputEntity.Fields()
	args = append(args, id)
	query := core.Join("UPDATE ", entity.Table(), " SET (", updateFieldList, ") = (", valueList, ") WHERE ", key, " = $", core.IntToStr(len(args)), " RETURNING ", fieldList)
	if rows, err := connection.Query(ctx, query, args...); err != nil {
		return err
	} else {
		defer rows.Close()
		if rows.Next() {
			if err := rows.Scan(fields...); err != nil {
				return err
			}
			return nil
		}
		return core.ErrNotFound
	}
}
