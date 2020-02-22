//go:generate protoc -I transactions --go_out=plugins=grpc:transactions transactions/transactions.proto

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "stream/transactions"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedTransactorServer
	savedTransactions []*pb.TransactionsReply

	transactions map[string][]*pb.Transactions
}

func (s *server) GetTransactions(account *pb.TransactionsRequest, in pb.Transactor_GetTransactionsServer) error {
	log.Printf("Received: Account: %v, Branch: %v", account.GetAccount(), account.GetBranch())
	for _, transaction := range s.savedTransactions {
		if err := in.Send(transaction); err != nil {
			return err
		}
	}
	return nil
}

func newServer() *server {
	s := &server{transactions: make(map[string][]*pb.Transactions)}
	s.loadFeatures()
	return s
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Printf("Server is running on port %s\n", port)
	s := grpc.NewServer()
	pb.RegisterTransactorServer(s, newServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *server) loadFeatures() {
	var data []byte

	data = exampleData

	if err := json.Unmarshal(data, &s.savedTransactions); err != nil {
		log.Fatalf("Failed to load default transactions: %v", err)
	}
}

var exampleData = []byte(`[{
		"id": "1"
		"transactions": {
			"date": "02/02/2020",
			"description": "TEST DESCRIPTION",
			"amount": "1200"
		}
	},
	{
		"id": "2"
		"transactions": {
			"date": "03/03/2020",
			"description": "TEST DESCRIPTION TWO",
			"amount": "1300"
		}
	}
]
`)
