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

// ShipmentContents holds the shipment content-related application interfaces
type ShipmentContents struct {
	AuthService        auth.AuthServiceInterface
	TokenService       auth.TokenInterface
	ShipmentContentApp application.ShipmentContentApplicationInterface
}

// NewShipmentContents returns a new instance of ShipmentContents
func NewShipmentContents(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, shipmentContentApp application.ShipmentContentApplicationInterface) *ShipmentContents {
	return &ShipmentContents{
		AuthService:        authService,
		TokenService:       tokenService,
		ShipmentContentApp: shipmentContentApp,
	}
}

// GetAllShipmentContents retrieves a paginated list of all shipment contents.
func (s *ShipmentContents) GetAllShipmentContents(ctx *gin.Context) {
	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	// Get the shipment content count from the shipment content application service.
	count, err := s.ShipmentContentApp.CountShipmentContents()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the shipmentContents from the shipment content application service.
	shipmentContents, err := s.ShipmentContentApp.GetAllShipmentContents(page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(shipmentContents) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No shipment contents found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(shipmentContents) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var shipmentContentPublicData []interface{}

	for _, shipmentContent := range shipmentContents {
		shipmentContentPublicData = append(shipmentContentPublicData, shipmentContent.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = shipmentContentPublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the shipmentContents as a response.
	response.SendOK(ctx, data, "")
}

// GetShipmentContentByID retrieves a single shipment content by ID.
func (s *ShipmentContents) GetShipmentContentByID(ctx *gin.Context) {
	// Parse the shipment content ID from the URL parameter.
	shipmentContentID, err := strconv.ParseUint(ctx.Param("shipment_content_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid shipment content ID."))
		return
	}

	// Get the shipment content from the shipment content application service.
	shipmentContent, err := s.ShipmentContentApp.GetShipmentContentByID(shipmentContentID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Shipment content not found."))
		return
	}

	shipmentContentPublicData := shipmentContent.PublicData(language.GetLanguage(ctx))

	// Send the shipment content as a response.
	response.SendOK(ctx, shipmentContentPublicData, "")
}
