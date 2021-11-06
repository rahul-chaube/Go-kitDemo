package main

import (
	"Profile/common/ipUtil"
	"Profile/post"
	userHttpClient "Profile/user/client/user/http"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	consulsd "github.com/go-kit/kit/sd/consul"
	stdconsul "github.com/hashicorp/consul/api"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	cfgKey      = "Post-Config"
	serviceName = "Post-Service"
)

var (
	ctx          = context.Background()
	tracer       = stdopentracing.GlobalTracer()
	commit       string
	buildVersion string
	apiVersion   string
)

const (
	retryMax     = 1
	retryTimeout = 10 * time.Second
)

var (
	tags = []string{"/v1.0"}
)

func main() {
	consulConfig := stdconsul.DefaultConfig()
	stdConsulClient, err := stdconsul.NewClient(consulConfig)
	if err != nil {
		panic(err)
	}
	kv, _, err := stdConsulClient.KV().Get(cfgKey, nil)
	if err != nil {
		panic(err)
	}
	if kv == nil {
		panic("Config not found")
	}
	//Service Discovery
	consulClient := consul.NewClient(stdConsulClient)
	eIp, err := ipUtil.ExternalIP()
	if err != nil {
		panic(err)
	}
	checkId := serviceName

	var (
		httpAddr = flag.String("http.addr", ":3002", "HTTP listen address")
	)
	flag.Parse()
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Health Check
	consulRegistar := consul.NewRegistrar(consulClient, &stdconsul.AgentServiceRegistration{
		ID:      checkId,
		Name:    serviceName,
		Tags:    tags,
		Port:    3002,
		Address: eIp,
		Check: &stdconsul.AgentServiceCheck{
			TTL:    "10s",
			Status: stdconsul.HealthPassing,
		},
	}, logger)

	fieldKeys := []string{"method", "error_code"}

	//service Descovery

	dfTransport := http.DefaultTransport.(*http.Transport)
	dfTransport.MaxIdleConnsPerHost = 1000
	dfTransport.MaxIdleConns = 1000
	dfClient := http.DefaultClient
	dfClient.Transport = dfTransport

	userIntancer := consulsd.NewInstancer(consulClient, logger, "User-Service", tags, true)
	userService := userHttpClient.NewWithLB(userIntancer, tracer, log.NewNopLogger(), retryMax, retryTimeout, dfClient)
	var s post.Service
	{
		s = post.NewService(userService)
		s = post.NewLoggingService(logger, s)
		s = post.NewInstrumentingService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "post_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "api",
				Subsystem: "post_service",
				Name:      "request_latency_seconds",
				Help:      "Total duration of requests in seconds.",
			}, fieldKeys), s)
	}
	var h http.Handler
	{
		h = post.MakeHandler(s, log.With(logger, "component", "HTTP"))
	}
	consulRegistar.Register()
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
