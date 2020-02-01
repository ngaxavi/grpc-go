package main

import (
	"fmt"
	"github.com/ngaxavi/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
