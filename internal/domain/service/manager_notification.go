package service

import (
	"encoding/json"
	"github.com/cristovaoolegario/tasks-api/internal/domain/dto"
	"github.com/cristovaoolegario/tasks-api/internal/infra/kafka"
)

type ManagerNotificationService interface {
	Notification(userId uint, task *dto.Task) error
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

func (s *ManagerNotificationServiceImp) Notification(userId uint, task *dto.Task) error {
	task.UserID = userId
	message, _ := json.Marshal(task)
	return s.publisher.PublishMessage(s.topic, message)
}
