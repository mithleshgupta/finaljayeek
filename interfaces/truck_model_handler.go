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

// TruckModels holds the truck model-related application interfaces
type TruckModels struct {
	AuthService   auth.AuthServiceInterface
	TokenService  auth.TokenInterface
	TruckModelApp application.TruckModelApplicationInterface
}

// NewTruckModels returns a new instance of TruckModels
func NewTruckModels(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, truckModelApp application.TruckModelApplicationInterface) *TruckModels {
	return &TruckModels{
		AuthService:   authService,
		TokenService:  tokenService,
		TruckModelApp: truckModelApp,
	}
}

// GetAllTruckModels retrieves a paginated list of all truck models.
func (t *TruckModels) GetAllTruckModels(ctx *gin.Context) {
	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	// Get the truck model count from the truck model application service.
	count, err := t.TruckModelApp.CountTruckModels()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the truckModels from the truck model application service.
	truckModels, err := t.TruckModelApp.GetAllTruckModels(page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(truckModels) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No truck models found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(truckModels) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var truckModelPublicData []interface{}

	for _, truckModel := range truckModels {
		truckModelPublicData = append(truckModelPublicData, truckModel.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = truckModelPublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the truckModels as a response.
	response.SendOK(ctx, data, "")
}

// GetTruckModelByID retrieves a single truck model by ID.
func (t *TruckModels) GetTruckModelByID(ctx *gin.Context) {
	// Parse the truck model ID from the URL parameter.
	truckModelID, err := strconv.ParseUint(ctx.Param("truck_model_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid truck model ID."))
		return
	}

	// Get the truck model from the truck model application service.
	truckModel, err := t.TruckModelApp.GetTruckModelByID(truckModelID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Truck model not found."))
		return
	}

	truckModelPublicData := truckModel.PublicData(language.GetLanguage(ctx))

	// Send the truck model as a response.
	response.SendOK(ctx, truckModelPublicData, "")
}
