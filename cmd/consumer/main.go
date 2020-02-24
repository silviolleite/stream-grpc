package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/segmentio/kafka-go"
	_ "github.com/segmentio/kafka-go/snappy"
	"stream-grpc/config"
	"stream-grpc/models"
)

func main() {
	brocker := fmt.Sprintf("%s:%s", config.Config.BrockerUrl, config.Config.BrockerPort)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{brocker},
		GroupID:   config.Config.GroupID,
		Topic:     config.Config.Topic,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	fmt.Println("Cosumer is running...")
	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			fmt.Println("Error to fetch message")
			fmt.Printf("Error: %v", err)
			break
		}

		// save on database
		SaveTransactions(m.Value)

		err = r.CommitMessages(context.Background(), m)
		if err != nil {
			fmt.Println("Error to commit message")
			fmt.Printf("Error: %v", err)
			break
		}
	}

	err := r.Close()
	if err != nil {
		fmt.Println("Error on Close reader")
		fmt.Printf("Error: %v\n", err)
	}
}

func SaveTransactions(value []byte)  {

	var extractedJson models.Statement

	if err := json.Unmarshal(value, &extractedJson); err != nil {
		fmt.Println("Error Unmarshal transactions")
		fmt.Printf("Error: %v\n", err)
	}

	db, err := gorm.Open("sqlite3", config.Config.SqlitePath)
	if err != nil {
		fmt.Println("Error to connect with database")
		fmt.Printf("Error: %v\n", err)
	}
	defer db.Close()
	db.AutoMigrate(&models.Transaction{})

	fmt.Printf("Received data %v\n", extractedJson)
	for _, transactions := range extractedJson.JsonTransaction {
		transaction := models.Transaction{
			Identity:          transactions.ID,
			Amount:      transactions.Amount,
			Date:        transactions.Date,
			Description: transactions.Description,
		}
		fmt.Printf("Insert data %v\n", transaction)
		if err := db.FirstOrCreate(&transaction).Error; err != nil {
			fmt.Println("Error to save data")
			fmt.Printf("Error: %v\n", err)
		}else{
			fmt.Println("Data entered successfully")
		}

	}

}
