package main

import (
	"fmt"
	"github.com/junhyuk0801/golang-grpc-practice/server_streaming/server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"net"
)

const HOST string = "0.0.0.0"
const PORT int = 8484

type ServerStreamingServer struct {
	pb.ServerStreamingServer
}

func (s *ServerStreamingServer) Call(input *pb.SingleData, stream pb.ServerStreaming_CallServer) error {
	p, _ := peer.FromContext(stream.Context())
	log.Printf("got %d from user %s\n", input.Num, p.Addr.String())
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
