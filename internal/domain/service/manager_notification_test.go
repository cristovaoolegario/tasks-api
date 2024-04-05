package service

import (
	"github.com/cristovaoolegario/tasks-api/internal/domain/dto"
	"github.com/cristovaoolegario/tasks-api/internal/infra/kafka"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManagerNotificationService_Notification(t *testing.T) {
	mockPublisher := &kafka.ProducerMock{
		PublishMessageMock: func(topic, message string) error {
			return nil
		},
	}
	service := NewManagerNotificationService("manager-notification", mockPublisher)

	err := service.Notification("test-user", &dto.Task{})

	assert.NoError(t, err)

}
