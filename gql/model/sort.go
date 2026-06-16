package model

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/99designs/gqlgen/graphql"

	"graphql-project/domain/model"
)

func MarshalSort(sort model.Sort) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		switch sort {
		case model.Asc:
			io.WriteString(w, "ASC")
		case model.Desc:
			io.WriteString(w, "DESC")
		}
	})
}

func UnmarshalSort(v interface{}) (model.Sort, error) {
	if sort, ok := v.(string); !ok {
		return 0, errors.New("invalid sort")
	} else {
		switch strings.ToLower(sort) {
		case "asc":
			return model.Asc, nil
		case "desc":
			return model.Desc, nil
		default:
			return 0, fmt.Errorf("invalid sort `%s`", sort)
		}
	}
}
