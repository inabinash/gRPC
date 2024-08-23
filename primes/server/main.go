package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/inabinash/grpc/primes/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedPrimeCalculatorServer
}

var (
	port = flag.Int("port", 50051, "The server port")
)

// var addrs = "0.0.0.0:5000"

func (s *Server) CalculatePrimes(in *pb.PrimeInput, stream pb.PrimeCalculator_CalculatePrimesServer) error {
	// for i := 0; i < 10; i++ {
	// 	res := fmt.Sprintf("Hello  %s in %d", in.GetFirstName(), i)
	// 	stream.Send(&pb.GreetResponse{Result: res})
	// }

	k := int32(2);
	N := in.Input

	for N > 1 {
		if N%k == 0 {
			log.Printf("Prime Number %d \n", k);
			stream.Send(&pb.PrimeOutput{Result: k});
			N = N / k
		} else {
			k = k + 1
		}
	}

	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen")
	}
	fmt.Printf("listning object ... %v\n", lis.Addr())
	s := grpc.NewServer()
	// fmt.Printf("grpc Server %v\n", s)
	pb.RegisterPrimeCalculatorServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve")
	}
}
