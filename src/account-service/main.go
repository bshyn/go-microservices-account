package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bshyn/go-microservices/account/config"
	"github.com/bshyn/go-microservices/account/repository"
	"github.com/bshyn/go-microservices/account/service"

	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"
)

func main() {
	var userLogger log.Logger
	{
		userLogger = log.NewLogfmtLogger(os.Stderr)
		userLogger = log.NewSyncLogger(userLogger)
		userLogger = log.With(userLogger,
			"logger", "userLogger",
			"service", "account",
			"time", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}
	var authLogger log.Logger
	{
		authLogger = log.NewLogfmtLogger(os.Stdout)
		authLogger = log.NewSyncLogger(authLogger)
		authLogger = log.With(authLogger,
			"loggerr", "authLogger",
			"service", "account",
			"time", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(userLogger).Log("msg", "service started")
	defer level.Info(userLogger).Log("msg", "service stopped")

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbSource := fmt.Sprintf("%s:%s@tcp(%s)/account", dbUser, dbPassword, dbHost)

	port := os.Getenv("PORT")

	var jwtExpiration time.Duration
	{
		jwtExpirationStr := os.Getenv("JWT_EXPIRATION")
		duration, err := time.ParseDuration(jwtExpirationStr)
		if err != nil {
			panic(err)
		}
		jwtExpiration = duration
	}

	var jwtKey []byte
	{
		jwtKeyStr := os.Getenv("JWT_KEY")
		jwtKey = []byte(jwtKeyStr)
	}

	var db *sql.DB
	{
		var err error

		db, err = sql.Open("mysql", dbSource)
		if err != nil {
			level.Error(userLogger).Log("exit", err)
			os.Exit(-1)
		}
	}

	ctx := context.Background()
	var userSrv service.UserService
	{
		repository := repository.NewRepo(db, userLogger)
		userSrv = service.NewUserService(repository, userLogger)
	}
	var authSrv service.AuthService
	{
		repository := repository.NewRepo(db, authLogger)
		authSrv = service.NewAuthService(jwtKey, jwtExpiration, repository, authLogger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := config.MakeEndpoints(userSrv, authSrv, jwtKey)

	go func() {
		fmt.Println("listening on port", port)
		handler := config.NewHTTPServer(ctx, endpoints)
		errs <- http.ListenAndServe(":"+port, handler)
	}()

	level.Error(userLogger).Log("exit", <-errs)
}
