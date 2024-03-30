package model

import (
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"

	"graphql-project/core"
	"graphql-project/domain/model"
)

func MarshalNullDate(v model.NullDate) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if v.State == model.Exists {
			io.WriteString(w, core.Quote((time.Time)(v.Value).Format(time.DateOnly)))
		} else {
			io.WriteString(w, "null")
		}
	})
}

func UnmarshalNullDate(v interface{}) (d model.NullDate, err error) {
	if v == nil {
		d.State = model.Null
		return
	}
	date, err := time.Parse(time.DateOnly, v.(string))
	if err == nil {
		d.Value = model.Date(date)
		d.State = model.Exists
	}
	return
}
