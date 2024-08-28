package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/inabinash/grpc/avg/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The Server Port")
)

type Server struct {
	pb.UnimplementedAvgCalculatorServer
}

func (*Server) CalculateAvg(stream pb.AvgCalculator_CalculateAvgServer) error {
	sum := 0
	cnt := 0
	for {
		
		msg, err := stream.Recv()
		fmt.Printf("Received message: %v\n", msg.GetInput());
		if err == io.EOF {
			if cnt == 0 {
				return stream.SendAndClose(&pb.AvgOutput{Result: 0})
			}
			avg := float32(sum) / float32(cnt)
			return stream.SendAndClose(&pb.AvgOutput{Result: avg})
		}

		if err != nil {
			log.Fatalf("Problem faced in receving message %v\n", err)
			break
		}
		cnt++
		sum += int(msg.GetInput())

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
	pb.RegisterAvgCalculatorServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve")
	}
}
