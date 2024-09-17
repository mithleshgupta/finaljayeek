package interfaces

import (
	"strconv"

	"github.com/OmarBader7/web-service-jayeek/application"
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/auth"
	"github.com/OmarBader7/web-service-jayeek/pkg/language"
	"github.com/OmarBader7/web-service-jayeek/pkg/pagination"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	"github.com/OmarBader7/web-service-jayeek/pkg/upload"
	"github.com/OmarBader7/web-service-jayeek/pkg/validator"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
)

// Drivers holds the driver-related application interfaces
type Drivers struct {
	AuthService           auth.AuthServiceInterface
	TokenService          auth.TokenInterface
	DriverApp             application.DriverApplicationInterface
	UserApp               application.UserApplicationInterface
	TransportationModeApp application.TransportationModeApplicationInterface
	IdentityDocumentApp   application.IdentityDocumentApplicationInterface
}

// NewDrivers returns a new instance of Drivers
func NewDrivers(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, driverApp application.DriverApplicationInterface, userApp application.UserApplicationInterface, transportationModeApp application.TransportationModeApplicationInterface, identityDocumentApp application.IdentityDocumentApplicationInterface) *Drivers {
	return &Drivers{
		AuthService:           authService,
		TokenService:          tokenService,
		DriverApp:             driverApp,
		UserApp:               userApp,
		TransportationModeApp: transportationModeApp,
		IdentityDocumentApp:   identityDocumentApp,
	}
}

// CreateDriver handles the creation of a new driver
func (d *Drivers) CreateDriver(c *gin.Context) {
	var driver entity.Driver

	// Create a new multipart form
	form, err := c.MultipartForm()
	if err != nil {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Invalid request body"))
		return
	}

	// Extract the fields from the multipart form
	transportationModeIDStr := form.Value["transportationModeID"][0]
	car := form.Value["car"][0]
	iDNumber := form.Value["iDNumber"][0]
	latitudeStr := form.Value["latitude"][0]
	longitudeStr := form.Value["longitude"][0]
	gender := form.Value["gender"][0]

	// Convert the values to the appropriate types
	transportationModeID, err := strconv.ParseUint(transportationModeIDStr, 10, 64)
	if err != nil {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Invalid request body"))
		return
	}

	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Invalid request body"))
		return
	}

	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Invalid request body"))
		return
	}

	// Assign the converted values to the driver struct fields
	driver.TransportationModeID = transportationModeID
	driver.Car = car
	driver.IDNumber = iDNumber
	driver.Latitude = latitude
	driver.Longitude = longitude

	switch gender {
	case "male":
		driver.Gender = entity.MaleGender
	case "female":
		driver.Gender = entity.FemaleGender
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

	// Validate all fields of the driver struct except for UserID, User, and TransportationMode
	validationErrors, _ := validator.ValidateExcept(c, &driver, "UserID", "User", "TransportationMode")
	if validationErrors != nil {
		response.SendUnprocessableEntity(c, validationErrors, "")
		return
	}

	if user.IsDriver {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("You have already registered as a driver."))
		return
	}

	if isDriverExists, _ := d.DriverApp.DriverWithFieldExists("id_number", driver.IDNumber); isDriverExists {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("The ID number you've entered already exists with another driver."))
		return
	}

	// Get the transportation mode by its ID
	transportationMode, err := d.TransportationModeApp.GetTransportationModeByID(driver.TransportationModeID)
	if err != nil {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Transportation mode not found."))
		return
	}

	// Extract the multipart file headers from the request
	drivingLicenseIDImage, err := c.FormFile("drivingLicenseIDImage")
	if err != nil {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Missing drivingLicenseIDImage"))
		return
	}

	vehicleRegistrationImage, err := c.FormFile("vehicleRegistrationImage")
	if err != nil {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Missing vehicleRegistrationImage"))
		return
	}

	vehicleFrontPhotoImage, err := c.FormFile("vehicleFrontPhotoImage")
	if err != nil {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Missing vehicleFrontPhotoImage"))
		return
	}

	livePhotoWithIDImage, err := c.FormFile("livePhotoWithIDImage")
	if err != nil {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Missing livePhotoWithIDImage"))
		return
	}

	// Set the UserID and TransportationModeID for the driver struct
	driver.UserID = user.ID
	driver.TransportationModeID = transportationMode.ID

	// Create the new driver
	createdDriver, err := d.DriverApp.CreateDriver(&driver)
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	drivingLicenseIDFileInfo, err := upload.UploadFile(drivingLicenseIDImage, "uploads")
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}
	drivingLicenseIDImageFilename := drivingLicenseIDFileInfo.Name()

	vehicleRegistrationFileInfo, err := upload.UploadFile(vehicleRegistrationImage, "uploads")
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}
	vehicleRegistrationImageFilename := vehicleRegistrationFileInfo.Name()

	vehicleFrontPhotoFileInfo, err := upload.UploadFile(vehicleFrontPhotoImage, "uploads")
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}
	vehicleFrontPhotoImageFilename := vehicleFrontPhotoFileInfo.Name()

	livePhotoWithIDFileInfo, err := upload.UploadFile(livePhotoWithIDImage, "uploads")
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}
	livePhotoWithIDImageFilename := livePhotoWithIDFileInfo.Name()

	var identityDocument entity.IdentityDocument

	identityDocument.UserID = user.ID
	identityDocument.DrivingLicenseIDImage = drivingLicenseIDImageFilename
	identityDocument.VehicleRegistrationImage = vehicleRegistrationImageFilename
	identityDocument.VehicleFrontPhotoImage = vehicleFrontPhotoImageFilename
	identityDocument.LivePhotoWithIDImage = livePhotoWithIDImageFilename

	// Create the new identity document
	_, err = d.IdentityDocumentApp.CreateIdentityDocument(&identityDocument)
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	response.SendOK(c, createdDriver.PublicData(language.GetLanguage(c)), "")
}

