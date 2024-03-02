package gql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/executor"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func handleGqlRequest(ctx *fiber.Ctx, r *http.Request, gqlExecutor graphql.GraphExecutor) error {
	requestCtx := r.Context()

	now := graphql.Now()
	params := graphql.RawParams{Headers: r.Header, ReadTime: graphql.TraceTiming{Start: now, End: now}}

	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()

	if err := decoder.Decode(&params); err != nil {
		gqlErr := gqlerror.Errorf("request decode: %+v; body: `%s`", err, ctx.BodyRaw())
		response := gqlExecutor.DispatchError(requestCtx, gqlerror.List{gqlErr})
		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if operation, err := gqlExecutor.CreateOperationContext(requestCtx, &params); err != nil {
		response := gqlExecutor.DispatchError(graphql.WithOperationContext(requestCtx, operation), err)
		return ctx.Status(http.StatusUnprocessableEntity).JSON(response)
	} else {
		responseHandler, responseCtx := gqlExecutor.DispatchOperation(requestCtx, operation)
		response := responseHandler(responseCtx)
		if len(response.Errors) > 0 {
			ctx.Status(http.StatusBadRequest)
		}
		return ctx.JSON(response)
	}
}

func GraphQL(gqlExecutor *executor.Executor) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if request, err := adaptor.ConvertRequest(ctx, true); err != nil {
			return err
		} else {
			requestCtx := graphql.StartOperationTrace(request.Context())
			return handleGqlRequest(ctx, request.WithContext(requestCtx), gqlExecutor)
		}
	}
}
