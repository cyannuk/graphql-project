package model

import (
	"errors"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"

	"graphql-project/domain/model"
)

func MarshalDate(dt model.Date) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		t := time.Time(dt)
		io.WriteString(w, t.Format(`"2006-01-02"`))
	})
}

func UnmarshalDate(v interface{}) (model.Date, error) {
	if str, ok := v.(string); !ok || (len(str) != 12) {
		return model.Date{}, errors.New("invalid date format")
	} else {
		if t, err := time.ParseInLocation(`"2006-01-02"`, str, time.UTC); err != nil {
			return model.Date{}, err
		} else {
			return model.Date(t), nil
		}
	}
}

func MarshalDateTime(dt model.DateTime) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		t := time.Time(dt)
		io.WriteString(w, t.Format(`"2006-01-02T15:04:05.999Z"`))
	})
}

func UnmarshalDateTime(v interface{}) (model.DateTime, error) {
	if str, ok := v.(string); !ok || (len(str) < 21) {
		return model.DateTime{}, errors.New("invalid datetime format")
	} else {
		if t, err := time.ParseInLocation(`"2006-01-02T15:04:05.999Z"`, str, time.UTC); err != nil {
			return model.DateTime{}, err
		} else {
			return model.DateTime(t), nil
		}
	}
}
