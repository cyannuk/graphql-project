package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/jackc/pgx/v5/pgxpool"

	"graphql-project/interface/model"
)

type FmtQuery = func(fields string) string
type NextRow = func()

func GetContextFields(ctx context.Context, entity model.Entity) (string, []any) {
	if graphql.GetFieldContext(ctx) == nil {
		return entity.Fields()
	}

	fields := graphql.CollectFieldsCtx(ctx, nil)
	args := make([]any, 0, len(fields)+1)

	var sb strings.Builder
	sb.Grow(256)
	sb.WriteByte(' ')

	hasIdentity := false
	identityName, identityArg := entity.Identity()

	for _, field := range fields {
		if name, arg := entity.Field(field.Name); arg != nil {
			if len(args) > 0 {
				sb.WriteByte(',')
			}
			sb.WriteByte('"')
			sb.WriteString(name)
			sb.WriteByte('"')
			args = append(args, arg)
			if name == identityName {
				hasIdentity = true
			}
		}
	}

	if !hasIdentity && identityArg != nil {
		if len(args) > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte('"')
		sb.WriteString(identityName)
		sb.WriteByte('"')
		args = append(args, identityArg)
	}

	sb.WriteByte(' ')
	return sb.String(), args
}

func FindEntity(ctx context.Context, dataSource *DataSource, entity model.Entity, fmtQuery FmtQuery, args ...any) error {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()
	fieldList, fields := GetContextFields(ctx, entity)
	if rows, err := connection.Query(ctx, fmtQuery(fieldList), args...); err != nil {
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

func FindEntities(ctx context.Context, dataSource *DataSource, entity model.Entity, fmtQuery FmtQuery, nextRow NextRow, offset int32, limit int32, args ...any) error {
	connection, err := (*pgxpool.Pool)(dataSource).Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()

	fieldList, fields := GetContextFields(ctx, entity)
	query := fmtQuery(fieldList)
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
			nextRow()
		}
		return nil
	}
}