// GetAllDrivers retrieves a paginated list of all drivers.
func (d *Drivers) GetAllDrivers(ctx *gin.Context) {
	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	// Get the driver count from the driver application service.
	count, err := d.DriverApp.CountDrivers()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the drivers from the driver application service.
	drivers, err := d.DriverApp.GetAllDrivers(page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(drivers) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No drivers found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(drivers) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var driverPublicData []interface{}

	for _, driver := range drivers {
		driverPublicData = append(driverPublicData, driver.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = driverPublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the drivers as a response.
	response.SendOK(ctx, data, "")
}

// GetDriverByID retrieves a single driver by ID.
func (d *Drivers) GetDriverByID(ctx *gin.Context) {
	// Parse the driver ID from the URL parameter.
	driverID, err := strconv.ParseUint(ctx.Param("driver_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid driver ID."))
		return
	}

	// Get the driver from the driver application service.
	driver, err := d.DriverApp.GetDriverByID(driverID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Driver not found."))
		return
	}

	// Send the driver as a response.
	response.SendOK(ctx, driver.PublicData(language.GetLanguage(ctx)), "")
}

func (d *Drivers) GetDriversByUserLocationID(ctx *gin.Context) {
	// Parse the driver ID from the URL parameter.
	userLocationID, err := strconv.ParseUint(ctx.Param("location_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid location ID."))
		return
	}

	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	// Get the driver count from the driver application service.
	count, err := d.DriverApp.CountDriversByUserLocationID(userLocationID)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	drivers, err := d.DriverApp.GetDriversByUserLocationID(userLocationID, page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(drivers) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No drivers found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(drivers) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var driverPublicData []interface{}

	for _, driver := range drivers {
		driverPublicData = append(driverPublicData, driver.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = driverPublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the drivers as a response.
	response.SendOK(ctx, data, "")
}

// This code is for testing purposes only
type LatLng struct {
	Latitude  float64
	Longitude float64
}
