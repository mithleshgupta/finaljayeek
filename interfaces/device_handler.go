package interfaces

import (
	"github.com/OmarBader7/web-service-jayeek/application"
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/auth"
	"github.com/OmarBader7/web-service-jayeek/pkg/language"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	"github.com/OmarBader7/web-service-jayeek/pkg/validator"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
)

// Devices holds the device-related application interfaces
type Devices struct {
	AuthService  auth.AuthServiceInterface
	TokenService auth.TokenInterface
	DeviceApp    application.DeviceApplicationInterface
	UserApp      application.UserApplicationInterface
}

// NewDevices returns a new instance of Devices
func NewDevices(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, deviceApp application.DeviceApplicationInterface, userApp application.UserApplicationInterface) *Devices {
	return &Devices{
		AuthService:  authService,
		TokenService: tokenService,
		DeviceApp:    deviceApp,
		UserApp:      userApp,
	}
}

// CreateDevice handles the creation of a new device
func (d *Devices) CreateDevice(c *gin.Context) {
	var device entity.Device

	// Bind the JSON body of the request to the Device struct
	if err := c.ShouldBindJSON(&device); err != nil {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Invalid request body"))
		return
	}

	// Extract the token metadata from the request
	metadata, err := d.TokenService.ExtractTokenMetadata(c.Request)
	if err != nil {
		response.SendUnauthorized(c, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Fetch the authenticated user's ID from the auth service
	userID, err := d.AuthService.FetchAuth(metadata.AccessTokenUUID)
	if err != nil {
		response.SendUnauthorized(c, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Get the user from the user application service
	user, err := d.UserApp.GetUserByID(userID)
	if err != nil {
		response.SendUnauthorized(c, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Validate all fields of the device struct except for UserID, User
	validationErrors, _ := validator.ValidateExcept(c, &device, "UserID", "User")
	if validationErrors != nil {
		response.SendUnprocessableEntity(c, validationErrors, "")
		return
	}

	// Set the UserID and TransportationModeID for the driver struct
	device.UserID = user.ID

	// Check if device with the given FCM token exists
	isDeviceExists, err := d.DeviceApp.DeviceWithFieldExists("fcm_token", device.FCMToken)
	if err != nil {
		// Handle the error here; this could be a genuine error in checking for the device's existence.
		response.SendInternalServerError(c, err.Error())
		return
	}

	if isDeviceExists {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Device with the provided FCM token already exists."))
		return
	}

	// Create the new device
	createdDevice, err := d.DeviceApp.CreateDevice(&device)
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	response.SendOK(c, createdDevice.PublicData(language.GetLanguage(c)), "")
}
