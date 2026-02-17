package models

import "gorm.io/gorm"

type PushSubscription struct {
	gorm.Model

	UserID   uint   `gorm:"index"`
	Endpoint string `gorm:"unique;not null"`

	P256dh string `gorm:"not null"`
	Auth   string `gorm:"not null"`
}
