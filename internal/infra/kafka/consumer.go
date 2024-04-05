package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerServiceImp struct {
	BootstrapServers string
	GroupId          string
	AutoOffsetReset  string
}

func NewConsumerServiceImp(bootstrapServers string, groupId string, autoOffsetReset string) *ConsumerServiceImp {
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

func ProcessEvents(msgChan chan *kafka.Message, LogFn func(*kafka.Message)) {
	for msg := range msgChan {
		LogFn(msg)
	}
}
