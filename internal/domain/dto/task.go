package dto

import (
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"time"
)

type Task struct {
	ID            uint       `json:"id" example:"1"`
	Summary       string     `json:"summary" example:"Task summary"`
	PerformedDate *time.Time `json:"performed_date" example:"2024-04-10T12:00:00Z"`
	UserID        uint       `json:"user_id" example:"1"`
}

func (t *Task) ToModel() *model.Task {
	return &model.Task{
		Summary:       t.Summary,
		PerformedDate: t.PerformedDate,
		UserID:        t.UserID,
	}
}
