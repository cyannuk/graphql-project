package model

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"

	"graphql-project/domain/model"
)

func MarshalNullInt(v model.NullInt) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if v.State == model.Exists {
			io.WriteString(w, strconv.FormatInt(int64(v.Value), 10))
		} else {
			io.WriteString(w, "null")
		}
	})
}

func UnmarshalNullInt(v interface{}) (r model.NullInt, err error) {
	if v == nil {
		r.State = model.Null
		return
	}
	if n, ok := v.(int64); ok {
		r.Value = int32(n)
		r.State = model.Exists
	} else if s, ok := v.(string); ok {
		if n, err = strconv.ParseInt(s, 10, 32); err == nil {
			r.Value = int32(n)
			r.State = model.Exists
		}
	} else {
		err = fmt.Errorf("invalid NullInt value '%v'", v)
	}
	return
}
