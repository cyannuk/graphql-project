package model

import (
	"errors"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalDate(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, t.Format(`"2006-01-02"`))
	})
}

func UnmarshalDate(v interface{}) (time.Time, error) {
	if str, ok := v.(string); !ok || (len(str) != 12) {
		return time.Time{}, errors.New("invalid date format")
	} else {
		return time.Parse(time.DateOnly, str[1:len(str)-1])
	}
}
