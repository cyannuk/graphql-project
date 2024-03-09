package repository

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"

	"graphql-project/interface/core"
	"graphql-project/interface/model"
)

const RepoContextKey = "repo_context"

type names []string
type fields []graphql.CollectedField

func (f *fields) Get(i int) string {
	return (*f)[i].Name
}

func (f *fields) Length() int {
	return len(*f)
}

func (n *names) Get(i int) string {
	return (*n)[i]
}

func (n *names) Length() int {
	return len(*n)
}

func getContextFields(ctx context.Context) core.StringArray {
	if graphql.GetFieldContext(ctx) != nil {
		f := fields(graphql.CollectFieldsCtx(ctx, nil))
		if len(f) > 0 {
			return &f
		}
	} else if s, ok := ctx.Value(RepoContextKey).([]string); ok {
		n := names(s)
		if len(n) > 0 {
			return &n
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

	var sb strings.Builder
	sb.Grow(256)
	sb.WriteByte(' ')

	hasIdentity := false
	identityName, identityArg := entity.Identity()

	for i := 0; i < fields.Length(); i++ {
		if name, arg := entity.Field(fields.Get(i)); arg != nil {
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
