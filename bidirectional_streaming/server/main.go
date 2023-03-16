package main

import (
	"fmt"
	"github.com/junhyuk0801/golang-grpc-practice/bidirectional_streaming/server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"io"
	"log"
	"net"
)

const HOST string = "0.0.0.0"
const PORT int = 8484

type BidirectionalStreamingServer struct {
	pb.BidirectionalStreamingServer
}

func (s *BidirectionalStreamingServer) Call(stream pb.BidirectionalStreaming_CallServer) error {
	p, _ := peer.FromContext(stream.Context())
	log.Printf("user %s sent stream", p.Addr.String())

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		for i := 0; i < int(in.Num); i++ {
			err := stream.Send(&pb.SingleData{Num: int32(i)})
			if err != nil {
				return err
			}
		}
	}
}

func newServer() *BidirectionalStreamingServer {
	s := &BidirectionalStreamingServer{}
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
	pb.RegisterBidirectionalStreamingServer(grpcServer, clientStreamingServer)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	} else {
		log.Printf("Server started at %s:%d", HOST, PORT)
	}
}
