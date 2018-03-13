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

	"alarm/pb/alarm"
	"alarm/alarmendpoint"
	"alarm/service"
)

type grpcServer struct {
	create grpctransport.Handler
	add    grpctransport.Handler
	end    grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer.
func NewGRPCServer(endpoints alarmendpoint.Set, tracer stdopentracing.Tracer, logger log.Logger) pb.AddServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		create: grpctransport.NewServer(
			endpoints.CreateEndpoint,
			decodeGRPCCreateRequest,
			encodeGRPCCreateResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "Create", logger)))...,
		),
		add: grpctransport.NewServer(
			endpoints.AddEndpoint,
			decodeGRPCAddRequest,
			encodeGRPCAddResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "Add", logger)))...,
		),
		end: grpctransport.NewServer(
			endpoints.EndEndpoint,
			decodeGRPCEndRequest,
			encodeGRPCEndResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "End", logger)))...,
		),
	}
}


func (s *grpcServer) Create(ctx oldcontext.Context, req *pb.CreateRequest) (*pb.CreateReply, error) {
	_, rep, err := s.create.ServeGRPC(ctx, req)
	if err != nil{
		return nil, err
	}
	return rep.(*pb.CreateReply), nil
}

func (s *grpcServer) Add(ctx oldcontext.Context, req *pb.AddRequest) (*pb.AddReply, error) {
	_, rep, err := s.create.ServeGRPC(ctx, req)
	if err != nil{
		return nil, err
	}
	return rep.(*pb.AddReply), nil
}

func (s *grpcServer) End(ctx oldcontext.Context, req *pb.EndRequest) (*pb.EndReply, error) {
	_, rep, err := s.create.ServeGRPC(ctx, req)
	if err != nil{
		return nil, err
	}
	return rep.(*pb.EndReply), nil
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

	var addEndpoint endpoint.Endpoint
	{
		addEndpoint = grpctransport.NewClient(
			conn,
			"pb.Add",
			"Add",
			encodeGRPCAddRequest,
			decodeGRPCAddResponse,
			pb.AddReply{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		addEndpoint = opentracing.TraceClient(tracer, "Add")(addEndpoint)
		addEndpoint = limiter(addEndpoint)
		addEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Add",
			Timeout: 30 * time.Second,
		}))(addEndpoint)
	}

	var endEndpoint endpoint.Endpoint
	{
		endEndpoint = grpctransport.NewClient(
			conn,
			"pb.Add",
			"End",
			encodeGRPCEndRequest,
			decodeGRPCEndResponse,
			pb.CreateReply{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		endEndpoint = opentracing.TraceClient(tracer, "End")(endEndpoint)
		endEndpoint = limiter(endEndpoint)
		endEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "End",
			Timeout: 30 * time.Second,
		}))(endEndpoint)
	}

	return alarmendpoint.Set{
		CreateEndpoint:       createEndpoint,
		AddEndpoint:          endEndpoint,
		EndEndpoint:          endEndpoint,
	}
}


func decodeGRPCCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateRequest)
	return alarmendpoint.CreateRequest{A: string(req.S)}, nil
}

func decodeGRPCCreateResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.CreateReply)
	return alarmendpoint.CreateResponse{V: string(reply.V), Err: str2err(reply.Err)}, nil
}

func encodeGRPCCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(alarmendpoint.CreateResponse)
	return &pb.CreateReply{V: string(resp.V), Err: err2str(resp.Err)}, nil
}

func encodeGRPCCreateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(alarmendpoint.CreateRequest)
	return &pb.CreateRequest{ID: string(req.A)}, nil
}


func decodeGRPCAddRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.AddRequest)
	return alarmendpoint.AddRequest{A: string(req.S)}, nil
}

func decodeGRPCAddResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.AddReply)
	return alarmendpoint.AddResponse{V: string(reply.V), Err: str2err(reply.Err)}, nil
}

func encodeGRPCAddResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(alarmendpoint.AddResponse)
	return &pb.AddReply{V: string(resp.V), Err: err2str(resp.Err)}, nil
}

func encodeGRPCAddRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(alarmendpoint.AddRequest)
	return &pb.AddRequest{ID: string(req.A)}, nil
}


func decodeGRPCEndRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.EndRequest)
	return alarmendpoint.EndRequest{A: string(req.S)}, nil
}

func decodeGRPCEndResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.EndReply)
	return alarmendpoint.EndResponse{V: string(reply.V), Err: str2err(reply.Err)}, nil
}

func encodeGRPCEndResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(alarmendpoint.EndResponse)
	return &pb.EndReply{V: string(resp.V), Err: err2str(resp.Err)}, nil
}

func encodeGRPCEndRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(alarmendpoint.EndRequest)
	return &pb.EndRequest{ID: string(req.A)}, nil
}

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