package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
	"os"
)

var p *kafka.Producer

func init() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	kafkaHost := os.Getenv("kafka_host")
	var err error
	p, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaHost})

	if err != nil {
		panic(err)
	}
}

func Publish(api []byte) {
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
	kafkaTopic := os.Getenv("kafka_topic")
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          api,
	}, nil)
}
