package transport

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"

	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	oldcontext "golang.org/x/net/context"
	"golang.org/x/time/rate"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	"stringsvc1/pb/stringsvc"
	"stringsvc1/svcendpoint"
	"stringsvc1/service"
)

type grpcServer struct {
	uppercase    grpctransport.Handler
	create grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer.
func NewGRPCServer(endpoints svcendpoint.Set, tracer stdopentracing.Tracer, logger log.Logger) pb.AddServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		uppercase: grpctransport.NewServer(
			endpoints.UppercaseEndpoint,
			decodeGRPCUppercaseRequest,
			encodeGRPCUppercaseResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "Uppercase", logger)))...,
		),
		create: grpctransport.NewServer(
			endpoints.CreateEndpoint,
			decodeGRPCCreateRequest,
			encodeGRPCCreateResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "Create", logger)))...,
		),
	}
}

func (s *grpcServer) Uppercase(ctx oldcontext.Context, req *pb.UppercaseRequest) (*pb.UppercaseReply, error) {
	_, rep, err := s.uppercase.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UppercaseReply), nil
}

func (s *grpcServer) Create(ctx oldcontext.Context, req *pb.CreateRequest) (*pb.CreateReply, error) {
	_, rep, err := s.create.ServeGRPC(ctx, req)
	if err != nil{
		return nil, err
	}
	return rep.(*pb.CreateReply), nil
}

// NewGRPCClient returns an AddService backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) service.Service {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	// Each individual endpoint is an http/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	var uppercaseEndpoint endpoint.Endpoint
	{
		uppercaseEndpoint = grpctransport.NewClient(
			conn,
			"pb.Add",
			"Uppercase",
			encodeGRPCUppercaseRequest,
			decodeGRPCUppercaseResponse,
			pb.UppercaseReply{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		uppercaseEndpoint = opentracing.TraceClient(tracer, "Uppercase")(uppercaseEndpoint)
		uppercaseEndpoint = limiter(uppercaseEndpoint)
		uppercaseEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Uppercase",
			Timeout: 30 * time.Second,
		}))(uppercaseEndpoint)
	}

	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = grpctransport.NewClient(
			conn,
			"pb.Add",
			"Create",
			encodeGRPCCreateRequest,
			decodeGRPCCreateResponse,
			pb.CreateReply{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		createEndpoint = opentracing.TraceClient(tracer, "Create")(createEndpoint)
		createEndpoint = limiter(createEndpoint)
		createEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Create",
			Timeout: 30 * time.Second,
		}))(createEndpoint)
	}

	// The Concat endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return svcendpoint.Set{
		UppercaseEndpoint:    uppercaseEndpoint,
		CreateEndpoint:       createEndpoint,
	}
}

// decodeGRPCUppercaseRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC uppercase request to a user-domain uppercase request. Primarily useful in a server.
func decodeGRPCUppercaseRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UppercaseRequest)
	return svcendpoint.UppercaseRequest{A: string(req.S)}, nil
}


// decodeGRPCUppercaseResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC uppercase reply to a user-domain uppercase response. Primarily useful in a client.
func decodeGRPCUppercaseResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UppercaseReply)
	return svcendpoint.UppercaseResponse{V: string(reply.V), Err: str2err(reply.Err)}, nil
}



// encodeGRPCUppercaseResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain uppercase response to a gRPC uppercase reply. Primarily useful in a server.
func encodeGRPCUppercaseResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(svcendpoint.UppercaseResponse)
	return &pb.UppercaseReply{V: string(resp.V), Err: err2str(resp.Err)}, nil
}

// encodeGRPCUppercaseRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain uppercase request to a gRPC uppercase request. Primarily useful in a client.
func encodeGRPCUppercaseRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(svcendpoint.UppercaseRequest)
	return &pb.UppercaseRequest{S: string(req.A)}, nil
}


func decodeGRPCCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateRequest)
	return svcendpoint.CreateRequest{A: string(req.S)}, nil
}


// decodeGRPCUppercaseResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC uppercase reply to a user-domain uppercase response. Primarily useful in a client.
func decodeGRPCCreateResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.CreateReply)
	return svcendpoint.CreateResponse{V: string(reply.V), Err: str2err(reply.Err)}, nil
}



// encodeGRPCUppercaseResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain uppercase response to a gRPC uppercase reply. Primarily useful in a server.
func encodeGRPCCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(svcendpoint.CreateResponse)
	return &pb.CreateReply{V: string(resp.V), Err: err2str(resp.Err)}, nil
}

// encodeGRPCUppercaseRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain uppercase request to a gRPC uppercase request. Primarily useful in a client.
func encodeGRPCCreateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(svcendpoint.CreateRequest)
	return &pb.CreateRequest{ID: string(req.A)}, nil
}


// These annoying helper functions are required to translate Go error types to
// and from strings, which is the type we use in our IDLs to represent errors.
// There is special casing to treat empty strings as nil errors.

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}