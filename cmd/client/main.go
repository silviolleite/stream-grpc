package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "stream-grpc/transactions"
)

const (
	address = "localhost:50051"
	account = "121221-1"
	branch  = "1212"
)

type data struct {
	Account string
	Branch  string
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTransactorClient(conn)

	accountInformation := data{
		Account: account,
		Branch:  branch,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.GetTransactions(ctx, &pb.TransactionsRequest{Account: accountInformation.Account, Branch: accountInformation.Branch})
	if err != nil {
		log.Fatalf("Could not get transactions: %v", err)
	}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetTransactions(_) = _, %v", c, err)
		}
		log.Println(data)
	}
}
