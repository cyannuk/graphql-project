package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzerolog"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"

	"graphql-project/config"
	"graphql-project/domain/repository"
	"graphql-project/gql/dataloader"
	"graphql-project/tracing"
)

type Application struct {
	app               *fiber.App
	tracerProvider    tracing.TracerProvider
	orderRepository   *repository.OrderRepository
	productRepository *repository.ProductRepository
	reviewRepository  *repository.ReviewRepository
	userRepository    *repository.UserRepository
	tokenRepository   *repository.TokenRepository
	config            *config.Config
}

func (a *Application) Start() error {
	return a.app.Listen(fmt.Sprintf("%s:%d", a.config.BindAddr(), a.config.Port()))
}

func (a *Application) Shutdown() error {
	ctx := context.Background()
	return errors.Join(a.tracerProvider.Shutdown(ctx), a.app.ShutdownWithContext(ctx))
}

func Default(_ *fiber.Ctx) error {
	return nil
}

func handleJwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": []fiber.Map{{"message": "Invalid or expired JWT"}}, "data": nil})
}

func NewApplication(cfg *config.Config) (application Application, err error) {
	tracerProvider, err := tracing.InitProviders(cfg)
	if err != nil {
		return
	}

	err = repository.ApplyMigrations(cfg)
	if err != nil {
		return
	}

	dataSource, err := repository.NewDataSource(cfg)
	if err != nil {
		return
	}

	orderRepository := repository.NewOrderRepository(dataSource)
	productRepository := repository.NewProductRepository(dataSource)
	reviewRepository := repository.NewReviewRepository(dataSource)
	userRepository := repository.NewUserRepository(dataSource)
	tokenRepository := repository.NewTokenRepository(dataSource)

	gqlExecutor := NewGqlExecutor(cfg, orderRepository, productRepository, reviewRepository, userRepository, tokenRepository)

	fiberCfg := fiber.Config{
		JSONEncoder:               json.Marshal,
		JSONDecoder:               json.Unmarshal,
		DisableKeepalive:          true,
		DisableStartupMessage:     true,
		DisableDefaultDate:        true,
		DisableDefaultContentType: true,
		DisableHeaderNormalizing:  true,
	}
	application = Application{
		fiber.New(fiberCfg),
		tracerProvider,
		orderRepository,
		productRepository,
		reviewRepository,
		userRepository,
		tokenRepository,
		cfg,
	}

	application.app.Use(fiberzerolog.New(FiberLogConfig()))
	// application.app.Use(tracing.Middleware(tracerProvider))
	application.app.Get("/", Default)
	application.app.Use(jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: cfg.JwtSecret()},
		ErrorHandler: handleJwtError,
	}))
	application.app.Use(dataloader.New(orderRepository, productRepository, reviewRepository, userRepository))
	application.app.Post("/graphql", GraphQL(gqlExecutor))
	// app.Use(compress.New(compress.Config{ Level: 1 }))

	err = nil
	return
}
