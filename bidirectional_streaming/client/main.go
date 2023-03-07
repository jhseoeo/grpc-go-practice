package main

import (
	"context"
	"fmt"
	"github.com/junhyuk0801/golang-grpc-practice/bidirectional_streaming/client/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"math/rand"
	"time"
)

const HOST string = "localhost"
const PORT int = 8484

type BidirectionalStreamingClient struct {
	client pb.BidirectionalStreamingClient
}

func (c *BidirectionalStreamingClient) Call(cnt int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stream, err := c.client.Call(ctx)
	if err != nil {
		return 0, err
	}

	totalSum := 0
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			} else if err != nil {
				log.Fatalf("failed to recv: %v", err)
				close(waitc)
				return
			}
			totalSum += int(in.Num)
		}
	}()

	for i := 0; i < cnt; i++ {
		err = stream.Send(&pb.SingleData{Num: int32(randInt(1, 10))})
		if err != nil {
			return 0, err
		}
	}

	err = stream.CloseSend()
	if err != nil {
		return 0, err
	}
	<-waitc

	return totalSum, nil
}

func newClient(conn *grpc.ClientConn) *BidirectionalStreamingClient {
	return &BidirectionalStreamingClient{
		client: pb.NewBidirectionalStreamingClient(conn),
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", HOST, PORT), opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := newClient(conn)
	rand.Seed(time.Now().UnixNano())

	for {
		num := randInt(1, 10)
		val, err := client.Call(num)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		} else {
			log.Printf("sent %d times, got %d", num, val)
		}

		time.Sleep(2 * time.Second)
	}
}
