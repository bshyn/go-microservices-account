package service

import (
	"context"
	"errors"
	"github.com/bshyn/go-microservices-account/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
	"time"
)

type userService struct {
	repository repository.Repository
	logger     log.Logger
}

func NewUserService(repo repository.Repository, logger log.Logger) UserService {
	return &userService{
		repository: repo,
		logger:     logger,
	}
}

type authService struct {
	signingKey   []byte
	timeToExpire time.Duration
	repository   repository.Repository
	logger       log.Logger
}

func NewAuthService(signingKey []byte, timeToExpire time.Duration, repo repository.Repository, logger log.Logger) AuthService {
	return &authService{
		signingKey:   signingKey,
		timeToExpire: timeToExpire,
		repository:   repo,
		logger:       logger,
	}
}

func (s userService) CreateUser(ctx context.Context, email string, password string) (repository.User, error) {
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

func (s userService) GetUser(id string) (repository.User, error) {
	logger := log.With(s.logger, "method", "GetUser")

	var user repository.User
	var err error

	if user, err = s.repository.GetUser(id); err != nil {
		level.Error(logger).Log("error", err)
		return repository.User{}, err
	}

	return user, nil
}

func (s authService) Login(email string, password string) (string, error) {
	user, err := s.repository.GetUserByEmailAndPassword(email, password)
	if err != nil {
		return "", err
	}

	token, err := generateToken(s.signingKey, s.timeToExpire, user.ID)
	if err != nil {
		return "", errors.New(err.Error())
	}
	return token, nil

}

func generateToken(signingKey []byte, timeToExpire time.Duration, id string) (string, error) {
	claims := customClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * timeToExpire).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
			Subject:   id,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

type customClaims struct {
	jwt.StandardClaims
}
