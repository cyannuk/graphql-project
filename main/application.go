package main

import (
	"github.com/99designs/gqlgen/graphql/executor"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/goccy/go-json"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"

	"graphql-pro/domain/repository"
	"graphql-pro/gql"
	"graphql-pro/gql/dataloader"
)

type Application struct {
	app             *fiber.App
	orderRepository *repository.OrderRepository
	userRepository  *repository.UserRepository
	jwtSecret       []byte
}

func (a *Application) Start(bindAddr string) error {
	return a.app.Listen(bindAddr)
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

func NewApplication(connectionString string, jwtSecret []byte) (Application, error) {
	dataSource, err := repository.NewDataSource(connectionString)
	if err != nil {
		return Application{}, err
	}

	orderRepository := repository.NewOrderRepository(dataSource)
	userRepository := repository.NewUserRepository(dataSource)

	resolver := gql.NewResolver(orderRepository, userRepository)
	gqlConfig := gql.Config{Resolvers: &resolver}
	gqlConfig.Complexity.Query.Orders = countComplexity
	gqlConfig.Complexity.Query.Users = countComplexity
	gqlConfig.Complexity.User.Orders = countComplexity
	gqlExecutor := executor.New(gql.NewExecutableSchema(gqlConfig))
	gqlExecutor.Use(extension.FixedComplexityLimit(2000))

	config := fiber.Config{
		JSONEncoder:               json.Marshal,
		JSONDecoder:               json.Unmarshal,
		DisableKeepalive:          true,
		DisableStartupMessage:     true,
		DisableDefaultDate:        true,
		DisableDefaultContentType: true,
	}
	application := Application{fiber.New(config), orderRepository, userRepository, jwtSecret}

	application.app.Get("/", Default)
	application.app.Post("/login", application.Login)
	application.app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: jwtSecret},
	}))
	application.app.Use(dataloader.New(orderRepository, userRepository))
	application.app.Post("/graphql", gql.GraphQL(gqlExecutor))
	// app.Use(compress.New(compress.Config{ Level: 1 }))

	return application, nil
}
