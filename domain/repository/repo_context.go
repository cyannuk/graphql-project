package repository

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	. "graphql-project/interface/core"
	"graphql-project/interface/model"
)

const repoContextKey = "repo_context"

type selection struct {
	Offset  int32
	Limit   int32
	Columns []string
}

type columnIterator struct {
	columns     []string
	current     string
	i           int
	identity    string
	hasIdentity bool
}

type entityFieldIterator struct {
	entity  model.Entity
	current string
	i       int
}

type contextFieldIterator struct {
	fields      []graphql.CollectedField
	current     string
	i           int
	entity      model.Entity
	hasIdentity bool
}

func (iter *contextFieldIterator) Get() string {
	return iter.current
}

func (iter *contextFieldIterator) Next() bool {
	l := len(iter.fields) - 1
	for iter.i < l {
		iter.i++
		f := iter.entity.Field(iter.fields[iter.i].Name)
		if f != "" {
			if f == iter.entity.Identity() {
				iter.hasIdentity = true
			}
			iter.current = f
			return true
		}
	}
	if !iter.hasIdentity && iter.i == l {
		iter.i++
		iter.current = iter.entity.Identity()
		return true
	}
	return false
}

func (iter *columnIterator) Get() string {
	return iter.current
}

func (iter *columnIterator) Next() bool {
	l := len(iter.columns) - 1
	if iter.i < l {
		iter.i++
		c := iter.columns[iter.i]
		if c == iter.identity {
			iter.hasIdentity = true
		}
		iter.current = c
		return true
	}
	if !iter.hasIdentity && iter.i == l {
		iter.i++
		iter.current = iter.identity
		return true
	}
	return false
}

func (iter *entityFieldIterator) Get() string {
	return iter.current
}

func (iter *entityFieldIterator) Next() bool {
	fields := iter.entity.Fields()
	if iter.i < len(fields)-1 {
		iter.i++
		iter.current = fields[iter.i]
		return true
	}
	return false
}

func getFields(ctx context.Context, entity model.Entity) Iterator {
	if s, ok := ctx.Value(repoContextKey).(selection); ok && len(s.Columns) > 0 {
		return &columnIterator{columns: s.Columns, identity: entity.Identity(), i: -1}
	} else if graphql.GetFieldContext(ctx) != nil {
		return &contextFieldIterator{fields: graphql.CollectFieldsCtx(ctx, nil), entity: entity, i: -1}
	} else {
		return &entityFieldIterator{entity: entity, i: -1}
	}
}

func getContextRange(ctx context.Context) (int32, int32) {
	if s, ok := ctx.Value(repoContextKey).(selection); ok {
		return s.Offset, s.Limit
	}
	return 0, 0
}

func With(ctx context.Context, offset int32, limit int32, columns ...string) context.Context {
	return context.WithValue(ctx, repoContextKey, selection{offset, limit, columns})
}
