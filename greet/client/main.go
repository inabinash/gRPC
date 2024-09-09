package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/inabinash/grpc/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	// "google.golang.org/grpc/credentials/insecure"
)

func doGreet(c pb.GreetClient) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.SayHello(ctx, &pb.GreetRequest{FirstName: "Abinash"})
	if err != nil {
		log.Fatalf("Couldn't say hello: %v", err)
	}
	fmt.Printf("Greeting Message : %v\n", res.Result)
}

func doGreetManyTimes(c pb.GreetClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GreetRequest{FirstName: "Abinash"}
	stream, err := c.GreetManyTimes(ctx, req)
	if err != nil {
		log.Fatalf("Couldn't say hello: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Problem faced in receving msg: %v", err)
		}

		log.Printf("Got message: %v", msg.Result)
	}
}

var addrs = "localhost:3000"

func main() {

	creds, err := credentials.NewClientTLSFromFile("server.crt", "")
	if err != nil {
		log.Fatalf("Failed to load server certificate: %v", err)
	}
	// create grpc client
	conn, err := grpc.NewClient(addrs, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	defer conn.Close()

	c := pb.NewGreetClient(conn)
	doGreet(c);
	// doGreetManyTimes(c)
}
