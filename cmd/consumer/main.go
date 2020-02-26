package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"stream-grpc/config"
	"stream-grpc/models"
)

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.Config.BrockerUrl,
		"group.id":          "machines",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	err = c.SubscribeTopics([]string{"transactions", "^aRegex.*[Tt]opic"}, nil)
	if err != nil {
		fmt.Printf("Error subscribing to topic: %v\n", err)
	}
	defer c.Close()

	fmt.Println("Consumer is listening...")

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			// save on database
			SaveTransactions(msg.Value)
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
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
		fmt.Printf("Inserting data %v\n", transaction)
		if err := db.FirstOrCreate(&transaction).Error; err != nil {
			fmt.Println("Error to save data")
			fmt.Printf("Error: %v\n", err)
		}else{
			fmt.Println("Data entered successfully")
		}

	}

}
