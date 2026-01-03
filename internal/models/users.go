package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName     *string    `json:"user_name,omitempty" gorm:"unique;"`
	PrimaryEmail *string    `json:"email,omitempty" gorm:"index"`
	Identities   []Identity `json:"-"`
}
