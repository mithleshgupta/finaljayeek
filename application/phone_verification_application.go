package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// PhoneVerificationApplication struct is responsible for handling setting-related business logic
type PhoneVerificationApplication struct {
	phoneVerificationRepo repository.PhoneVerificationRepository
}

var _ PhoneVerificationApplicationInterface = &PhoneVerificationApplication{}

// PhoneVerificationApplicationInterface defines the methods that PhoneVerificationApplication should implement
type PhoneVerificationApplicationInterface interface {
	GetPhoneVerificationByPhoneAndCode(phone string, code string) (*entity.PhoneVerification, error)
	CreatePhoneVerification(phoneVerification *entity.PhoneVerification) (*entity.PhoneVerification, error)
	UpdatePhoneVerificationByID(id uint64, phoneVerification *entity.PhoneVerification) (*entity.PhoneVerification, error)
}

func (a *PhoneVerificationApplication) GetPhoneVerificationByPhoneAndCode(phone string, code string) (*entity.PhoneVerification, error) {
	return a.phoneVerificationRepo.GetPhoneVerificationByPhoneAndCode(phone, code)
}

func (a *PhoneVerificationApplication) CreatePhoneVerification(phoneVerification *entity.PhoneVerification) (*entity.PhoneVerification, error) {
	return a.phoneVerificationRepo.CreatePhoneVerification(phoneVerification)
}

func (a *PhoneVerificationApplication) UpdatePhoneVerificationByID(id uint64, phoneVerification *entity.PhoneVerification) (*entity.PhoneVerification, error) {
	return a.phoneVerificationRepo.UpdatePhoneVerificationByID(id, phoneVerification)
}
