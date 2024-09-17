package interfaces

import (
	"log"

	"github.com/OmarBader7/web-service-jayeek/application"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/auth"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/user_setting"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	"github.com/OmarBader7/web-service-jayeek/pkg/validator"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
)

type UserSetting struct {
	AuthService        auth.AuthServiceInterface
	TokenService       auth.TokenInterface
	UserApp            application.UserApplicationInterface
	UserSettingService user_setting.UserSettingServiceInterface
}

// NewUserSetting creates and returns a new instance of UserSetting.
func NewUserSetting(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, userApp application.UserApplicationInterface, userSettingService user_setting.UserSettingServiceInterface) *UserSetting {
	return &UserSetting{
		AuthService:        authService,
		TokenService:       tokenService,
		UserApp:            userApp,
		UserSettingService: userSettingService,
	}
}

func (p *UserSetting) GetUserSettings(ctx *gin.Context) {
	// Extract the token metadata from the request
	metadata, err := p.TokenService.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Fetch the authenticated user's ID from the auth service
	userID, err := p.AuthService.FetchAuth(metadata.AccessTokenUUID)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Get the user from the user application service
	user, err := p.UserApp.GetUserByID(userID)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	userSetting, err := p.UserSettingService.GetUserSetting(ctx, user)

	if err != nil {
		log.Fatalf("Error recreating user settings: %v", err)
		response.SendInternalServerError(ctx, ginI18n.MustGetMessage("Failed to retrieve profile details."))
		return
	}

	// Send the response.
	response.SendOK(ctx, userSetting, "")
}

func (p *UserSetting) UpdateUserSettings(ctx *gin.Context) {
	var userSetting user_setting.UserSetting

	// Bind the JSON body of the request to the UserSetting struct
	if err := ctx.ShouldBindJSON(&userSetting); err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid request body"))
		return
	}

	// Extract the token metadata from the request
	metadata, err := p.TokenService.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Fetch the authenticated user's ID from the auth service
	userID, err := p.AuthService.FetchAuth(metadata.AccessTokenUUID)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Get the user from the user application service
	user, err := p.UserApp.GetUserByID(userID)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	validationErrors, _ := validator.ValidateExcept(ctx, &userSetting)
	if validationErrors != nil {
		response.SendUnprocessableEntity(ctx, validationErrors, "")
		return
	}

	// Update the specific settings provided in the request
	if userSetting.IsAvailable != nil {
		user.AddSetting("is_available", *&userSetting.IsAvailable)
	}
	if userSetting.IsDarkMode != nil {
		user.AddSetting("is_dark_mode", *&userSetting.IsDarkMode)
	}
	if userSetting.Is24HourFormat != nil {
		user.AddSetting("is_24_hour_format", *&userSetting.Is24HourFormat)
	}

	updatedUserSetting, err := p.UserSettingService.UpdateUserSetting(ctx, user)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Send the response.
	response.SendOK(ctx, updatedUserSetting, "")
}
