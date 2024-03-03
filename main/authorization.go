package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"graphql-pro/core"

	"graphql-pro/domain/model"
)

func (a *Application) Login(ctx *fiber.Ctx) error {
	var login model.Login
	err := ctx.BodyParser(&login)
	if err != nil {
		log.Error().Err(err).Msg("parse request")
		return fiber.ErrBadRequest
	}
	user, err := a.userRepository.GetUserByEmail(ctx.Context(), login.Email)
	if err != nil {
		log.Error().Err(err).Msg("get user by email")
		return fiber.ErrInternalServerError
	}
	if user == nil || user.Password != login.Password {
		return fiber.ErrUnauthorized
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, core.UserClaims(user, time.Now().Add(time.Hour*72)))
	// Generate encoded token and send it as response.
	if t, err := token.SignedString(a.jwtSecret); err != nil {
		log.Error().Err(err).Msg("generate jwt token")
		return fiber.ErrInternalServerError
	} else {
		return ctx.JSON(fiber.Map{"token": t})
	}
}
