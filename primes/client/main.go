package main

import (
	"context"
	// "fmt"
	"io"
	"log"
	"time"

	pb "github.com/inabinash/grpc/primes/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)



func doCalculate(c pb.PrimeCalculatorClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.PrimeInput{Input: 210}
	stream, err := c.CalculatePrimes(ctx, req);
	if err != nil {
		log.Fatalf("Couldn't Do the calculation: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Problem faced in receving msg: %v", err)
		}

		log.Printf("Next Prime Number: %v", msg.Result)
	}
}

var addrs = "127.0.0.1:3000"

func main() {
	// create grpc client
	conn, err := grpc.NewClient(addrs, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	defer conn.Close()

	c := pb.NewPrimeCalculatorClient(conn)
	// doGreet(c);
	doCalculate(c)
}
