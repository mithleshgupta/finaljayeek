package interfaces

import (
	"strconv"

	"github.com/OmarBader7/web-service-jayeek/application"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/auth"
	"github.com/OmarBader7/web-service-jayeek/pkg/language"
	"github.com/OmarBader7/web-service-jayeek/pkg/pagination"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
)

// Locations holds the location-related application interfaces
type Locations struct {
	AuthService  auth.AuthServiceInterface
	TokenService auth.TokenInterface
	LocationApp  application.LocationApplicationInterface
}

// NewLocations returns a new instance of Locations
func NewLocations(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, locationApp application.LocationApplicationInterface) *Locations {
	return &Locations{
		AuthService:  authService,
		TokenService: tokenService,
		LocationApp:  locationApp,
	}
}

// GetAllLocations retrieves a paginated list of all locations.
func (l *Locations) GetAllLocations(ctx *gin.Context) {
	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	// Get the location count from the location application service.
	count, err := l.LocationApp.CountLocations()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the locations from the location application service.
	locations, err := l.LocationApp.GetAllLocations(page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(locations) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No locations found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(locations) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var locationPublicData []interface{}

	for _, location := range locations {
		locationPublicData = append(locationPublicData, location.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = locationPublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the locations as a response.
	response.SendOK(ctx, data, "")
}

// GetLocationByID retrieves a single location by ID.
func (l *Locations) GetLocationByID(ctx *gin.Context) {
	// Parse the location ID from the URL parameter.
	locationID, err := strconv.ParseUint(ctx.Param("location_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid location ID."))
		return
	}

	// Get the location from the location application service.
	location, err := l.LocationApp.GetLocationByID(locationID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Location not found."))
		return
	}

	locationPublicData := location.PublicData(language.GetLanguage(ctx))

	// Send the location as a response.
	response.SendOK(ctx, locationPublicData, "")
}

// GetLocationByCoordinates retrieves a location based on the provided latitude, longitude, and radius.
func (l *Locations) GetLocationByCoordinates(ctx *gin.Context) {
	// Parse the latitude from the URL parameter.
	latitude, err := strconv.ParseFloat(ctx.Param("latitude"), 64)
	if err != nil {
		// Send an error response if the latitude is invalid.
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid latitude."))
		return
	}

	// Parse the longitude from the URL parameter.
	longitude, err := strconv.ParseFloat(ctx.Param("longitude"), 64)
	if err != nil {
		// Send an error response if the longitude is invalid.
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid longitude."))
		return
	}

	// Parse the radius from the URL parameter.
	radius, err := strconv.ParseFloat(ctx.Param("radius"), 64)
	if err != nil {
		// Send an error response if the radius is invalid.
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid radius."))
		return
	}

	// Get the location by coordinates from the LocationApp using the parsed latitude, longitude, and radius.
	location, err := l.LocationApp.GetLocationByCoordinates(longitude, latitude, &radius)

	// If there's an error or the location is not found, send an error response.
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Location not found."))
		return
	}

	locationPublicData := location.PublicData(language.GetLanguage(ctx))

	// Send the location as a response.
	response.SendOK(ctx, locationPublicData, "")
}
