package dto

import (
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"time"
)

type Task struct {
	ID            uint       `json:"id"`
	Summary       string     `json:"summary"`
	PerformedDate *time.Time `json:"performed_date"`
	UserID        uint       `json:"user_id"`
}

func (t *Task) ToModel() *model.Task {
	return &model.Task{
		Summary:       t.Summary,
		PerformedDate: t.PerformedDate,
		UserID:        t.UserID,
	}
}
