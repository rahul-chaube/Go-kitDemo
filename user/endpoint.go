package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type AddUserEndpoint endpoint.Endpoint
type GetUserEndpoint endpoint.Endpoint

type Endpoints struct {
	AddUserEndpoint
	GetUserEndpoint
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

type getUserRequest struct {
	UserId string `json:"profile"`
}

type getUserResponse struct {
	Err error `json:"error"`
}

func MakeGetUserEndpoint(s GetUserSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getUserRequest)
		err = s.GetUser(ctx, req.UserId)
		return getUserResponse{Err: err}, nil
	}
}

func (e GetUserEndpoint) GetUser(ctx context.Context, userId string) (err error) {
	request := getUserRequest{
		UserId: userId,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	return response.(getUserResponse).Err
}
