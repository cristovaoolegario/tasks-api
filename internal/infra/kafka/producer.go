package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProducerServiceImp struct {
	Topic            string
	BootstrapServers string //"localhost:9092"
}

type ProducerService interface {
	PublishMessage(topic string, message []byte) error
}

func NewProducerServiceImp(bootstrapServers string) *ProducerServiceImp {
	return &ProducerServiceImp{
		BootstrapServers: bootstrapServers,
	}
}

func (s *ProducerServiceImp) PublishMessage(topic string, message []byte) error {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		BootStrapServers: s.BootstrapServers,
	})
	if err != nil {
		return fmt.Errorf("failed to create producer: %v", err)
	}

	defer producer.Close()

	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to produce message: %v", err)
	}

	// Wait for message deliveries before shutting down
	producer.Flush(15 * 1000)

	return nil
}
