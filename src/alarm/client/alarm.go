package main

import (
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "alarm/pb/alarm"
	"time"
)

const (
	// StrategyEnd 结束
	StrategyEnd = "end"
	// StrategyAdd 增加点
	StrategyAdd = "add"
	// StrategyCreate 创建
	StrategyCreate = "create"
	address     = "localhost:8081"
)

// Message 报警消息
type Message struct {
	ID     string
	FlowID uint32
	Source string
	Type   string
	Time   time.Time
	// new add end
	Strategy string

	Target   string
	SourceID string
	Location []float64
	Data     map[string]interface{}
}

var message chan *Message


// Alarm 报警消息通道
func Alarm() chan<- *Message {
	return message
}

// Run 运行报警服务
func Run() {
	for msg := range message {
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewAddClient(conn)
		switch msg.Strategy {
		case StrategyAdd:
			_ = "Add"
			if len(os.Args) > 1 {
				_ = os.Args[1]
			}
			r, err := c.Add(context.Background(), &pb.AddRequest{ID:msg.ID,FlowID:msg.FlowID,Source:msg.Source,Type:msg.Type,Strategy:msg.Strategy,Target:msg.Target,SourceID:msg.SourceID})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			log.Printf("Create: %s", r.V)
		case StrategyCreate:
			_ = "Create"
			if len(os.Args) > 1 {
				_ = os.Args[1]
			}
			r, err := c.Create(context.Background(), &pb.CreateRequest{ID:msg.ID,FlowID:msg.FlowID,Source:msg.Source,Type:msg.Type,Strategy:msg.Strategy,Target:msg.Target,SourceID:msg.SourceID})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			log.Printf("Create: %s", r.V)
		case StrategyEnd:
			_ = "End"
			if len(os.Args) > 1 {
				_ = os.Args[1]
			}
			r, err := c.End(context.Background(), &pb.EndRequest{ID:msg.ID,FlowID:msg.FlowID,Source:msg.Source,Type:msg.Type,Strategy:msg.Strategy,Target:msg.Target,SourceID:msg.SourceID})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			log.Printf("Create: %s", r.V)
		}
	}
}
