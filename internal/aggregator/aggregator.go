package aggregator

import (
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	pb "github.com/net-reply-future-networks/k8s-log-aggregator/api"
	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/utils"
)

type Aggregator struct {
	NatsClient utils.NatsClient
	LogWriter  LogWriter
	GrpcServer GrpcServer
}

func (a *Aggregator) Run() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go a.RunNatsService()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterLogStreamServer(grpcServer, &a.GrpcServer)
	fmt.Println("Listening on port 8000")
	go grpcServer.Serve(lis)
	wg.Wait()
}

func (a *Aggregator) RunNatsService() {
	_, ch, err := a.NatsClient.GetSubscription()
	if err != nil {
		log.Fatal(err)
	}
	if err = a.LogWriter.DbHandlers.Startup(); err != nil {
		log.Fatal(err)
	}
	for {
		msg := <-ch
		fmt.Printf("Sending message to streams, %d\n", len(a.GrpcServer.Streams))
		for _, x := range a.GrpcServer.Streams {
			x.Stream <- msg.Data
		}
		a.LogWriter.WriteLog(msg.Data)
	}
}
