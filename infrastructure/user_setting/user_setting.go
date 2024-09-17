package user_setting

import (
	"context"
	"encoding/json"

	"github.com/OmarBader7/web-service-jayeek/application"
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

// UserSettingServiceInterface defines the methods that a chat service should implement.
type UserSettingServiceInterface interface {
	GetUserSetting(ctx *gin.Context, user *entity.User) (*UserSetting, error)
	UpdateUserSetting(ctx *gin.Context, user *entity.User) (*UserSetting, error)
}

// UserSettingService represents the profile service implementation
type UserSettingService struct {
	UserApp application.UserApplicationInterface
}

// UserSetting represents the details of a profile.
type UserSetting struct {
	IsAvailable    *bool `json:"is_available" validate:"omitempty"`
	IsDarkMode     *bool `json:"is_dark_mode" validate:"omitempty"`
	Is24HourFormat *bool `json:"is_24_hour_format" validate:"omitempty"`
}

var _ UserSettingServiceInterface = &UserSettingService{}

// NewUserSettingService creates and returns a new instance of UserSettingService
func NewUserSettingService(userApp application.UserApplicationInterface) *UserSettingService {
	return &UserSettingService{
		UserApp: userApp,
	}
}

func (s *UserSettingService) GetUserSetting(ctx *gin.Context, user *entity.User) (*UserSetting, error) {
	var settings UserSetting
	if user.Settings != nil {
		err := json.Unmarshal([]byte(user.Settings), &settings)
		if err != nil {
			return nil, err
		}
	}

	return &settings, nil
}

func (s *UserSettingService) UpdateUserSetting(ctx *gin.Context, user *entity.User) (*UserSetting, error) {
	updatedUser, err := s.UserApp.UpdateUserByID(user.ID, user)
	if err != nil {
		return nil, err
	}

	var settings UserSetting
	err = json.Unmarshal([]byte(updatedUser.Settings), &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}
