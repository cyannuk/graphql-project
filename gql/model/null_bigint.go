package model

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"

	"graphql-project/domain/model"
)

func MarshalNullBigInt(v model.NullBigInt) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if v.State == model.Exists {
			io.WriteString(w, strconv.FormatInt(v.Value, 10))
		} else {
			io.WriteString(w, "null")
		}
	})
}

func UnmarshalNullBigInt(v interface{}) (r model.NullBigInt, err error) {
	if v == nil {
		r.State = model.Null
		return
	}
	if n, ok := v.(int64); ok {
		r.Value = n
		r.State = model.Exists
	} else if s, ok := v.(string); ok {
		if n, err = strconv.ParseInt(s, 10, 64); err == nil {
			r.Value = n
			r.State = model.Exists
		}
	} else {
		err = fmt.Errorf("invalid NullBigInt value '%v'", v)
	}
	return
}
