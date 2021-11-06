package post

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) AddUser(ctx context.Context, p Post) (err error) {
	defer func(begin time.Time) {
		s.logger.Log("mathod", "Add User",
			"userProfile", p,
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return s.Service.AddPost(ctx, p)
}
