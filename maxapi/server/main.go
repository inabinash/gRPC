package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/inabinash/grpc/maxapi/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The Server Port")
)

type Server struct {
	pb.UnimplementedMaxApiCalculatorServer
}

func (*Server) CalculateMax(stream pb.MaxApiCalculator_CalculateMaxServer) error {
	maxi := 0
	for {
		msg, err := stream.Recv()
		fmt.Printf("Received message: %v\n", msg.GetInput())
		if err == io.EOF {
			return stream.Send(&pb.MaxApiOutput{Result: int32(maxi)})
		}

		if err != nil {
			log.Fatalf("Problem faced in receving message %v\n", err)
			break
		}
		maxi = max(maxi, int(msg.GetInput()))
		stream.Send(&pb.MaxApiOutput{Result: int32(maxi)})
	}

	return nil
}

//  func (UnimplementedAvgCalculatorServer) CalculateAvg(grpc.ClientStreamingServer[PrimeInput, PrimeOutput]) error {

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		log.Fatal("Failed to open a network with the provided address .", err)
	}
	fmt.Printf("listning object ... %v\n", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterMaxApiCalculatorServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve")
	}
}
