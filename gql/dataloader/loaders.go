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
	OrderLoader   *dataloadgen.Loader[int64, *model.Order]
	ProductLoader *dataloadgen.Loader[int64, *model.Product]
	ReviewLoader  *dataloadgen.Loader[int64, *model.Review]
	UserLoader    *dataloadgen.Loader[int64, *model.User]
}

func NewLoaders(orderRepository *repository.OrderRepository, productRepository *repository.ProductRepository, reviewRepository *repository.ReviewRepository, userRepository *repository.UserRepository) Loaders {
	return Loaders{
		OrderLoader:   dataloadgen.NewLoader(orderRepository.GetOrderByIds, dataloadgen.WithWait(time.Millisecond)),
		ProductLoader: dataloadgen.NewLoader(productRepository.GetProductByIds, dataloadgen.WithWait(time.Millisecond)),
		ReviewLoader:  dataloadgen.NewLoader(reviewRepository.GetReviewByIds, dataloadgen.WithWait(time.Millisecond)),
		UserLoader:    dataloadgen.NewLoader(userRepository.GetUserByIds, dataloadgen.WithWait(time.Millisecond)),
	}
}

func New(orderRepository *repository.OrderRepository, productRepository *repository.ProductRepository, reviewRepository *repository.ReviewRepository, userRepository *repository.UserRepository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals(loadersKey, NewLoaders(orderRepository, productRepository, reviewRepository, userRepository))
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

func PrimeProduct(ctx context.Context, product *model.Product) bool {
	loaders := FromContext(ctx)
	return loaders.ProductLoader.Prime(product.ID, product)
}

func PrimeReview(ctx context.Context, review *model.Review) bool {
	loaders := FromContext(ctx)
	return loaders.ReviewLoader.Prime(review.ID, review)
}

func PrimeUser(ctx context.Context, user *model.User) bool {
	loaders := FromContext(ctx)
	return loaders.UserLoader.Prime(user.ID, user)
}
