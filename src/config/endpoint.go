package config

import (
	"context"
	"github.com/bshyn/go-microservices-account/model"
	"github.com/bshyn/go-microservices-account/service"
	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateUser endpoint.Endpoint
	GetUser    endpoint.Endpoint
	Login      endpoint.Endpoint
}

func MakeEndpoints(userService service.UserService, authService service.AuthService, jwtKey []byte) Endpoints {
	var getUserEndpoint endpoint.Endpoint
	{
		kf := func(token *stdjwt.Token) (interface{}, error) { return jwtKey, nil }
		getUserEndpoint = makeGetUserEndpoint(userService)
		getUserEndpoint = jwt.NewParser(kf, stdjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(getUserEndpoint)
	}

	return Endpoints{
		Login:      makeLoginEndpoint(authService),
		CreateUser: makeCreateUserEndpoint(userService),
		GetUser:    getUserEndpoint,
	}
}

func makeCreateUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.CreateUserRequest)
		user, err := s.CreateUser(ctx, req.Email, req.Password)
		return model.CreateUserResponse{Id: user.ID, Email: user.Email}, err
	}
}

func makeGetUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.GetUserRequest)
		user, err := s.GetUser(req.Id)
		return model.GetUserResponse{Id: user.ID, Email: user.Email}, err
	}
}

func makeLoginEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.LoginRequest)
		token, err := s.Login(req.Email, req.Password)
		return model.LoginResponse{Jwt: token}, err
	}
}
