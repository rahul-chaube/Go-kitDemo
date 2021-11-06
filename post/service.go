package post

import (
	userService "Profile/user"
	"context"
	"errors"
	"fmt"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type service struct {
	data        []Post
	userService userService.Service
}
type AddPostSvc interface {
	AddPost(ctx context.Context, p Post) error
}
type Service interface {
	AddPostSvc
}

func NewService(userService userService.Service) Service {
	return &service{
		data:        []Post{},
		userService: userService,
	}
}
func (s *service) AddPost(ctx context.Context, p Post) (err error) {
	fmt.Println("Add Post Service is called")
	s.userService.GetUser(ctx, "")
	return
}
