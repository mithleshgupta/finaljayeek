package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// PhoneVerificationRepository implements the repository.PhoneVerificationRepository interface
type PhoneVerificationRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewPhoneVerificationRepository creates a new instance of the PhoneVerificationRepository
func NewPhoneVerificationRepository(db *gorm.DB) *PhoneVerificationRepository {
	return &PhoneVerificationRepository{db: db}
}

func (r *PhoneVerificationRepository) CreatePhoneVerification(phoneVerification *entity.PhoneVerification) (*entity.PhoneVerification, error) {
	if err := r.db.Debug().Model(&phoneVerification).Create(&phoneVerification).Error; err != nil {
		return nil, err
	}
	if err := r.db.Debug().Model(&phoneVerification).Take(&phoneVerification).Error; err != nil {
		return nil, err
	}
	return phoneVerification, nil
}

func (r *PhoneVerificationRepository) UpdatePhoneVerificationByID(id uint64, phoneVerification *entity.PhoneVerification) (*entity.PhoneVerification, error) {
	if err := r.db.Debug().Model(&phoneVerification).Updates(phoneVerification).Where("id = ?", id).Error; err != nil {
		return nil, err
	}

	return phoneVerification, nil
}

func (r *PhoneVerificationRepository) GetPhoneVerificationByPhoneAndCode(phone string, code string) (*entity.PhoneVerification, error) {
	var phoneVerification entity.PhoneVerification
	if err := r.db.Debug().Model(&phoneVerification).Where("phone = ?", phone).Where("code = ?", code).Order("created_at DESC").Take(&phoneVerification).Error; err != nil {
		return nil, err
	}

	return &phoneVerification, nil
}
