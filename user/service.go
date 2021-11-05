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
type Service interface {
	AddUserSvc
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
