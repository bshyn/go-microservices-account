package service

import (
	"context"
	"github.com/bshyn/go-microservices/account/repository"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
)

type service struct {
	repository repository.Repository
	logger     log.Logger
}

func NewService(repo repository.Repository, logger log.Logger) Service {
	return &service{
		repository: repo,
		logger:     logger,
	}
}

func (s service) CreateUser(ctx context.Context, email string, password string) (repository.User, error) {
	logger := log.With(s.logger, "method", "CreateUser")

	uuid, _ := uuid.NewV4()
	id := uuid.String()

	user := repository.User{
		ID:       id,
		Email:    email,
		Password: password,
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		level.Error(logger).Log("error", err)
		return repository.User{}, err
	}

	logger.Log("newUser", id)

	return user, nil
}

func (s service) GetUser(id string) (repository.User, error) {
	logger := log.With(s.logger, "method", "GetUser")

	var user repository.User
	var err error

	if user, err = s.repository.GetUser(id); err != nil {
		level.Error(logger).Log("error", err)
		return repository.User{}, err
	}

	return user, nil
}
