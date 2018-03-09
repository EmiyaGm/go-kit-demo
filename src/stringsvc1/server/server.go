package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"

	pb "stringsvc1/pb/stringsvc"
	"google.golang.org/grpc"
	"strings"
)

const (
	port = ":8080"
)

type stringSvcServer struct {}

func (s *stringSvcServer) Uppercase(ctx context.Context,in *pb.UppercaseRequest) (*pb.UppercaseReply, error) {
	return &pb.UppercaseReply{V:strings.ToUpper(in.S)}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAddServer(s,&stringSvcServer{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

