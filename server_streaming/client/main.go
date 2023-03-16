package main

import (
	"context"
	"fmt"
	"github.com/junhyuk0801/golang-grpc-practice/server_streaming/client/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"math/rand"
	"time"
)

const HOST string = "server"
const PORT int = 8484

type ServerStreamingClient struct {
	client pb.ServerStreamingClient
}

func (c *ServerStreamingClient) Call(num int) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	input := &pb.SingleData{Num: int32(num)}
	stream, err := c.client.Call(ctx, input)
	if err != nil {
		return nil, nil
	}

	res := make([]int, 0, 0)
	for {
		num, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, nil
		} else {
			res = append(res, int(num.Num))
		}
	}

	return res, nil
}

func newClient(conn *grpc.ClientConn) *ServerStreamingClient {
	return &ServerStreamingClient{
		client: pb.NewServerStreamingClient(conn),
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
