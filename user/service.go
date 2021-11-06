package user

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type service struct {
	data []Profile
}
type AddUserSvc interface {
	AddUser(ctx context.Context, p Profile) error
}
type GetUserSvc interface {
	GetUser(ctx context.Context, userId string) error
}
type Service interface {
	AddUserSvc
	GetUserSvc
}

func NewService() Service {
	return &service{
		data: []Profile{},
	}
}
func (s *service) AddUser(ctx context.Context, p Profile) (err error) {
	fmt.Println("Add Service is called")
	return
}
func (s *service) GetUser(ctx context.Context, userId string) (err error) {
	fmt.Println("Get User Called *********************  ")
	return
}
