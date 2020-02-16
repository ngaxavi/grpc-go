package main

import (
	"context"
	"fmt"
	"github.com/ngaxavi/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Couldn't connect: %v", err)
	}

	defer connection.Close()

	serviceClient := calculatorpb.NewSumServiceClient(connection)

	fmt.Println("Client has been created")

	//calculateSum(serviceClient)
	//primeNumberDecomposition(serviceClient)

	//computeAverage(serviceClient)

	findMaximum(serviceClient)
}

func calculateSum(serviceClient calculatorpb.SumServiceClient)  {
	req := &calculatorpb.SumRequest{
		X: 7,
		Y: 10,
	}

	res, err := serviceClient.Sum(context.Background(), req)
	if err!= nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}

	log.Printf("%v + %v = %v", req.GetX(), req.GetY(), res.Sum)
}

func primeNumberDecomposition(serviceClient calculatorpb.SumServiceClient) {
	req := &calculatorpb.PNDRequest{
		Number: 120,
	}

	res, err := serviceClient.PrimeNumberDecomposition(context.Background(), req)

	if err!= nil {
		log.Fatalf("Error while calling Prime Number Decompostion RPC: %v", err)
	}

	for {
		msg, err := res.Recv()
		if  err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("%v", msg.GetPrime())
	}
}

func computeAverage(serviceClient calculatorpb.SumServiceClient) {
	fmt.Println("Streaming to do a client streaming RPC...")

	requests := []*calculatorpb.NumberRequest{
		&calculatorpb.NumberRequest{
			Number: 1,
		},
		&calculatorpb.NumberRequest{
			Number: 2,
		},
		&calculatorpb.NumberRequest{
			Number: 3,
		},
		&calculatorpb.NumberRequest{
			Number: 4,
		},
	}

	stream, err := serviceClient.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while reading stream: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending number: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("Compute Average: %v\n", res.GetAverage())
}

func findMaximum(serviceClient calculatorpb.SumServiceClient)  {
	fmt.Println("Streaming to do a bidirectional streaming RPC...")

	stream, err := serviceClient.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	waitc := make(chan struct{})

	// we send a bunch of messages to the client ( go routine)
	go func() {
		numbers := []int32{1, 5, 3, 6, 2, 20}
		for _, number := range numbers {
			stream.Send(&calculatorpb.NumberRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// we receive a bunch of messages to the server
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetMax())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc

}