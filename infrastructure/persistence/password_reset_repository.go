package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// PasswordResetRepository implements the repository.PasswordResetRepository interface
type PasswordResetRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewPasswordResetRepository creates a new instance of the PasswordResetRepository
func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) CreatePasswordReset(passwordReset *entity.PasswordReset) (*entity.PasswordReset, error) {
	if err := r.db.Debug().Model(&passwordReset).Create(&passwordReset).Error; err != nil {
		return nil, err
	}
	if err := r.db.Debug().Model(&passwordReset).Take(&passwordReset).Error; err != nil {
		return nil, err
	}
	return passwordReset, nil
}

func (r *PasswordResetRepository) UpdatePasswordResetByID(id uint64, passwordReset *entity.PasswordReset) (*entity.PasswordReset, error) {
	if err := r.db.Debug().Model(&passwordReset).Updates(passwordReset).Where("id = ?", id).Error; err != nil {
		return nil, err
	}

	return passwordReset, nil
}

func (r *PasswordResetRepository) GetPasswordResetByUserIDAndVerificationCode(userID uint64, verificationCode string) (*entity.PasswordReset, error) {
	var passwordReset entity.PasswordReset
	if err := r.db.Debug().Model(&passwordReset).Where("user_id = ?", userID).Where("verification_code = ?", verificationCode).Order("created_at DESC").Take(&passwordReset).Error; err != nil {
		return nil, err
	}

	return &passwordReset, nil
}
