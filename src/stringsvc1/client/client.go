package main

import (
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "stringsvc1/pb/stringsvc"
)

const (
	address     = "localhost:8080"
	defaultName = "Uppercase"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAddClient(conn)

	// Contact the server and print out its response.
	_ = defaultName
	if len(os.Args) > 1 {
		_ = os.Args[1]
	}
	r, err := c.Uppercase(context.Background(), &pb.UppercaseRequest{S: "abcdefg"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Uppercase: %s", r.V)

	r2, err := c.Create(context.Background(), &pb.CreateRequest{ID:"gmtest",FlowID:1234,Source:"testSource",Type:"testType"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Create: %s", r2.V)
}
