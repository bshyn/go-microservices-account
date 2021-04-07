package main

import (
	"context"
	"database/sql"
	"fmt"

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
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "account",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service stopped")


	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dbSource := fmt.Sprintf("%s:%s@tcp(%s)/account",  dbUser, dbPassword, dbHost)

	level.Debug(logger).Log("dbSource", dbSource)

	port := os.Getenv("PORT")

	var db *sql.DB
	{
		var err error

		db, err = sql.Open("mysql", dbSource)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	ctx := context.Background()
	var srv service.Service
	{
		repository := repository.NewRepo(db, logger)
		srv = service.NewService(repository, logger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := config.MakeEndpoints(srv)

	go func() {
		fmt.Println("listening on port", port)
		handler := config.NewHTTPServer(ctx, endpoints)
		errs <- http.ListenAndServe(":" + port, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
}
