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

// TransportationModes holds the transportation mode-related application interfaces
type TransportationModes struct {
	AuthService           auth.AuthServiceInterface
	TokenService          auth.TokenInterface
	TransportationModeApp application.TransportationModeApplicationInterface
}

// NewTransportationModes returns a new instance of TransportationModes
func NewTransportationModes(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, transportationModeApp application.TransportationModeApplicationInterface) *TransportationModes {
	return &TransportationModes{
		AuthService:           authService,
		TokenService:          tokenService,
		TransportationModeApp: transportationModeApp,
	}
}

// GetAllTransportationModes retrieves a paginated list of all transportation modes.
func (t *TransportationModes) GetAllTransportationModes(ctx *gin.Context) {
	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	// Get the transportation mode count from the transportation mode application service.
	count, err := t.TransportationModeApp.CountTransportationModes()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the transportationModes from the transportation mode application service.
	transportationModes, err := t.TransportationModeApp.GetAllTransportationModes(page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(transportationModes) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No transportation modes found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(transportationModes) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var transportationModePublicData []interface{}

	for _, transportationMode := range transportationModes {
		transportationModePublicData = append(transportationModePublicData, transportationMode.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = transportationModePublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the transportationModes as a response.
	response.SendOK(ctx, data, "")
}

// GetTransportationModeByID retrieves a single transportation mode by ID.
func (t *TransportationModes) GetTransportationModeByID(ctx *gin.Context) {
	// Parse the transportation mode ID from the URL parameter.
	transportationModeID, err := strconv.ParseUint(ctx.Param("transportation_mode_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid transportation mode ID."))
		return
	}

	// Get the transportation mode from the transportation mode application service.
	transportationMode, err := t.TransportationModeApp.GetTransportationModeByID(transportationModeID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Transportation mode not found."))
		return
	}

	transportationModePublicData := transportationMode.PublicData(language.GetLanguage(ctx))

	// Send the transportation mode as a response.
	response.SendOK(ctx, transportationModePublicData, "")
}
