package main

import (
	profilesvc "Profile/user"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	"github.com/go-kit/log"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	fieldKeys := []string{"method", "error_code"}

	var s profilesvc.Service
	{
		s = profilesvc.NewService()
		s = profilesvc.NewLoggingService(logger, s)
		s = profilesvc.NewInstrumentingService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "appexpiry_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "api",
				Subsystem: "appexpiry_service",
				Name:      "request_latency_seconds",
				Help:      "Total duration of requests in seconds.",
			}, fieldKeys), s)
	}
	var h http.Handler
	{
		h = profilesvc.MakeHandler(s, log.With(logger, "component", "HTTP"))
	}
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
