package client

import (
	"os"
	"flag"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
	"fmt"

	"alarm/alarmservice"
	"alarm/alarmtransport"
	"text/tabwriter"
)

const (
	// StrategyEnd 结束
	StrategyEnd = "end"
	// StrategyAdd 增加点
	StrategyAdd = "add"
	// StrategyCreate 创建
	StrategyCreate = "create"
	port = "0.0.0.0:8081"
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

// Init 初始化报警服务
func Init() error {
	message = make(chan *Message)
	return nil
}

// Alarm 报警消息通道
func Alarm() chan<- *Message {
	return message
}

// Run 运行报警服务
func Run() {
	fs := flag.NewFlagSet("alarmcli", flag.ExitOnError)
	var grpcAddr = fs.String("grpc-addr", port, "gRPC address of alarmsvc")
	fs.Usage = usageFor(fs, os.Args[0]+" [flags] <a> <b>")
	fs.Parse(os.Args[1:])
	var (
		svc alarmservice.Service
		err error
	)
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	svc = alarmtransport.NewGRPCClient(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	for msg := range message {
		switch msg.Strategy {
		case StrategyAdd:
			v,err:= svc.Add(context.Background(), string(msg.ID),uint32(msg.FlowID),string(msg.Source),string(msg.Type),string(msg.Strategy),string(msg.Target),string(msg.SourceID))
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				//os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "%s\n", v)
		case StrategyCreate:
			v,err:= svc.Create(context.Background(), string(msg.ID),uint32(msg.FlowID),string(msg.Source),string(msg.Type),string(msg.Strategy),string(msg.Target),string(msg.SourceID))
			panic(err)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				//os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "%s\n", v)
		case StrategyEnd:
			v,err:= svc.End(context.Background(), string(msg.ID),uint32(msg.FlowID),string(msg.Source),string(msg.Type),string(msg.Strategy),string(msg.Target),string(msg.SourceID))
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				//os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "%s\n", v)
		}
	}
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
