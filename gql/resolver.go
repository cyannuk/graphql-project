package gql

import (
	"graphql-project/domain/repository"
)

type Resolver struct {
	orderRepository *repository.OrderRepository
	userRepository  *repository.UserRepository
}

func NewResolver(orderRepository *repository.OrderRepository, userRepository *repository.UserRepository) Resolver {
	return Resolver{orderRepository, userRepository}
}
