package db

import (
	"github.com/syed.fazil/vtask/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Identity{}, &models.User{}, &models.Task{}); err != nil {
		return err
	}
	return db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS identities_password_email_unique
		ON identities (email)
		WHERE issuer = 'password';
	`).Error
}
