package dataloader

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vikstrous/dataloadgen"

	"graphql-project/domain/model"
	"graphql-project/domain/repository"
)

const (
	loadersKey = "data-loaders"
)

type Loaders struct {
	OrderLoader *dataloadgen.Loader[int64, *model.Order]
	UserLoader  *dataloadgen.Loader[int64, *model.User]
}

func NewLoaders(orderRepository *repository.OrderRepository, userRepository *repository.UserRepository) Loaders {
	return Loaders{
		OrderLoader: dataloadgen.NewLoader(orderRepository.GetOrderByIds, dataloadgen.WithWait(time.Millisecond)),
		UserLoader:  dataloadgen.NewLoader(userRepository.GetUserByIds, dataloadgen.WithWait(time.Millisecond)),
	}
}

func New(orderRepository *repository.OrderRepository, userRepository *repository.UserRepository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals(loadersKey, NewLoaders(orderRepository, userRepository))
		return ctx.Next()
	}
}

func FromContext(ctx context.Context) Loaders {
	return ctx.Value(loadersKey).(Loaders)
}

func PrimeOrder(ctx context.Context, order *model.Order) bool {
	loaders := FromContext(ctx)
	return loaders.OrderLoader.Prime(order.ID, order)
}

func PrimeUser(ctx context.Context, user *model.User) bool {
	loaders := FromContext(ctx)
	return loaders.UserLoader.Prime(user.ID, user)
}
