package model

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/99designs/gqlgen/graphql"

	"graphql-project/domain/model"
)

func MarshalRole(role model.Role) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		switch role {
		case model.RoleRefresh:
			io.WriteString(w, "REFRESH")
		case model.RoleAnon:
			io.WriteString(w, "ANON")
		case model.RoleUser:
			io.WriteString(w, "USER")
		case model.RoleAdmin:
			io.WriteString(w, "ADMIN")
		}
	})
}

func UnmarshalRole(v interface{}) (model.Role, error) {
	if role, ok := v.(string); !ok {
		return 0, errors.New("invalid role")
	} else {
		switch strings.ToLower(role) {
		case "refresh":
			return model.RoleRefresh, nil
		case "anon":
			return model.RoleAnon, nil
		case "user":
			return model.RoleUser, nil
		case "admin":
			return model.RoleAdmin, nil
		default:
			return 0, fmt.Errorf("invalid role `%s`", role)
		}
	}
}
