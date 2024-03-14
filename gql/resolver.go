package gql

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"graphql-project/config"
	"graphql-project/core"
	"graphql-project/domain/repository"
)

type Resolver struct {
	cfg               *config.Config
	orderRepository   *repository.OrderRepository
	productRepository *repository.ProductRepository
	reviewRepository  *repository.ReviewRepository
	userRepository    *repository.UserRepository
}

func NewResolver(cfg *config.Config, orderRepository *repository.OrderRepository, productRepository *repository.ProductRepository, reviewRepository *repository.ReviewRepository, userRepository *repository.UserRepository) Resolver {
	return Resolver{cfg, orderRepository, productRepository, reviewRepository, userRepository}
}

func (r *Resolver) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := r.userRepository.GetUserByEmail(ctx, email)
	if err != nil || user == nil || user.Password != password {
		return "", fiber.ErrUnauthorized
	}
	if token, err := core.NewJwt(user, r.cfg.JwtExpiration(), r.cfg.JwtSecret()); err != nil {
		return "", fiber.ErrInternalServerError
	} else {
		return token, nil
	}
}
