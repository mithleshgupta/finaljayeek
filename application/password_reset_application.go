package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// PasswordResetApplication struct is responsible for handling setting-related business logic
type PasswordResetApplication struct {
	passwordResetRepo repository.PasswordResetRepository
}

var _ PasswordResetApplicationInterface = &PasswordResetApplication{}

// PasswordResetApplicationInterface defines the methods that PasswordResetApplication should implement
type PasswordResetApplicationInterface interface {
	GetPasswordResetByUserIDAndVerificationCode(userID uint64, verificationCode string) (*entity.PasswordReset, error)
	CreatePasswordReset(passwordReset *entity.PasswordReset) (*entity.PasswordReset, error)
	UpdatePasswordResetByID(id uint64, passwordReset *entity.PasswordReset) (*entity.PasswordReset, error)
}

func (a *PasswordResetApplication) GetPasswordResetByUserIDAndVerificationCode(userID uint64, verificationCode string) (*entity.PasswordReset, error) {
	return a.passwordResetRepo.GetPasswordResetByUserIDAndVerificationCode(userID, verificationCode)
}

func (a *PasswordResetApplication) CreatePasswordReset(passwordReset *entity.PasswordReset) (*entity.PasswordReset, error) {
	return a.passwordResetRepo.CreatePasswordReset(passwordReset)
}

func (a *PasswordResetApplication) UpdatePasswordResetByID(id uint64, passwordReset *entity.PasswordReset) (*entity.PasswordReset, error) {
	return a.passwordResetRepo.UpdatePasswordResetByID(id, passwordReset)
}
