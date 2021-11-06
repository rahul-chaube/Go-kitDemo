package main

import (
	ipUtil "Profile/common/ipUtil"
	profilesvc "Profile/user"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	stdconsul "github.com/hashicorp/consul/api"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/sd/consul"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	"github.com/go-kit/log"
)

const (
	cfgKey      = "User-Config"
	serviceName = "User-Service"
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
	// //Config
	// config, err := configAppExpiry.New(bytes.NewBuffer(kv.Value), true)
	// if err != nil {
	// 	panic(err)
	// }

	//Service Discovery
	consulClient := consul.NewClient(stdConsulClient)
	eIp, err := ipUtil.ExternalIP()
	if err != nil {
		panic(err)
	}
	checkId := serviceName

	var (
		httpAddr = flag.String("http.addr", ":3001", "HTTP listen address")
	)
	flag.Parse()
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	fieldKeys := []string{"method", "error_code"}

	consulRegistar := consul.NewRegistrar(consulClient, &stdconsul.AgentServiceRegistration{
		ID:      checkId,
		Name:    serviceName,
		Tags:    []string{"/V1.0"},
		Port:    3001,
		Address: eIp,
		Check: &stdconsul.AgentServiceCheck{
			TTL:    "10s",
			Status: stdconsul.HealthPassing,
		},
	}, logger)

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
