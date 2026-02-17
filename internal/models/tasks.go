package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name      string     `json:"name" binding:"required"`
	Completed bool       `json:"completed"`
	Content   string     `json:"content"`
	DueAt     *time.Time `json:"due_at,omitempty"`
	UserId    uint       `gorm:"not null;index"`
	Reminders []Reminder `gorm:"constraint:OnDelete:CASCADE"`
}
