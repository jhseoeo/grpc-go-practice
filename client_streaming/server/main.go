package main

import (
	"fmt"
	"github.com/junhyuk0801/golang-grpc-practice/client_streaming/server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"io"
	"log"
	"net"
)

const HOST string = "0.0.0.0"
const PORT int = 8484

type ClientStreamingServer struct {
	pb.ClientStreamingServer
}

func (s *ClientStreamingServer) Call(stream pb.ClientStreaming_CallServer) error {
	p, _ := peer.FromContext(stream.Context())

	numList := make([]int, 0, 0)
	for {
		val, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		} else {
			numList = append(numList, int(val.Num))
		}
	}
	log.Printf("got %v from user %s", numList, p.Addr.String())

	sum := 0
	for _, num := range numList {
		sum += num
	}

	err := stream.SendAndClose(&pb.SingleData{Num: int32(sum)})
	if err != nil {
		return err
	}

	return nil
}

func newServer() *ClientStreamingServer {
	s := &ClientStreamingServer{}
	return s
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", HOST, PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	clientStreamingServer := newServer()
	pb.RegisterClientStreamingServer(grpcServer, clientStreamingServer)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	} else {
		log.Printf("Server started at %s:%d", HOST, PORT)
	}
}
