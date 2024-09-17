package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// SettingApplication struct is responsible for handling setting-related business logic
type SettingApplication struct {
	settingRepo repository.SettingRepository
}

var _ SettingApplicationInterface = &SettingApplication{}

// SettingApplicationInterface defines the methods that SettingApplication should implement
type SettingApplicationInterface interface {
	GetAllSettings() ([]entity.Setting, error)
	GetSettingByKey(key string) (string, error)
}

func (a *SettingApplication) GetAllSettings() ([]entity.Setting, error) {
	return a.settingRepo.GetAllSettings()
}

func (a *SettingApplication) GetSettingByKey(key string) (string, error) {
	return a.settingRepo.GetSettingByKey(key)
}
