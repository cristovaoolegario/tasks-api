package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerServiceImp struct {
	BootstrapServers string //"localhost:9094"
	GroupId          string //"gostats"
	AutoOffsetReset  string //"earliest"
}

func NewConsumerServiceImp(bootstrapServers, groupId, autoOffsetReset string) *ConsumerServiceImp {
	return &ConsumerServiceImp{
		BootstrapServers: bootstrapServers,
		GroupId:          groupId,
		AutoOffsetReset:  autoOffsetReset,
	}
}

type ConsumerService interface {
	Consume(topics []string, servers string, msgChan chan *kafka.Message)
}

func (cs *ConsumerServiceImp) Consume(topics []string, msgChan chan *kafka.Message) {
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		BootStrapServers: cs.BootstrapServers,
		GroupId:          cs.GroupId,
		AutoOffsetReset:  cs.AutoOffsetReset,
	})

	if err != nil {
		panic(err)
	}

	kafkaConsumer.SubscribeTopics(topics, nil)

	for {
		msg, err := kafkaConsumer.ReadMessage(-1)
		if err == nil {
			msgChan <- msg
		}
	}
}

func ProcessEvents(msgChan chan *kafka.Message) {
	for msg := range msgChan {
		fmt.Println("Received message", string(msg.Value), "on topic", *msg.TopicPartition.Topic)
	}
}
