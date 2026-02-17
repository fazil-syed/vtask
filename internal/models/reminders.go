package models

import (
	"time"

	"gorm.io/gorm"
)

type Reminder struct {
	gorm.Model
	TaskID   uint `gorm:"index;not null"`
	RemindAt time.Time
	Type     string
}
