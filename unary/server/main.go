package main

import (
	"context"
	"fmt"
	"github.com/junhyuk0801/golang-grpc-practice/unary/server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"net"
)

const HOST string = "localhost"
const PORT int = 8484

type UnaryServer struct {
	pb.UnimplementedUnaryServer
}

func (s *UnaryServer) Call(ctx context.Context, input *pb.SingleData) (*pb.SingleData, error) {
	p, _ := peer.FromContext(ctx)
	log.Printf("got %d from user %s", input.Num, p.Addr.String())
	return &pb.SingleData{Num: input.Num * input.Num}, nil
}

func newServer() *UnaryServer {
	s := &UnaryServer{}
	return s
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", HOST, PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	unaryServer := newServer()
	pb.RegisterUnaryServer(grpcServer, unaryServer)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	} else {
		log.Printf("Server started at %s:%d", HOST, PORT)
	}
}
