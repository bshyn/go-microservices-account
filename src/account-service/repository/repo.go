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
	query := `INSERT INTO USERS (ID, EMAIL, PASSWORD) VALUES (?, ?, ?)`

	if user.Email == "" || user.Password == "" {
		return EmptyUserErr
	}

	savedUser, err := repo.GetUserByEmail(user.Email)

	if err != nil {
		return err
	}

	if savedUser.ID != "" {
		return UserWithMailExistsErr
	}

	_, err = repo.db.ExecContext(ctx, query, user.ID, user.Email, user.Password)

	if err != nil {
		return err
	}
	return nil
}

func (repo *repo) GetUser(id string) (User, error) {
	query := `SELECT EMAIL, PASSWORD FROM USERS WHERE ID = ?`

	var email string
	var password string

	err := repo.db.QueryRow(query, id).Scan(&email, &password)
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

func (repo *repo) GetUserByEmail(email string) (User, error) {
	query := `SELECT ID, PASSWORD FROM USERS WHERE EMAIL = ?`

	var id string
	var password string

	err := repo.db.QueryRow(query, email).Scan(&id, &password)
	if err != sql.ErrNoRows && err != nil {
		return User{}, RepoErr
	}

	user := User{
		ID:       id,
		Email:    email,
		Password: password,
	}

	return user, nil
}
