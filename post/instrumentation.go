package post

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentingService(requestCount metrics.Counter, requestLatency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   requestCount,
		requestLatency: requestLatency,
		Service:        s,
	}
}
func (s *instrumentingService) AddUser(ctx context.Context, p Post) (err error) {
	fmt.Println("Add Post Instrumentation ")
	defer func(begin time.Time) {
		methodField := []string{"AddPost "}
		s.requestCount.With(methodField...).Add(1)
		s.requestLatency.With(methodField...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.AddPost(ctx, p)
}
