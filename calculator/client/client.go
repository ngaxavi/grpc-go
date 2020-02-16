package main

import (
	"context"
	"fmt"
	"github.com/ngaxavi/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
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
	primeNumberDecomposition(serviceClient)
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