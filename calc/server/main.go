package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/inabinash/grpc/calc/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedCalculatorServer
}

var (
	port = flag.Int("port", 50051, "The server port")
)
// var addrs = "0.0.0.0:5000"
// func (s *Server)SayHello(ctx context.Context , in *pb.GreetRequest)(*pb.GreetResponse , error){
// 	log.Printf("Received : %v\n", in.GetFirstName());
// 	return &pb.GreetResponse{Result: "Hello " + in.GetFirstName()} , nil;
// }

func (s *Server) Calculate(ctx context.Context , in *pb.CalcInput)(*pb.CalcResult, error){
	res:= in.FirstInput+ in.SecondInput;
	fmt.Printf("Result : %v\n", res);
	return &pb.CalcResult{Result: res} ,nil;
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
	pb.RegisterCalculatorServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve")
	}
}
