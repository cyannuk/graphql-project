package core

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"graphql-pro/domain/model"
)

func GetUserSession(ctx context.Context) model.UserSession {
	user := ctx.Value("user")
	if user != nil {
		if token, ok := user.(*jwt.Token); ok {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				return model.UserSession{
					Name:   claims["name"].(string),
					Email:  claims["email"].(string),
					UserId: int64(claims["user_id"].(float64)),
					Admin:  claims["admin"].(bool),
					Exp:    time.Unix(int64(claims["exp"].(float64)), 0),
				}
			}
		}
	}
	return model.UserSession{}
}

func UserClaims(user *model.User, expiration time.Time) jwt.MapClaims {
	return jwt.MapClaims{
		"name":    user.Name,
		"email":   user.Email,
		"user_id": user.ID,
		"admin":   false,
		"exp":     expiration.Unix(),
	}
}

func CheckUserId(ctx context.Context, id int64) (int64, error) {
	userSession := GetUserSession(ctx)
	if id > 0 {
		if !userSession.Admin && id != userSession.UserId {
			return 0, fiber.ErrForbidden
		}
	} else {
		id = userSession.UserId
	}
	return id, nil
}

func CheckUserEmail(ctx context.Context, email string) (string, error) {
	userSession := GetUserSession(ctx)
	if len(email) > 0 {
		if !userSession.Admin && email != userSession.Email {
			return "", fiber.ErrForbidden
		}
	} else {
		email = userSession.Email
	}
	return email, nil
}
