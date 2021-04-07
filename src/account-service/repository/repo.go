package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-kit/kit/log"
)

var (
	RepoErr               = errors.New("unable to handle Repo Request")
	EmptyUserErr          = errors.New("password and email are required")
	UserWithMailExistsErr = errors.New("an user with this mail already exists")
)

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepo(db *sql.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}

func (repo *repo) CreateUser(ctx context.Context, user User) error {
	sql := `INSERT INTO USERS (ID, EMAIL, PASSWORD) VALUES (?, ?, ?)`

	if user.Email == "" || user.Password == "" {
		return EmptyUserErr
	}

	savedUser, err := repo.GetUserByEmail(ctx, user.Email)

	if err != nil {
		return err
	}

	if savedUser.ID != "" {
		return UserWithMailExistsErr
	}

	_, err = repo.db.ExecContext(ctx, sql, user.ID, user.Email, user.Password)

	if err != nil {
		return err
	}
	return nil
}

func (repo *repo) GetUser(ctx context.Context, id string) (User, error) {
	sql := `SELECT EMAIL, PASSWORD FROM USERS WHERE ID = ?`

	var email string
	var password string

	err := repo.db.QueryRow(sql, id).Scan(&email, &password)
	if err != nil {
		return User{}, RepoErr
	}

	user := User{
		ID:       id,
		Email:    email,
		Password: password,
	}

	return user, nil
}

func (repo *repo) GetUserByEmail(ctx context.Context, email string) (User, error) {
	sql := `SELECT ID, PASSWORD FROM USERS WHERE EMAIL = ?`

	var id string
	var password string

	err := repo.db.QueryRow(sql, email).Scan(&id, &password)
	if err != nil {
		return User{}, RepoErr
	}

	user := User{
		ID:       id,
		Email:    email,
		Password: password,
	}

	return user, nil
}
