package main

import (
	"fmt"
	"github.com/junhyuk0801/golang-grpc-practice/server_streaming/server/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

const HOST string = "localhost"
const PORT int = 8484

type ServerStreamingServer struct {
	pb.ServerStreamingServer
}

func (s *ServerStreamingServer) Call(input *pb.SingleData, stream pb.ServerStreaming_CallServer) error {
	log.Printf("got %d\n", input.Num)
	for i := int32(0); i < input.Num; i++ {
		err := stream.Send(&pb.SingleData{Num: i})
		if err != nil {
			return err
		}
	}
	return nil
}

func newServer() *ServerStreamingServer {
	s := &ServerStreamingServer{}
	return s
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", HOST, PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	serverStreamingServer := newServer()
	pb.RegisterServerStreamingServer(grpcServer, serverStreamingServer)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	} else {
		log.Printf("Server started at %s:%d", HOST, PORT)
	}
}
