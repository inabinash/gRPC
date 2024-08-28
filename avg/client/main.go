package main

import (
	"context"
	"log"
	"time"

	pb "github.com/inabinash/grpc/avg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addrs = "127.0.0.1:3000"

func doCalculate(c pb.AvgCalculatorClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.CalculateAvg(ctx)
	// req = []*pb.AvgInput{
	// 	{Input: 1},
	// 	{Input: 2},
	// 	{Input: 3},
	// }

	if err != nil {
		log.Fatalf("Error in creating the stream  %v", err)
	}
	// for i := 0; i < 5; i++ {
	// 	var input int
	// 	fmt.Scan(&input)
	// 	print("input :", input)
	// 	req := &pb.AvgInput{Input: int32(input)}
	// 	// time.Sleep(time.Second);
	// 	if err:= stream.Send(req) ; err != nil {
	// 		log.Fatalf("Error sending in number %v\n", err)
	// 	}
	// }
	numbers := []int32{10, 20, 30, 40}

	for _, number := range numbers {
		if err := stream.Send(&pb.AvgInput{Input: number}); err != nil {
			log.Fatalf("could not send number: %v", err)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("could not receive response: %v", err)
	}

	log.Printf("Average: %v", res.GetResult())
	// msg, err := stream.RecvMsg()
	// log.Printf("Average of the numbers is %f")
}

func main() {
	conn, err := grpc.NewClient(addrs, grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer conn.Close()
	c := pb.NewAvgCalculatorClient(conn)
	doCalculate(c)
}
