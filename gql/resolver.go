package gql

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"graphql-project/config"
	"graphql-project/core"
	"graphql-project/domain/model"
	"graphql-project/domain/repository"
)

type Resolver struct {
	cfg               *config.Config
	orderRepository   *repository.OrderRepository
	productRepository *repository.ProductRepository
	reviewRepository  *repository.ReviewRepository
	userRepository    *repository.UserRepository
	tokenRepository   *repository.TokenRepository
}

func NewResolver(cfg *config.Config, orderRepository *repository.OrderRepository, productRepository *repository.ProductRepository,
	reviewRepository *repository.ReviewRepository, userRepository *repository.UserRepository, tokenRepository *repository.TokenRepository,
) Resolver {
	return Resolver{
		cfg,
		orderRepository,
		productRepository,
		reviewRepository,
		userRepository,
		tokenRepository,
	}
}

func (r *Resolver) Login(ctx context.Context, email string, password string) (tokens model.Tokens, err error) {
	user, err := r.userRepository.GetUserByEmail(repository.With(ctx, 0, 0, "id", "password", "email", "name", "role"), email)
	if err != nil || user == nil || !core.VerifyPassword(password, user.Password) {
		err = fiber.ErrUnauthorized
		return
	}
	tokens, err = core.NewJwt(user, r.cfg.JwtExpiration(), r.cfg.JwtRefreshExpiration(), r.cfg.JwtSecret())
	if err == nil {
		err = r.tokenRepository.CreateToken(ctx, user.ID, tokens.RefreshToken)
	}
	return
}

func (r *Resolver) Refresh(ctx context.Context) (tokens model.Tokens, err error) {
	userId, _ := core.GetContextUser(ctx)
	user, err := r.userRepository.GetUserByID(repository.With(ctx, 0, 0, "id", "password", "email", "name", "role"), userId)
	if err != nil || user == nil {
		err = fiber.ErrUnauthorized
		return
	}
	refreshToken, err := r.tokenRepository.GetTokenByID(ctx, userId)
	if err != nil || refreshToken != core.GetJwt(ctx) {
		err = fiber.ErrUnauthorized
		return
	}
	tokens, err = core.NewJwt(user, r.cfg.JwtExpiration(), r.cfg.JwtRefreshExpiration(), r.cfg.JwtSecret())
	if err == nil {
		err = r.tokenRepository.CreateToken(ctx, user.ID, tokens.RefreshToken)
	}
	return
}
