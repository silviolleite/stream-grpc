//go:generate protoc -I transactions --go_out=plugins=grpc:transactions transactions/transactions.proto

package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"google.golang.org/grpc"
	"log"
	"net"
	"stream-grpc/config"
	"stream-grpc/models"
	pb "stream-grpc/transactions"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedTransactorServer

	transactions map[string][]*pb.Transactions
}

func (s *server) GetTransactions(account *pb.TransactionsRequest, in pb.Transactor_GetTransactionsServer) error {
	log.Printf("Received: Account: %v, Branch: %v", account.GetAccount(), account.GetBranch())
	transactions := GetTransactionsOnDB()
	if err := in.Send(&transactions); err != nil {
		return err
	}
	return nil
}

func newServer() *server {
	s := &server{transactions: make(map[string][]*pb.Transactions)}
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


func GetTransactionsOnDB() pb.TransactionsReply {
	db, err := gorm.Open("sqlite3", config.Config.SqlitePath)
	if err != nil {
		fmt.Println("Error to connect with database")
		fmt.Printf("Error: %v\n", err)
	}
	defer db.Close()
	transactions := make([]models.Transaction, 0)
	if err := db.Find(&transactions).Error; err != nil {
		fmt.Println(err)
	}
	var transactionsReply pb.TransactionsReply
	for _, transaction := range transactions {
		trans := pb.Transactions{
			Id:                   transaction.Identity,
			Date:                 transaction.Date,
			Description:          transaction.Description,
			Amount:               transaction.Amount,
		}
		transactionsReply.Transactions = append(transactionsReply.Transactions, &trans)
	}
	return transactionsReply
}

