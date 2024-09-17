package interfaces

import (
	"log"

	"github.com/OmarBader7/web-service-jayeek/application"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/auth"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/profile"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	"github.com/OmarBader7/web-service-jayeek/pkg/validator"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
)

type Profile struct {
	AuthService    auth.AuthServiceInterface
	TokenService   auth.TokenInterface
	UserApp        application.UserApplicationInterface
	DriverApp      application.DriverApplicationInterface
	LocationApp    application.LocationApplicationInterface
	ProfileService profile.ProfileServiceInterface
}

// NewProfile creates and returns a new instance of Profile.
func NewProfile(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, userApp application.UserApplicationInterface, driverApp application.DriverApplicationInterface, locationApp application.LocationApplicationInterface, profileService profile.ProfileServiceInterface) *Profile {
	return &Profile{
		AuthService:    authService,
		TokenService:   tokenService,
		UserApp:        userApp,
		DriverApp:      driverApp,
		LocationApp:    locationApp,
		ProfileService: profileService,
	}
}

func (p *Profile) GetProfileDetails(ctx *gin.Context) {
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

	profileDetails, err := p.ProfileService.GetProfileDetails(ctx, user)

	if err != nil {
		log.Fatalf("Error recreating profile details: %v", err)
		response.SendInternalServerError(ctx, ginI18n.MustGetMessage("Failed to retrieve profile details."))
		return
	}

	// Send the response.
	response.SendOK(ctx, profileDetails, "")
}

func (p *Profile) UpdateProfileDetails(ctx *gin.Context) {
	var profileDetails profile.ProfileDetails

	// Bind the JSON body of the request to the Order struct
	if err := ctx.ShouldBindJSON(&profileDetails); err != nil {
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

	// Validate all fields of the profileDetails struct
	validationErrors, _ := validator.ValidateExcept(ctx, &profileDetails, "Phone")
	if validationErrors != nil {
		response.SendUnprocessableEntity(ctx, validationErrors, "")
		return
	}

	// Get the location by its ID
	location, err := p.LocationApp.GetLocationByID(profileDetails.LocationID)
	if err != nil {
		response.SendUnprocessableEntity(ctx, nil, ginI18n.MustGetMessage("Location not found."))
		return
	}

	user.LocationID = location.ID
	user.Name = profileDetails.Name
	user.Email = profileDetails.Email
	user.BirthYear = profileDetails.BirthYear
	user.Location = *location

	updatedProfileDetails, err := p.ProfileService.UpdateProfileDetails(ctx, user)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Send the response.
	response.SendOK(ctx, updatedProfileDetails, "")
}

func (p *Profile) UpdatePhoneNumber(ctx *gin.Context) {
	var profileDetails profile.ProfileDetails

	// Bind the JSON body of the request to the Order struct
	if err := ctx.ShouldBindJSON(&profileDetails); err != nil {
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

	// Validate all fields of the profileDetails struct
	validationErrors, _ := validator.ValidateExcept(ctx, &profileDetails, "Name", "Email", "BirthYear", "LocationID")
	if validationErrors != nil {
		response.SendUnprocessableEntity(ctx, validationErrors, "")
		return
	}

	user.Phone = profileDetails.Phone

	updatedProfileDetails, err := p.ProfileService.UpdateProfileDetails(ctx, user)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Send the response.
	response.SendOK(ctx, updatedProfileDetails, "")
}
