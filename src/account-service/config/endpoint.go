package config

import (
	"context"
	"github.com/bshyn/go-microservices/account/model"
	"github.com/bshyn/go-microservices/account/service"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateUser endpoint.Endpoint
	GetUser    endpoint.Endpoint
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		CreateUser: makeCreateUserEndpoint(s),
		GetUser:    makeGetUserEndpoint(s),
	}
}

func makeCreateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.CreateUserRequest)
		user, err := s.CreateUser(ctx, req.Email, req.Password)
		return model.CreateUserResponse{Id: user.ID, Email: user.Email}, err
	}
}

func makeGetUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.GetUserRequest)
		user, err := s.GetUser(req.Id)
		return model.GetUserResponse{Id: user.ID, Email: user.Email}, err
	}
}
