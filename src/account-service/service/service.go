package service

import (
	"context"
	"github.com/bshyn/go-microservices/account/repository"
)

type Service interface {
	CreateUser(ctx context.Context, email string, password string) (repository.User, error)
	GetUser(ctx context.Context, id string) (repository.User, error)
}
