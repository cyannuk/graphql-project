package repository

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	gotils "github.com/savsgio/gotils/strconv"

	"graphql-project/interface/core"
	"graphql-project/interface/model"
)

const repoContextKey = "repo_context"

type selection struct {
	Offset  int32
	Limit   int32
	Columns []string
}

type names []string
type fields []graphql.CollectedField

func (f fields) Get(i int) string {
	return f[i].Name
}

func (f fields) Length() int {
	return len(f)
}

func (n names) Get(i int) string {
	return n[i]
}

func (n names) Length() int {
	return len(n)
}

func getContextFields(ctx context.Context) core.StringArray {
	if s, ok := ctx.Value(repoContextKey).(selection); ok {
		n := names(s.Columns)
		if len(n) > 0 {
			return &n
		}
	} else if graphql.GetFieldContext(ctx) != nil {
		f := fields(graphql.CollectFieldsCtx(ctx, nil))
		if len(f) > 0 {
			return &f
		}
	}
	return nil
}

func getFields(ctx context.Context, entity model.Entity) (string, []any) {
	fields := getContextFields(ctx)
	if fields == nil {
		return entity.Fields()
	}
	args := make([]any, 0, fields.Length()+1)
	names := make([]byte, 0, 128)

	hasIdentity := false
	identityName, identityArg := entity.Identity()

	for i := 0; i < fields.Length(); i++ {
		if name, arg := entity.Field(fields.Get(i)); arg != nil {
			if len(args) > 0 {
				names = append(names, ',')
			}
			names = append(names, '"')
			names = append(names, name...)
			names = append(names, '"')
			args = append(args, arg)
			if name == identityName {
				hasIdentity = true
			}
		}
	}

	if !hasIdentity && identityArg != nil {
		if len(args) > 0 {
			names = append(names, ',')
		}
		names = append(names, '"')
		names = append(names, identityName...)
		names = append(names, '"')
		args = append(args, identityArg)
	}

	return gotils.B2S(names), args
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
