package main

import (
	"context"
	"fmt"
	"github.com/junhyuk0801/golang-grpc-practice/client_streaming/client/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"time"
)

const HOST string = "localhost"
const PORT int = 8484

type ClientStreamingClient struct {
	client pb.ClientStreamingClient
}

func (c *ClientStreamingClient) Call(num int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stream, err := c.client.Call(ctx)
	if err != nil {
		return 0, err
	}

	for i := 1; i <= num; i++ {
		err := stream.Send(&pb.SingleData{Num: int32(i)})
		if err != nil {
			return 0, err
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		return 0, err
	}

	return int(reply.Num), nil
}

func newClient(conn *grpc.ClientConn) *ClientStreamingClient {
	return &ClientStreamingClient{
		client: pb.NewClientStreamingClient(conn),
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
			log.Printf("sent %d, got %d", num, val)
		}

		time.Sleep(2 * time.Second)
	}
}
