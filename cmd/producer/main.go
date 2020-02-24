package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"stream-grpc/config"
	"time"
)

func main() {
	brocker := fmt.Sprintf("%s:%s", config.Config.BrockerUrl, config.Config.BrockerPort)
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brocker},
		Topic:   config.Config.Topic,
		Balancer: &kafka.LeastBytes{},
	})
	transactions := []byte(`{
		"transactions": [
			{
				"id": "1",
				"date": "02/02/2020",
				"description": "TEST DESCRIPTION",
				"amount": "1200"
			},
			{
				"id": "2",
				"date": "03/03/2020",
				"description": "TEST DESCRIPTION TWO",
				"amount": "1300"
			},
{
				"id": "3",
				"date": "03/04/2020",
				"description": "TEST DESCRIPTION THREE",
				"amount": "1400"
			}
	]}`)

	message := kafka.Message{
		Key:   []byte("transactions"),
		Value: transactions,
		Time:  time.Now(),
	}
	err := w.WriteMessages(context.Background(), message)
	if err != nil {
		fmt.Println("Error to push a new message")
		fmt.Printf("Error: %v\n", err)
	}else {
		fmt.Println("Message was published")
	}

	err = w.Close()
	if err != nil {
		fmt.Println("Error on Close writer")
		fmt.Printf("Error: %v\n", err)
	}

}