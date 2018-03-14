package alarmendpoint

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

	"alarm/alarmservice"
	"fmt"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	CreateEndpoint    endpoint.Endpoint
	AddEndpoint       endpoint.Endpoint
	EndEndpoint       endpoint.Endpoint

}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(svc alarmservice.Service, trace stdopentracing.Tracer) Set {
	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = MakeCreateEndpoint(svc)
		createEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(createEndpoint)
		createEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(createEndpoint)
		createEndpoint = opentracing.TraceServer(trace, "Create")(createEndpoint)
	}
	var addEndpoint endpoint.Endpoint
	{
		addEndpoint = MakeAddEndpoint(svc)
		addEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(addEndpoint)
		addEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(addEndpoint)
		addEndpoint = opentracing.TraceServer(trace, "Add")(addEndpoint)
	}
	var endEndpoint endpoint.Endpoint
	{
		endEndpoint = MakeEndEndpoint(svc)
		endEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(endEndpoint)
		endEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(endEndpoint)
		endEndpoint = opentracing.TraceServer(trace, "End")(endEndpoint)
	}
	return Set{
		CreateEndpoint:    createEndpoint,
		AddEndpoint:       addEndpoint,
		EndEndpoint:       endEndpoint,
	}
}


func (s Set) Create(ctx context.Context, ID string, FlowID uint32, Source string, Type string, Strategy string, Target string, SourceID string) (string, error){
	fmt.Print("create alarm data")
	resp, err := s.CreateEndpoint(ctx, CreateRequest{ID: ID,FlowID: FlowID,Source: Source,Type: Type,Strategy: Strategy,Target: Target,SourceID: SourceID})
	if err != nil {
		return "", err
	}
	response := resp.(CreateResponse)
	return response.V, response.Err
}

func (s Set) Add(ctx context.Context, ID string, FlowID uint32, Source string, Type string, Strategy string, Target string, SourceID string) (string, error){
	fmt.Print("add alarm data")
	resp, err := s.AddEndpoint(ctx, AddRequest{ID: ID,FlowID: FlowID,Source: Source,Type: Type,Strategy: Strategy,Target: Target,SourceID: SourceID})
	if err != nil {
		return "", err
	}
	response := resp.(AddResponse)
	return response.V, response.Err
}

func (s Set) End(ctx context.Context, ID string, FlowID uint32, Source string, Type string, Strategy string, Target string, SourceID string) (string, error){
	fmt.Print("end alarm data")
	resp, err := s.EndEndpoint(ctx, AddRequest{ID: ID,FlowID: FlowID,Source: Source,Type: Type,Strategy: Strategy,Target: Target,SourceID: SourceID})
	if err != nil {
		return "", err
	}
	response := resp.(EndResponse)
	return response.V, response.Err
}

func MakeCreateEndpoint(s alarmservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateRequest)
		err = s.Create(ctx, req.ID, req.FlowID, req.Source, req.Type,req.Strategy,req.Target,req.SourceID)
		return CreateResponse{V:"create alarm data", Err: err}, nil
	}
}

func MakeAddEndpoint(s alarmservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddRequest)
		err = s.Add(ctx, req.ID, req.FlowID, req.Source, req.Type,req.Strategy,req.Target,req.SourceID)
		return AddResponse{V:"add alarm data", Err: err}, nil
	}
}

func MakeEndEndpoint(s alarmservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(EndRequest)
		err = s.End(ctx, req.ID, req.FlowID, req.Source, req.Type,req.Strategy,req.Target,req.SourceID)
		return EndResponse{V:"end alarm data", Err: err}, nil
	}
}

type Failer interface {
	Failed() error
}


type CreateRequest struct {
	ID string `json:"ID"`
	FlowID uint32 `json:"FlowID"`
	Source string `json:"Source"`
	Type string `json:"Type"`
	Strategy string `json:"Strategy"`
	Target string `json:"Target"`
	SourceID string `json:"SourceID"`
}

type AddRequest struct {
	ID string `json:"ID"`
	FlowID uint32 `json:"FlowID"`
	Source string `json:"Source"`
	Type string `json:"Type"`
	Strategy string `json:"Strategy"`
	Target string `json:"Target"`
	SourceID string `json:"SourceID"`
}

type EndRequest struct {
	ID string `json:"ID"`
	FlowID uint32 `json:"FlowID"`
	Source string `json:"Source"`
	Type string `json:"Type"`
	Strategy string `json:"Strategy"`
	Target string `json:"Target"`
	SourceID string `json:"SourceID"`
}

type CreateResponse struct {
	V   string   `json:"v"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

type AddResponse struct {
	V   string   `json:"v"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

type EndResponse struct {
	V   string   `json:"v"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

// Failed implements Failer.
func (r CreateResponse) Failed() error { return r.Err }

func (r AddResponse) Failed() error { return r.Err }

func (r EndResponse) Failed() error { return r.Err }

