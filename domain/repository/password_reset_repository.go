package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// PasswordResetRepository defines the methods that a setting repository should implement
type PasswordResetRepository interface {
	CreatePasswordReset(passwordReset *entity.PasswordReset) (*entity.PasswordReset, error)
	UpdatePasswordResetByID(id uint64, passwordReset *entity.PasswordReset) (*entity.PasswordReset, error)
	GetPasswordResetByUserIDAndVerificationCode(userID uint64, verificationCode string) (*entity.PasswordReset, error)
}
