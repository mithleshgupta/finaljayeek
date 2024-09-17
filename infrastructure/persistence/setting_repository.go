package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// SettingRepository implements repository.SettingRepository
// and handles CRUD operations for Setting entities
type SettingRepository struct {
	db *gorm.DB
}

// NewSettingRepository returns a new instance of SettingRepository
func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db}
}

func (r *SettingRepository) GetAllSettings() ([]entity.Setting, error) {
	var settings []entity.Setting
	if err := r.db.Debug().Model(&entity.Setting{}).Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *SettingRepository) GetSettingByKey(key string) (string, error) {
	var setting entity.Setting
	if err := r.db.Debug().Model(&setting).Where("key = ?", key).Take(&setting).Error; err != nil {
		return "", err
	}
	return setting.Value, nil
}
