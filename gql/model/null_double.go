package model

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"

	"graphql-project/domain/model"
)

func MarshalNullDouble(v model.NullDouble) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if v.State == model.Exists {
			io.WriteString(w, strconv.FormatFloat(v.Value, 'f', -1, 64))
		} else {
			io.WriteString(w, "null")
		}
	})
}

func UnmarshalNullDouble(v interface{}) (r model.NullDouble, err error) {
	if v == nil {
		r.State = model.Null
		return
	}
	if f, ok := v.(float64); ok {
		r.Value = f
		r.State = model.Exists
	} else if s, ok := v.(string); ok {
		if f, err = strconv.ParseFloat(s, 64); err == nil {
			r.Value = f
			r.State = model.Exists
		}
	} else {
		err = fmt.Errorf("invalid NullDouble value '%v'", v)
	}
	return
}
