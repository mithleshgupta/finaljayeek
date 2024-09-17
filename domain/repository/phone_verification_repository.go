package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// PhoneVerificationRepository defines the methods that a setting repository should implement
type PhoneVerificationRepository interface {
	CreatePhoneVerification(phoneVerification *entity.PhoneVerification) (*entity.PhoneVerification, error)
	UpdatePhoneVerificationByID(id uint64, phoneVerification *entity.PhoneVerification) (*entity.PhoneVerification, error)
	GetPhoneVerificationByPhoneAndCode(phone string, code string) (*entity.PhoneVerification, error)
}
