package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/inabinash/grpc/calc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addrs = "127.0.0.1:3000"
func main() {
	// create grpc client
	conn , err := grpc.NewClient(addrs ,grpc.WithTransportCredentials(insecure.NewCredentials()));
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// defer conn.Close()

	c := pb.NewCalculatorClient(conn);
	ctx, cancel := context.WithTimeout(context.Background() ,time.Second);
	
	defer cancel();
	res , err:= c.Calculate(ctx,&pb.CalcInput{FirstInput: 2 , SecondInput: 3});
	if err != nil {
		log.Fatalf("Couldn't Do the summation: %v", err);
	}
	fmt.Printf("Calculated Result : %v\n", res.Result);

}