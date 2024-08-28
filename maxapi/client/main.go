package main

import (
	"context"
	"log"
	"sync"
	"time"

	pb "github.com/inabinash/grpc/maxapi/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addrs = "127.0.0.1:3000"

func doCalculate(c pb.MaxApiCalculatorClient) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()

	stream, err := c.CalculateMax(ctx)

	if err != nil {
		log.Fatalf("Error in creating the stream  %v", err)
	}

	var wg sync.WaitGroup

	firstDone := make(chan bool)
	secondDone := make(chan bool)

	numbers := []int32{10, 20, 10, 40,30,50,100}
	wg.Add(1)
	go func() {
		for _, number := range numbers {
			if err := stream.Send(&pb.MaxApiInput{Input: number}); err != nil {
				log.Fatalf("could not send number: %v", err)
			}
			time.Sleep(time.Second)
			firstDone <- true
			<-secondDone
		}
		close(firstDone)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			_, ok := <-firstDone
			if !ok {
				break
				// Exit loop when first goroutine completes
			}
			res, err := stream.Recv()
			if err != nil {
				log.Fatalf("could not receive response: %v", err)
			}
			secondDone <- true
			log.Printf("Max till Now: %v", res.GetResult())
		}

		close(secondDone);

	}()
	wg.Wait()
	log.Print("Both goroutines have completed.")
	// msg, err := stream.RecvMsg()
	// log.Printf("Average of the numbers is %f")
}

func main() {
	conn, err := grpc.NewClient(addrs, grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer conn.Close()
	c := pb.NewMaxApiCalculatorClient(conn)
	doCalculate(c)
}
