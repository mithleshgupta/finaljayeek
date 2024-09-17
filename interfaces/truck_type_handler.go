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

// TruckTypes holds the truck type-related application interfaces
type TruckTypes struct {
	AuthService  auth.AuthServiceInterface
	TokenService auth.TokenInterface
	TruckTypeApp application.TruckTypeApplicationInterface
}

// NewTruckTypes returns a new instance of TruckTypes
func NewTruckTypes(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, truckTypeApp application.TruckTypeApplicationInterface) *TruckTypes {
	return &TruckTypes{
		AuthService:  authService,
		TokenService: tokenService,
		TruckTypeApp: truckTypeApp,
	}
}

// GetAllTruckTypes retrieves a paginated list of all truck types.
func (t *TruckTypes) GetAllTruckTypes(ctx *gin.Context) {
	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	// Get the truck type count from the truck type application service.
	count, err := t.TruckTypeApp.CountTruckTypes()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the truckTypes from the truck type application service.
	truckTypes, err := t.TruckTypeApp.GetAllTruckTypes(page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(truckTypes) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No truck types found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(truckTypes) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var truckTypePublicData []interface{}

	for _, truckType := range truckTypes {
		truckTypePublicData = append(truckTypePublicData, truckType.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = truckTypePublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the truckTypes as a response.
	response.SendOK(ctx, data, "")
}

// GetTruckTypeByID retrieves a single truck type by ID.
func (t *TruckTypes) GetTruckTypeByID(ctx *gin.Context) {
	// Parse the truck type ID from the URL parameter.
	truckTypeID, err := strconv.ParseUint(ctx.Param("truck_type_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid truck type ID."))
		return
	}

	// Get the truck type from the truck type application service.
	truckType, err := t.TruckTypeApp.GetTruckTypeByID(truckTypeID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Truck type not found."))
		return
	}

	truckTypePublicData := truckType.PublicData(language.GetLanguage(ctx))

	// Send the truck type as a response.
	response.SendOK(ctx, truckTypePublicData, "")
}
