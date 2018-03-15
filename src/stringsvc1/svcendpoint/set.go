package svcendpoint

import (
	"context"
	"time"

	"golang.org/x/time/rate"

	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"

	"stringsvc1/service"
	"fmt"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	UppercaseEndpoint    endpoint.Endpoint
	CreateEndpoint    endpoint.Endpoint
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(svc service.Service, trace stdopentracing.Tracer) Set {
	var uppercaseEndpoint endpoint.Endpoint
	{
		uppercaseEndpoint = MakeUppercaseEndpoint(svc)
		uppercaseEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(uppercaseEndpoint)
		uppercaseEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(uppercaseEndpoint)
		uppercaseEndpoint = opentracing.TraceServer(trace, "Uppercase")(uppercaseEndpoint)
	}
	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = MakeCreateEndpoint(svc)
		createEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(createEndpoint)
		createEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(createEndpoint)
		createEndpoint = opentracing.TraceServer(trace, "Uppercase")(createEndpoint)
	}
	return Set{
		UppercaseEndpoint:    uppercaseEndpoint,
		CreateEndpoint:    createEndpoint,
	}
}

// Sum implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) Uppercase(ctx context.Context, a string) (string, error) {
	resp, err := s.UppercaseEndpoint(ctx, UppercaseRequest{A: a})
	if err != nil {
		return "", err
	}
	response := resp.(UppercaseResponse)
	return response.V, response.Err
}

func (s Set) Create(ctx context.Context, ID string, FlowID uint32, Source string, Type string) (string, error){
	fmt.Print("create alarm data")
	resp, err := s.CreateEndpoint(ctx, CreateRequest{ID: ID,FlowID: FlowID,Source: Source,Type: Type})
	if err != nil {
		return "", err
	}
	response := resp.(CreateResponse)
	return response.V, response.Err
}

// Concat implements the service interface, so Set may be used as a
// service. This is primarily useful in the context of a client library.


// MakeSumEndpoint constructs a Sum endpoint wrapping the service.
func MakeUppercaseEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UppercaseRequest)
		v, err := s.Uppercase(ctx, req.A)
		return UppercaseResponse{V: v, Err: err}, nil
	}
}

func MakeCreateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateRequest)
		err = s.Create(ctx, req.ID, req.FlowID, req.Source, req.Type)
		return CreateResponse{V:"create alarm data", Err: err}, nil
	}
}

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}

// SumRequest collects the request parameters for the Sum method.
type UppercaseRequest struct {
	A string
}

// SumResponse collects the response values for the Sum method.
type UppercaseResponse struct {
	V   string   `json:"v"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

// Failed implements Failer.
func (r UppercaseResponse) Failed() error { return r.Err }

type CreateRequest struct {
	ID string `json:"ID"`
	FlowID uint32 `json:"FlowID"`
	Source string `json:"Source"`
	Type string `json:"Type"`
}

// SumResponse collects the response values for the Sum method.
type CreateResponse struct {
	V   string   `json:"v"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

// Failed implements Failer.
func (r CreateResponse) Failed() error { return r.Err }
