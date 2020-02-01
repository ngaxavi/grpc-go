package main

import (
	"context"
	"fmt"
	"github.com/ngaxavi/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
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

	calculateSum(serviceClient)
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