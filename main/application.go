package main

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql/executor"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzerolog"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"

	"graphql-project/config"
	"graphql-project/domain/repository"
	"graphql-project/gql"
	"graphql-project/gql/dataloader"
)

type Application struct {
	app             *fiber.App
	orderRepository *repository.OrderRepository
	userRepository  *repository.UserRepository
	config          *config.Config
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

func countComplexity(childComplexity int, _ int32, limit int32) int {
	return int(limit) * childComplexity
}

func NewGqlExecutor(orderRepository *repository.OrderRepository, userRepository *repository.UserRepository) *executor.Executor {
	resolver := gql.NewResolver(orderRepository, userRepository)
	gqlConfig := gql.Config{Resolvers: &resolver}
	gqlConfig.Complexity.Query.Orders = countComplexity
	gqlConfig.Complexity.Query.Users = countComplexity
	gqlConfig.Complexity.User.Orders = countComplexity
	return executor.New(gql.NewExecutableSchema(gqlConfig))
}

func NewApplication(config *config.Config) (Application, error) {
	if err := repository.ApplyMigrations(config); err != nil {
		return Application{}, err
	}

	dataSource, err := repository.NewDataSource(config)
	if err != nil {
		return Application{}, err
	}

	orderRepository := repository.NewOrderRepository(dataSource)
	userRepository := repository.NewUserRepository(dataSource)

	gqlExecutor := NewGqlExecutor(orderRepository, userRepository)
	gqlExecutor.Use(extension.FixedComplexityLimit(config.QueryComplexity()))

	fiberCfg := fiber.Config{
		JSONEncoder:               json.Marshal,
		JSONDecoder:               json.Unmarshal,
		DisableKeepalive:          true,
		DisableStartupMessage:     true,
		DisableDefaultDate:        true,
		DisableDefaultContentType: true,
	}
	application := Application{fiber.New(fiberCfg), orderRepository, userRepository, config}

	application.app.Use(fiberzerolog.New(FiberLogConfig()))
	application.app.Get("/", Default)
	application.app.Post("/login", application.Login)
	application.app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: config.JwtSecret()},
	}))
	application.app.Use(dataloader.New(orderRepository, userRepository))
	application.app.Post("/graphql", gql.GraphQL(gqlExecutor))
	// app.Use(compress.New(compress.Config{ Level: 1 }))

	return application, nil
}
