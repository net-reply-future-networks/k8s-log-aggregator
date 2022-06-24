package aggregator

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/net-reply-future-networks/k8s-log-aggregator/api"
)

type GrpcServer struct {
	pb.UnimplementedLogStreamServer
	Streams Streams
	i       int
}

type StreamKey struct {
	Key    int
	Stream chan []byte
}

type Streams []StreamKey

func (s *Streams) DropStream(key int) {
	stream := *s
	for i, x := range stream {
		if x.Key == key {
			*s = append(stream[:i], stream[i+1:]...)
		}
	}
}

func (g *GrpcServer) StreamLogs(_ *pb.StreamsRequest, stream pb.LogStream_StreamLogsServer) error {
	fmt.Println("New stream request received")
	g.i++
	key := g.i
	in := make(chan []byte)
	g.Streams = append(g.Streams, StreamKey{
		Key:    key,
		Stream: in,
	})
	defer g.Streams.DropStream(key)
	for {
		fmt.Println("Waiting for log")
		log := <-in
		logpb := pb.Log{}
		err := protojson.Unmarshal(log, &logpb)
		if err != nil {
			return err
		}
		fmt.Println("sending log")
		stream.Send(&logpb)
	}
}

func NewGrpcServer() *GrpcServer {
	srv := &GrpcServer{}
	return srv
}
