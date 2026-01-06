package services

import (
	"context"
	"errors"

	"github.com/syed.fazil/vtask/internal/models"
	"github.com/syed.fazil/vtask/internal/schemas"
	"gorm.io/gorm"
)

type IdentityIntent int

const (
	IntentLogin IdentityIntent = iota
	IntentRegister
)

var ErrIdentityAlreadyExists = errors.New("identity already exists")

func FindOrCreateUserWithIdentity(ctx context.Context, db *gorm.DB, input schemas.IdentityInput, intent IdentityIntent) (*models.User, error) {
	dbWithCtx := db.WithContext(ctx)
	var resultUser *models.User
	err := dbWithCtx.Transaction(func(tx *gorm.DB) error {
		var identity models.Identity
		err := tx.Model(&models.Identity{}).Where("issuer = ? AND subject = ?", input.Issuer, input.Subject).Preload("User").First(&identity).Error
		if err == nil {
			if intent == IntentRegister {
				return ErrIdentityAlreadyExists
			}
			resultUser = &identity.User
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		user := models.User{
			PrimaryEmail: input.Email,
			UserName:     input.UserName,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		newIdentity := models.Identity{
			Issuer:        input.Issuer,
			Subject:       input.Subject,
			Email:         input.Email,
			EmailVerified: input.EmailVerified,
			UserId:        user.ID,
			PasswordHash:  input.PasswordHash,
		}

		if err := tx.Create(&newIdentity).Error; err != nil {
			return err
		}
		resultUser = &user
		return nil
	})
	return resultUser, err
}
