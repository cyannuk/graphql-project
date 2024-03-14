package main

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzerolog"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"

	"graphql-project/config"
	"graphql-project/domain/repository"
	"graphql-project/gql/dataloader"
)

type Application struct {
	app               *fiber.App
	orderRepository   *repository.OrderRepository
	productRepository *repository.ProductRepository
	reviewRepository  *repository.ReviewRepository
	userRepository    *repository.UserRepository
	config            *config.Config
}

func (a *Application) Start() error {
	return a.app.Listen(fmt.Sprintf("%s:%d", a.config.BindAddr(), a.config.Port()))
}

func (a *Application) Shutdown() error {
	return a.app.Shutdown()
}

func Default(_ *fiber.Ctx) error {
	return nil
}

func NewApplication(cfg *config.Config) (Application, error) {
	if err := repository.ApplyMigrations(cfg); err != nil {
		return Application{}, err
	}

	dataSource, err := repository.NewDataSource(cfg)
	if err != nil {
		return Application{}, err
	}

	orderRepository := repository.NewOrderRepository(dataSource)
	productRepository := repository.NewProductRepository(dataSource)
	reviewRepository := repository.NewReviewRepository(dataSource)
	userRepository := repository.NewUserRepository(dataSource)

	gqlExecutor := NewGqlExecutor(cfg, orderRepository, productRepository, reviewRepository, userRepository)

	fiberCfg := fiber.Config{
		JSONEncoder:               json.Marshal,
		JSONDecoder:               json.Unmarshal,
		DisableKeepalive:          true,
		DisableStartupMessage:     true,
		DisableDefaultDate:        true,
		DisableDefaultContentType: true,
	}
	application := Application{fiber.New(fiberCfg), orderRepository, productRepository, reviewRepository, userRepository, cfg}

	application.app.Use(fiberzerolog.New(FiberLogConfig()))
	application.app.Get("/", Default)
	application.app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: cfg.JwtSecret()},
	}))
	application.app.Use(dataloader.New(orderRepository, productRepository, reviewRepository, userRepository))
	application.app.Post("/graphql", GraphQL(gqlExecutor))
	// app.Use(compress.New(compress.Config{ Level: 1 }))

	return application, nil
}
