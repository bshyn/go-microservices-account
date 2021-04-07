package model

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type (
	CreateUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	CreateUserResponse struct {
		Id    string `json:"id"`
		Email string `json:"email"`
	}

	GetUserRequest struct {
		Id string `json:"id"`
	}
	GetUserResponse struct {
		Id    string `json:"id"`
		Email string `json:"email"`
	}
)

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func DecodeCreateUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return nil, err
	}

	return req, nil
}

func DecodeGetUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	req := GetUserRequest{Id: vars["id"]}

	return req, nil
}
