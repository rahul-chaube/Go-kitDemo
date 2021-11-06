package post

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type AddPostEndpoint endpoint.Endpoint

type Endpoint struct {
	AddPostEndpoint
}

type addPostRequest struct {
	Post Post `json:"profile"`
}

type addPostResponse struct {
	Err error `json:"profile"`
}

func MakeAddPostEndpoint(s AddPostSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addPostRequest)
		err = s.AddPost(ctx, req.Post)
		return addPostResponse{Err: err}, nil
	}
}

func (e AddPostEndpoint) AddUser(ctx context.Context, profile Post) (err error) {
	request := addPostRequest{
		Post: profile,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	return response.(addPostResponse).Err
}
