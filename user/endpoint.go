package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type AddUserEndpoint endpoint.Endpoint

type Endpoint struct {
	AddUserEndpoint
}

type addUserRequest struct {
	UserProfile Profile `json:"profile"`
}

type addUserResponse struct {
	Err error `json:"profile"`
}

func MakeAddUserEndpoint(s AddUserSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addUserRequest)
		err = s.AddUser(ctx, req.UserProfile)
		return addUserResponse{Err: err}, nil
	}
}

func (e AddUserEndpoint) AddUser(ctx context.Context, profile Profile) (err error) {
	request := addUserRequest{
		UserProfile: profile,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	return response.(addUserResponse).Err
}
