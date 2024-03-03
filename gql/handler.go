package gql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/executor"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/savsgio/gotils/strconv"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

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
		if err, ok := response.Errors[l-1].Err.(*fiber.Error); ok {
			status = err.Code
		} else {
			status = http.StatusBadRequest
		}
	} else {
		status = http.StatusOK
	}
	return status
}

func GraphQL(gqlExecutor *executor.Executor) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		now := graphql.Now()
		params := graphql.RawParams{Headers: requestHeaders(ctx), ReadTime: graphql.TraceTiming{Start: now, End: now}}
		requestCtx := graphql.StartOperationTrace(ctx.Context())

		body := ctx.BodyRaw()
		if err := json.Unmarshal(body, &params); err != nil {
			gqlErr := gqlerror.Errorf("request decode: %+v; body: `%s`", err, body)
			response := gqlExecutor.DispatchError(requestCtx, gqlerror.List{gqlErr})
			return ctx.Status(http.StatusBadRequest).JSON(response)
		}

		if operation, err := gqlExecutor.CreateOperationContext(requestCtx, &params); err != nil {
			response := gqlExecutor.DispatchError(graphql.WithOperationContext(requestCtx, operation), err)
			return ctx.Status(http.StatusUnprocessableEntity).JSON(response)
		} else {
			responseHandler, responseCtx := gqlExecutor.DispatchOperation(requestCtx, operation)
			response := responseHandler(responseCtx)
			return ctx.Status(responseStatus(response)).JSON(response)
		}
	}
}
