package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/cristovaoolegario/tasks-api/internal/domain/dto"
	kafkaService "github.com/cristovaoolegario/tasks-api/internal/infra/kafka"
	"os"
)

const (
	BrokerHost   = "BROKER_HOST"
	ManagerTopic = "MANAGER_NOTIFICATION_TOPIC"
)

func main() {
	msgChan := make(chan *kafka.Message)

	var topics = []string{os.Getenv(ManagerTopic)}
	notificationConsumer := kafkaService.NewConsumerServiceImp(
		os.Getenv(BrokerHost), "gostats", "earliest")

	go notificationConsumer.Consume(topics, msgChan)
	kafkaService.ProcessEvents(msgChan, func(message *kafka.Message) {
		var updatedTaskDto *dto.Task
		if err := json.Unmarshal(message.Value, &updatedTaskDto); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("The tech %d performed the task %d on date %s \n", updatedTaskDto.UserID, updatedTaskDto.ID, updatedTaskDto.PerformedDate)
	})
}
