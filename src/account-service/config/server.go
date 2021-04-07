package config

import (
	"context"
	"github.com/bshyn/go-microservices/account/model"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	router := mux.NewRouter()
	router.Use(commonMiddleware)

	router.Methods("POST").Path("/user").Handler(httptransport.NewServer(
		endpoints.CreateUser,
		model.DecodeCreateUserReq,
		model.EncodeResponse,
	))

	router.Methods("GET").Path("/user/{id}").Handler(httptransport.NewServer(
		endpoints.GetUser,
		model.DecodeGetUserReq,
		model.EncodeResponse,
	))

	return router
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
