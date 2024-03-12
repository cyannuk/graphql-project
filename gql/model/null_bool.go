package model

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"

	"graphql-project/domain/model"
)

func MarshalNullBool(v model.NullBool) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if v.State == model.Exists {
			io.WriteString(w, strconv.FormatBool(v.Value))
		} else {
			io.WriteString(w, "null")
		}
	})
}

func UnmarshalNullBool(v interface{}) (r model.NullBool, err error) {
	if v == nil {
		r.State = model.Null
		return
	}
	if b, ok := v.(bool); ok {
		r.Value = b
		r.State = model.Exists
	} else if s, ok := v.(string); ok {
		if b, err = strconv.ParseBool(s); err == nil {
			r.Value = b
			r.State = model.Exists
		}
	} else {
		err = fmt.Errorf("invalid NullBool value '%v'", v)
	}
	return
}
