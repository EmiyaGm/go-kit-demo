package alarmtransport

import (
	"context"

	oldcontext "golang.org/x/net/context"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	"alarm/pb/alarm"
	"google.golang.org/grpc"
	"github.com/go-kit/kit/endpoint"
	"alarm/alarmendpoint"
	"alarm/alarmservice"
	"errors"
	"fmt"
)

type grpcServer struct {
	create grpctransport.Handler
	add    grpctransport.Handler
	end    grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer.
func NewGRPCServer(endpoints alarmendpoint.Set) pb.AddServer {
	return &grpcServer{
		create: grpctransport.NewServer(
			endpoints.CreateEndpoint,
			decodeGRPCCreateRequest,
			encodeGRPCCreateResponse,
		),
		add: grpctransport.NewServer(
			endpoints.AddEndpoint,
			decodeGRPCAddRequest,
			encodeGRPCAddResponse,
		),
		end: grpctransport.NewServer(
			endpoints.EndEndpoint,
			decodeGRPCEndRequest,
			encodeGRPCEndResponse,
		),
	}
}

func decodeGRPCCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	fmt.Print("decodeGRPCCreateRequest")
	req := grpcReq.(*pb.CreateRequest)
	return alarmendpoint.CreateRequest{ID: string(req.ID),FlowID:uint32(req.FlowID),Source:string(req.Source),Type:string(req.Type),Strategy:string(req.Strategy),Target:string(req.Target),SourceID:string(req.SourceID)}, nil
}

func encodeGRPCCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	fmt.Print("encodeGRPCCreateResponse")
	resp := response.(alarmendpoint.CreateResponse)
	return &pb.CreateReply{V: string(resp.V), Err: err2str(resp.Err)}, nil
}

func decodeGRPCAddRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.AddRequest)
	return alarmendpoint.AddRequest{ID: string(req.ID),FlowID:uint32(req.FlowID),Source:string(req.Source),Type:string(req.Type),Strategy:string(req.Strategy),Target:string(req.Target),SourceID:string(req.SourceID)}, nil
}

func encodeGRPCAddResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(alarmendpoint.AddResponse)
	return &pb.AddReply{V: string(resp.V), Err: err2str(resp.Err)}, nil
}

func decodeGRPCEndRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.EndRequest)
	return alarmendpoint.EndRequest{ID: string(req.ID),FlowID:uint32(req.FlowID),Source:string(req.Source),Type:string(req.Type),Strategy:string(req.Strategy),Target:string(req.Target),SourceID:string(req.SourceID)}, nil
}

func encodeGRPCEndResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(alarmendpoint.EndResponse)
	return &pb.EndReply{V: string(resp.V), Err: err2str(resp.Err)}, nil
}


func (s *grpcServer) Create(ctx oldcontext.Context, req *pb.CreateRequest) (*pb.CreateReply, error) {
	fmt.Print("grpc Create")
	_, rep, err := s.create.ServeGRPC(ctx, req)
	if err != nil{
		return nil, err
	}
	return rep.(*pb.CreateReply), nil
}

func (s *grpcServer) Add(ctx oldcontext.Context, req *pb.AddRequest) (*pb.AddReply, error) {
	_, rep, err := s.add.ServeGRPC(ctx, req)
	if err != nil{
		return nil, err
	}
	return rep.(*pb.AddReply), nil
}

func (s *grpcServer) End(ctx oldcontext.Context, req *pb.EndRequest) (*pb.EndReply, error) {
	_, rep, err := s.end.ServeGRPC(ctx, req)
	if err != nil{
		return nil, err
	}
	return rep.(*pb.EndReply), nil
}


func NewGRPCClient(conn *grpc.ClientConn) alarmservice.Service {
	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = grpctransport.NewClient(
			conn,
			"pb.Add",
			"Create",
			encodeGRPCCreateRequest,
			decodeGRPCCreateResponse,
			pb.CreateReply{},
		).Endpoint()
	}

	// The Concat endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var addEndpoint endpoint.Endpoint
	{
		addEndpoint = grpctransport.NewClient(
			conn,
			"pb.Add",
			"Add",
			encodeGRPCAddRequest,
			decodeGRPCAddResponse,
			pb.AddReply{},
		).Endpoint()
	}

	var endEndpoint endpoint.Endpoint
	{
		addEndpoint = grpctransport.NewClient(
			conn,
			"pb.Add",
			"End",
			encodeGRPCEndRequest,
			decodeGRPCEndResponse,
			pb.AddReply{},
		).Endpoint()
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.

	return alarmendpoint.Set{
		CreateEndpoint: createEndpoint,
		AddEndpoint: addEndpoint,
		EndEndpoint: endEndpoint,
	}
}


func encodeGRPCCreateRequest(_ context.Context, request interface{}) (interface{}, error) {
	fmt.Print("encodeGRPCCreateRequest")
	req := request.(alarmendpoint.CreateRequest)
	return &pb.CreateRequest{ID: string(req.ID),FlowID:uint32(req.FlowID),Source:string(req.Source),Type:string(req.Type),Strategy:string(req.Strategy),Target:string(req.Target),SourceID:string(req.SourceID)}, nil
	}

func decodeGRPCCreateResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	fmt.Print("decodeGRPCCreateResponse")
	reply := grpcReply.(*pb.CreateReply)
	return alarmendpoint.CreateResponse{V: string(reply.V), Err: str2err(reply.Err)}, nil
}

func encodeGRPCAddRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(alarmendpoint.AddRequest)
	return &pb.AddRequest{ID: string(req.ID),FlowID:uint32(req.FlowID),Source:string(req.Source),Type:string(req.Type),Strategy:string(req.Strategy),Target:string(req.Target),SourceID:string(req.SourceID)}, nil
}

func decodeGRPCAddResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.AddReply)
	return alarmendpoint.AddResponse{V: string(reply.V), Err: str2err(reply.Err)}, nil
}

func encodeGRPCEndRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(alarmendpoint.EndRequest)
	return &pb.EndRequest{ID: string(req.ID),FlowID:uint32(req.FlowID),Source:string(req.Source),Type:string(req.Type),Strategy:string(req.Strategy),Target:string(req.Target),SourceID:string(req.SourceID)}, nil
}

func decodeGRPCEndResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.EndReply)
	return alarmendpoint.EndResponse{V: string(reply.V), Err: str2err(reply.Err)}, nil
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