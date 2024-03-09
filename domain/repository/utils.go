package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"graphql-project/interface/model"
)

func makeQuery(query string, fieldList string) string {
	i := strings.Index(query, "{fields}")
	if i < 0 {
		return query
	}
	return query[:i] + fieldList + query[i+8:]
}

func FindEntity(ctx context.Context, dataSource *DataSource, entity model.Entity, query string, args ...any) error {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()
	fieldList, fields := getFields(ctx, entity)
	if rows, err := connection.Query(ctx, makeQuery(query, fieldList), args...); err != nil {
		return err
	} else {
		defer rows.Close()
		if rows.Next() {
			if err := rows.Scan(fields...); err != nil {
				return err
			}
		}
		return nil
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
	query = makeQuery(query, fieldList)
	if offset > 0 {
		query = fmt.Sprint(query, " OFFSET ", offset)
	}
	if limit > 0 {
		query = fmt.Sprint(query, " LIMIT ", limit)
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
