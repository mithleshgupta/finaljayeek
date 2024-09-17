package interfaces

import (
	"github.com/OmarBader7/web-service-jayeek/application"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/auth"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	"github.com/gin-gonic/gin"
)

// Settings holds the setting-related application interfaces
type Settings struct {
	AuthService  auth.AuthServiceInterface
	TokenService auth.TokenInterface
	SettingApp   application.SettingApplicationInterface
}

// NewSettings returns a new instance of Settings
func NewSettings(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, settingApp application.SettingApplicationInterface) *Settings {
	return &Settings{
		AuthService:  authService,
		TokenService: tokenService,
		SettingApp:   settingApp,
	}
}

// GetAllSettings retrieves a paginated list of all settings.
func (u *Settings) GetAllSettings(ctx *gin.Context) {
	// Get the settings from the setting application service.
	settings, err := u.SettingApp.GetAllSettings()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	var settingPublicData []interface{}

	for _, setting := range settings {
		settingPublicData = append(settingPublicData, setting.PublicData())
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = settingPublicData

	// Send the settings as a response.
	response.SendOK(ctx, data, "")
}
