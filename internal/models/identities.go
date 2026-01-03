package models

import "gorm.io/gorm"

type Identity struct {
	gorm.Model
	// Forign key fields to link with User
	UserId uint `gorm:"not null;index" json:"-"`
	User   User `gorm:"constraint:OnDelete:CASCADE" json:"-"`

	// Issuer Fields
	Issuer  string `gorm:"uniqueIndex:idx_issuer_subject;not null" json:"issuer"`  // "password", "google", "github", etc
	Subject string `gorm:"uniqueIndex:idx_issuer_subject;not null" json:"subject"` // password: internal id, oidc: sub

	// Email from this identity (may be nil)
	Email         *string `gorm:"index" json:"email,omitempty"`
	EmailVerified bool    `gorm:"default:false" json:"email_verified"`

	// Password auth only
	PasswordHash *string `json:"-"`
}
