package main

import (
	"context"
	"fmt"
	"github.com/ngaxavi/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)
	x := req.GetX()
	y := req.GetY()
	res := &calculatorpb.SumResponse{
		Sum: x + y,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PNDRequest, stream calculatorpb.SumService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Prime Number Decomposition function was invoked with %v\n", req)
	number := req.GetNumber()
	var k int32
	k = 2
	for number > 1 {
		if number % k == 0 {
			res := &calculatorpb.PNDResponse{
				Prime: k,
			}
			stream.Send(res)
			number /= k
		} else {
			k++
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Println("The server is running n port 50051")
}
