package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/net-reply-future-networks/k8s-log-aggregator/api"
)

type Client struct {
	client pb.LogStreamClient
}

func (c *Client) Startup() {
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	c.client = pb.NewLogStreamClient(conn)
}

func (c *Client) Stream() {
	c.Startup()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	stream, err := c.client.StreamLogs(ctx, &pb.StreamsRequest{})
	if err != nil {
		log.Fatalf("client.StreamLog failed: %v", err)
	}
	for {
		logpb, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.StreamLog failed: %v", err)
		}
		fmt.Println(logpb)
	}
}
