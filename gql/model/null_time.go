package model

import (
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"

	"graphql-project/core"
	"graphql-project/domain/model"
)

func MarshalNullTime(v model.NullTime) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if v.State == model.Exists {
			io.WriteString(w, core.Quote(v.Value.Format(time.RFC3339Nano)))
		} else {
			io.WriteString(w, "null")
		}
	})
}

func UnmarshalNullTime(v interface{}) (r model.NullTime, err error) {
	if v == nil {
		r.State = model.Null
		return
	}
	if r.Value, err = time.Parse(time.RFC3339Nano, v.(string)); err == nil {
		r.State = model.Exists
	}
	return
}
