package main

import (
	"github.com/go-kit/kit/log"
	"net"
	"flag"
	alarmpb "alarm/pb/alarm"
	"google.golang.org/grpc"
	"alarm/alarmservice"
	"alarm/alarmtransport"
	"alarm/alarmendpoint"
	"os"
	"fmt"
	"text/tabwriter"
	"github.com/oklog/oklog/pkg/group"
)

const (
	port = "0.0.0.0:8081"
)

func main() {
	fs := flag.NewFlagSet("alarm", flag.ExitOnError)
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	var service = alarmservice.New()
	var endpoints = alarmendpoint.New(service)
	var grpcServer = alarmtransport.NewGRPCServer(endpoints)
	var grpcAddr = fs.String("grpc-addr", port, "gRPC listen address")
	var g group.Group
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])
	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
		os.Exit(1)
	}
	g.Add(func() error {
		logger.Log("transport", "gRPC", "addr", *grpcAddr)
		// we add the Go Kit gRPC Interceptor to our gRPC service as it is used by
		baseServer := grpc.NewServer()
		alarmpb.RegisterAddServer(baseServer, grpcServer)
		return baseServer.Serve(grpcListener)
	}, func(error) {
		grpcListener.Close()
	})
	logger.Log("exit", g.Run())
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

