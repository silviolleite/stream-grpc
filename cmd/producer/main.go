package main

import (
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"stream-grpc/config"
)

func main() {
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
		fmt.Printf("Message could not be enqueued: %v\n", err)
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}