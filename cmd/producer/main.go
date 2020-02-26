package main

import (
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"stream-grpc/config"
)

func main() {
//	brocker := fmt.Sprintf("%s:%s", config.Config.BrockerUrl, config.Config.BrockerPort)
//	w := kafka.NewWriter(kafka.WriterConfig{
//		Brokers: []string{brocker},
//		Topic:   config.Config.Topic,
//		Balancer: &kafka.LeastBytes{},
//	})
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
//
//	message := kafka.Message{
//		Key:   []byte("transactions"),
//		Value: transactions,
//		Time:  time.Now(),
//	}
//	err := w.WriteMessages(context.Background(), message)
//	if err != nil {
//		fmt.Println("Error to push a new message")
//		fmt.Printf("Error: %v\n", err)
//	}else {
//		fmt.Println("Message was published")
//	}
//
//	err = w.Close()
//	if err != nil {
//		fmt.Println("Error on Close writer")
//		fmt.Printf("Error: %v\n", err)
//	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": config.Config.BrockerUrl})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "transactions"

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          transactions,
	}, nil)
	if err != nil {
		fmt.Printf("Produce error: %v\n", err)
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}