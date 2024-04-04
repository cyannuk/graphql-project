package main

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/executor"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/savsgio/gotils/strconv"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"graphql-project/auth"
	"graphql-project/config"
	"graphql-project/core"
	"graphql-project/domain/model"
	"graphql-project/domain/repository"
	"graphql-project/gql"
)

func countComplexity(childComplexity int, _ int32, limit int32) int {
	return int(limit) * childComplexity
}

func hasRole(ctx context.Context, _ interface{}, next graphql.Resolver, roles []model.Role) (interface{}, error) {
	if ok := auth.UserHasRole(ctx, roles); !ok {
		return nil, fiber.ErrForbidden
	}
	return next(ctx)
}

func NewGqlExecutor(cfg *config.Config, orderRepository *repository.OrderRepository, productRepository *repository.ProductRepository,
	reviewRepository *repository.ReviewRepository, userRepository *repository.UserRepository, tokenRepository *repository.TokenRepository,
) *executor.Executor {
	resolver := gql.NewResolver(cfg, orderRepository, productRepository, reviewRepository, userRepository, tokenRepository)
	gqlConfig := gql.Config{Resolvers: &resolver}
	gqlConfig.Directives.HasRole = hasRole
	gqlConfig.Complexity.Query.Orders = countComplexity
	gqlConfig.Complexity.Query.Users = countComplexity
	gqlConfig.Complexity.User.Orders = countComplexity
	gqlExecutor := executor.New(gql.NewExecutableSchema(gqlConfig))
	gqlExecutor.Use(extension.FixedComplexityLimit(cfg.QueryComplexity()))
	return gqlExecutor
}

func requestHeaders(ctx *fiber.Ctx) http.Header {
	headers := make(http.Header, 16)
	ctx.Request().Header.VisitAll(func(k, v []byte) {
		headers.Set(strconv.B2S(k), strconv.B2S(v))
	})
	return headers
}

func responseStatus(response *graphql.Response) int {
	var status int
	l := len(response.Errors)
	if l > 0 {
		err := response.Errors[l-1].Err
		if err == core.ErrNotFound {
			status = http.StatusNotFound
		} else if err, ok := err.(*fiber.Error); ok {
			status = err.Code
		} else {
			status = http.StatusBadRequest
		}
		log.Error().Err(err).Msg("dispatch graphql operation")
	} else {
		status = http.StatusOK
	}
	return status
}

func GraphQL(gqlExecutor *executor.Executor) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		now := graphql.Now()
		params := graphql.RawParams{ /*Headers: requestHeaders(ctx),*/ ReadTime: graphql.TraceTiming{Start: now, End: now}}
		requestCtx := graphql.StartOperationTrace(ctx.Context())

		body := ctx.BodyRaw()
		if err := json.Unmarshal(body, &params); err != nil {
			log.Error().Err(err).Msg("unmarshal graphql request")
			gqlErr := gqlerror.Errorf("request decode: %+v; body: `%s`", err, body)
			response := gqlExecutor.DispatchError(requestCtx, gqlerror.List{gqlErr})
			return ctx.Status(http.StatusBadRequest).JSON(response)
		}

		if operation, err := gqlExecutor.CreateOperationContext(requestCtx, &params); err != nil {
			log.Error().Err(err).Msg("create graphql operation")
			response := gqlExecutor.DispatchError(graphql.WithOperationContext(requestCtx, operation), err)
			return ctx.Status(http.StatusUnprocessableEntity).JSON(response)
		} else {
			responseHandler, responseCtx := gqlExecutor.DispatchOperation(requestCtx, operation)
			response := responseHandler(responseCtx)
			return ctx.Status(responseStatus(response)).JSON(response)
		}
	}
}
