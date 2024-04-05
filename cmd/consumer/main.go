package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	kafkaService "github.com/cristovaoolegario/tasks-api/internal/infra/kafka"
)

func main() {
	msgChan := make(chan *kafka.Message)
	var topics = []string{"managerNotification"}
	notificationConsumer := kafkaService.NewConsumerServiceImp("localhost:9092", "gostats", "earliest")

	go notificationConsumer.Consume(topics, msgChan)
	kafkaService.ProcessEvents(msgChan)
}
