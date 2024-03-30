package model

import (
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"

	"graphql-project/core"
	"graphql-project/domain/model"
)

func MarshalDate(v model.Date) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, core.Quote((time.Time)(v).Format(time.DateOnly)))
	})
}

func UnmarshalDate(v interface{}) (d model.Date, err error) {
	date, err := time.Parse(time.DateOnly, v.(string))
	if err != nil {
		return
	}
	return model.Date(date), nil
}
