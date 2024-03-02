package main

import (
	"errors"
	"io"

	"github.com/99designs/gqlgen/graphql/executor"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"

	"graphql-pro/domain/repository"
	"graphql-pro/gql"
)

func onShutdown(resources ...io.Closer) fiber.OnShutdownHandler {
	return func() error {
		var errs error
		for _, resource := range resources {
			if err := resource.Close(); err != nil {
				errs = errors.Join(errs, err)
			}
		}
		return errs
	}
}

func Default(_ *fiber.Ctx) error {
	return nil
}

func NewApplication(connectionString string) (*fiber.App, error) {
	dataSource, err := repository.NewDataSource(connectionString)
	if err != nil {
		return nil, err
	}

	gqlExecutor := executor.New(gql.NewExecutableSchema(gql.Config{Resolvers: &gql.Resolver{UsersRepo: repository.NewUserRepository(dataSource)}}))

	config := fiber.Config{JSONEncoder: json.Marshal, JSONDecoder: json.Unmarshal, DisableKeepalive: true, DisableStartupMessage: true, DisableDefaultDate: true, DisableDefaultContentType: true}
	app := fiber.New(config)

	app.Get("/", Default)
	app.Post("/graphql", gql.GraphQL(gqlExecutor))
	// app.Hooks().OnShutdown(onShutdown(dataSource))

	return app, nil
}
