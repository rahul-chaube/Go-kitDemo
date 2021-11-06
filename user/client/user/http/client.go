package userhttpclient

import (
	chttp "Profile/transport/http"
	"Profile/user"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/tracing/opentracing"

	httptransport "github.com/go-kit/kit/transport/http"

	stdopentracing "github.com/opentracing/opentracing-go"
)

func New(instance string, tracer stdopentracing.Tracer, logger log.Logger, client *http.Client) (user.Service, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}
	var opts = func(funcName string) (options []httptransport.ClientOption) {
		options = []httptransport.ClientOption{
			httptransport.ClientBefore(opentracing.HTTPToContext(tracer, funcName, logger)),
			httptransport.SetClient(client),
		}
		return
	}
	addUserEndpoint := httptransport.NewClient(
		http.MethodPost,
		copyURL(u, "/user"),
		chttp.EncodeHTTPGenericRequest,
		user.DecodeAddUserResponse,
		opts("Add User")...,
	)
	getUserEndpoint := httptransport.NewClient(
		http.MethodGet,
		copyURL(u, "/user"),
		chttp.EncodeHTTPGenericRequest,
		user.DecodeGetUserResponse,
		opts("Get User")...,
	)
	return user.Endpoints{
		AddUserEndpoint: user.AddUserEndpoint(addUserEndpoint),
		GetUserEndpoint: user.GetUserEndpoint(getUserEndpoint),
	}, nil
}
func NewWithLB(instancer sd.Instancer, tracer stdopentracing.Tracer, logger log.Logger, retryMax int, retryTimeout time.Duration, client *http.Client) user.Service {
	endpoints := user.Endpoints{}
	{
		factory := newFactory(func(s user.Service) endpoint.Endpoint {
			return user.MakeAddUserEndpoint(s)
		}, tracer, logger, client)
		subcriber := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(subcriber)
		retry := chttp.BalancerToEndpoint(balancer)
		endpoints.AddUserEndpoint = user.AddUserEndpoint(retry)
	}
	{
		factory := newFactory(func(s user.Service) endpoint.Endpoint {
			return user.MakeGetUserEndpoint(s)
		}, tracer, logger, client)
		subcriber := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(subcriber)
		retry := chttp.BalancerToEndpoint(balancer)
		endpoints.GetUserEndpoint = user.GetUserEndpoint(retry)
	}
	return endpoints
}
func newFactory(makeEndpoint func(user.Service) endpoint.Endpoint, tracer stdopentracing.Tracer, logger log.Logger, client *http.Client) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := New(instance, tracer, logger, client)
		if err != nil {
			return nil, nil, err
		}
		endpoint := makeEndpoint(service)
		return endpoint, nil, nil
	}
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}
