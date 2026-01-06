package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name       string     `json:"name" binding:"required"`
	Completed  bool       `json:"completed"`
	DueAt      *time.Time `json:"due_at,omitempty"`
	ReminderAt *time.Time `json:"reminder_at,omitempty"`
}
