package service

import (
	"context"
	"github.com/bshyn/go-microservices/account/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, email string, password string) (repository.User, error)
	GetUser(id string) (repository.User, error)
}
