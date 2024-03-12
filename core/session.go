package core

import (
	"context"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"graphql-project/domain/model"
)

const ctxUserKey = "user"

func NewJwt(user *model.User, expiration time.Duration, jwtSecret []byte) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"name":  user.Name,
		"email": user.Email,
		"uid":   user.ID,
		"role":  user.Role,
		"iat":   now.Unix(),
		"exp":   now.Add(expiration * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func getString(claims jwt.MapClaims, key string) string {
	if claim, ok := claims[key]; ok {
		if s, ok := claim.(string); ok {
			return s
		}
	}
	return ""
}

func getInt(claims jwt.MapClaims, key string) int64 {
	if claim, ok := claims[key]; ok {
		if v, ok := claim.(float64); ok {
			return int64(v)
		}
	}
	return 0
}

func getRole(claims jwt.MapClaims, key string) model.Role {
	role := model.RoleAnon
	if claim, ok := claims[key]; ok {
		if v, ok := claim.(float64); ok {
			role = model.Role(int64(v))
		}
	}
	return role
}

func getTime(claims jwt.MapClaims, key string) time.Time {
	var t int64 = 0
	if claim, ok := claims[key]; ok {
		if v, ok := claim.(float64); ok {
			t = int64(v)
		}
	}
	return time.Unix(t, 0)
}

func getJwtClaims(ctx context.Context, key string) jwt.MapClaims {
	user := ctx.Value(key)
	if user != nil {
		if token, ok := user.(*jwt.Token); ok {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				return claims
			}
		}
	}
	return nil
}

func GetContextUser(ctx context.Context) (int64, model.Role) {
	if claims := getJwtClaims(ctx, ctxUserKey); claims != nil {
		return getInt(claims, "uid"), getRole(claims, "role")
	}
	return -1, model.RoleAnon
}

func CheckUserId(ctx context.Context, id int64) (int64, bool) {
	if claims := getJwtClaims(ctx, ctxUserKey); claims != nil {
		userId := getInt(claims, "uid")
		role := getRole(claims, "role")
		if id > 0 {
			if id == userId || role == model.RoleAdmin {
				return id, true
			}
		} else {
			return userId, true
		}
	}
	return 0, false
}

func CheckUserEmail(ctx context.Context, email string) (string, bool) {
	if claims := getJwtClaims(ctx, ctxUserKey); claims != nil {
		userEmail := getString(claims, "email")
		role := getRole(claims, "role")
		if len(email) > 0 {
			if email == userEmail || role == model.RoleAdmin {
				return email, true
			}
		} else {
			return userEmail, true
		}
	}
	return "", false
}

func UserHasRole(ctx context.Context, roles []model.Role) bool {
	if claims := getJwtClaims(ctx, ctxUserKey); claims != nil {
		role := getRole(claims, "role")
		return slices.Contains(roles, role)
	}
	return false
}
