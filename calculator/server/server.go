package main

import (
	"context"
	"fmt"
	"github.com/ngaxavi/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
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
	k := int64(2)
	for number > 1 {
		if number%k == 0 {
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

func (*server) ComputeAverage(stream calculatorpb.SumService_ComputeAverageServer) error {
	fmt.Println("Compute average function was invoked with a streaming request")
	counter := 0
	sum := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Average: float64(sum) / float64(counter),
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		sum += req.GetNumber()
		counter++
	}
}

func (*server) FindMaximum(stream calculatorpb.SumService_FindMaximumServer) error {
	fmt.Println("FindMaximum function was invoked with a streaming request")

	max := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		number := req.GetNumber()

		if max < number {
			max = number
			sendErr := stream.Send(&calculatorpb.MaximumResponse{
				Max: max,
			})

			if sendErr != nil {
				log.Fatalf("Error while sending data to cleint: %v", sendErr)
				return sendErr
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.NumberRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument, fmt.Sprintf("Received a negative number: %v", number),
			)
	}

	return &calculatorpb.SquareRootResponse{
		Root: math.Sqrt(float64(number)),
	}, nil

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
