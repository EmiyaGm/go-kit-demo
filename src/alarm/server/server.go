package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"

	pb "alarm/pb/alarm"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"
)

const (
	port = ":8081"
)

const (
	collection = "vehicle_warning"
)

var c *mgo.Collection

type alarmServer struct {}

func (s *alarmServer) Create(ctx context.Context,in *pb.CreateRequest) (*pb.CreateReply , error){
	return &pb.CreateReply{V:"create alarm data"}, nil
}

func (s *alarmServer) Add(ctx context.Context,in *pb.AddRequest) (*pb.AddReply , error){
	return &pb.AddReply{V:"add alarm data"}, nil
}

func (s *alarmServer) End(ctx context.Context,in *pb.EndRequest) (*pb.EndReply , error){
	return &pb.EndReply{V:"end alarm data"}, nil
}

func main() {
	session, err := mgo.Dial("")
	db := session.DB("parse_vehicle")
	c = db.C(collection)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAddServer(s,&alarmServer{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

