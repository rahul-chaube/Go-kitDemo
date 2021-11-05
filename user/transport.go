package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	ErrBadRequest = errors.New("Bad Request")
)
var errBadRoute = errors.New("bad route")

func MakeHandler(s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}
	AddUserHandler := kithttp.NewServer(
		MakeAddUserEndpoint(s),
		DecodeAddUserRequest,
		encodeResponse,
		opts...,
	)
	r := mux.NewRouter()
	r.Handle("/user", AddUserHandler)
	return r
}
func DecodeAddUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req addUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}
	return req, err
}

// func DecodeAddUserResponse(ctx context.Context, r *http.Response) (interface{}, error) {
// 	var resp addUserResponse
// 	err := chttp.DecodeResponse(ctx, r, &resp)
// 	return resp, err
// }

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	// case cargo.ErrUnknown:
	// 	w.WriteHeader(http.StatusNotFound)
	case ErrBadRequest:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
