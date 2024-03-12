package model

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"

	"graphql-project/domain/model"
)

func MarshalNullFloat(v model.NullFloat) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if v.State == model.Exists {
			io.WriteString(w, strconv.FormatFloat(float64(v.Value), 'f', -1, 32))
		} else {
			io.WriteString(w, "null")
		}
	})
}

func UnmarshalNullFloat(v interface{}) (r model.NullFloat, err error) {
	if v == nil {
		r.State = model.Null
		return
	}
	if f, ok := v.(float64); ok {
		r.Value = float32(f)
		r.State = model.Exists
	} else if s, ok := v.(string); ok {
		if f, err = strconv.ParseFloat(s, 32); err == nil {
			r.Value = float32(f)
			r.State = model.Exists
		}
	} else {
		err = fmt.Errorf("invalid NullFloat value '%v'", v)
	}
	return
}
