package model

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"

	"graphql-project/core"
	"graphql-project/domain/model"
)

func MarshalNullString(str model.NullString) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if str.State == model.Exists {
			io.WriteString(w, core.Quote(str.Value))
		} else {
			io.WriteString(w, "null")
		}
	})
}

func UnmarshalNullString(v interface{}) (r model.NullString, err error) {
	if v == nil {
		r.State = model.Null
		return
	}
	if s, ok := v.(string); ok {
		r.Value = s
		r.State = model.Exists
	} else {
		err = fmt.Errorf("invalid NullString value '%v'", v)
	}
	return
}
