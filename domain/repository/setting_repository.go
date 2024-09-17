package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// SettingRepository defines the methods that a setting repository should implement
type SettingRepository interface {
	GetAllSettings() ([]entity.Setting, error)
	GetSettingByKey(string) (string, error)
}
