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
type VoiceNote struct {
	gorm.Model
	Filename      string `json:"filename"`
	ContentType   string `json:"content_type"`
	Size          int64  `json:"size"`
	Transcription string `json:"transcription"`
	Status        string `json:"status"`
	FilePath      string `json:"file_path"`
}

type User struct {
	gorm.Model
	UserName     string `db:"user_name" json:"user_name" gorm:"unique;not null"`
	Email        string `db:"email" json:"email" gorm:"unique;not null"`
	PasswordHash string `db:"password_hash" json:"-" gorm:"not null"`
}
