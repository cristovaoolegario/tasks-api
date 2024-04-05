package service

import (
	"fmt"
	"github.com/cristovaoolegario/tasks-api/internal/domain/dto"
	"github.com/cristovaoolegario/tasks-api/internal/infra/kafka"
)

type ManagerNotificationService interface {
	Notification(userId string, task *dto.Task) error
}

type ManagerNotificationServiceImp struct {
	topic     string
	publisher kafka.ProducerService
}

func NewManagerNotificationService(topic string, publisher kafka.ProducerService) *ManagerNotificationServiceImp {
	return &ManagerNotificationServiceImp{
		topic:     topic,
		publisher: publisher,
	}
}

func (s *ManagerNotificationServiceImp) Notification(userId string, task *dto.Task) error {
	message := fmt.Sprintf("User %s updated task %v", userId, task)
	return s.publisher.PublishMessage(s.topic, message)
}
