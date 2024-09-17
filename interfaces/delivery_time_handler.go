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

// DeliveryTimes holds the deliveryTime-related application interfaces
type DeliveryTimes struct {
	AuthService     auth.AuthServiceInterface
	TokenService    auth.TokenInterface
	DeliveryTimeApp application.DeliveryTimeApplicationInterface
}

// NewDeliveryTimes returns a new instance of DeliveryTimes
func NewDeliveryTimes(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, deliveryTimeApp application.DeliveryTimeApplicationInterface) *DeliveryTimes {
	return &DeliveryTimes{
		AuthService:     authService,
		TokenService:    tokenService,
		DeliveryTimeApp: deliveryTimeApp,
	}
}

// GetAllDeliveryTimes retrieves a paginated list of all deliveryTimes.
func (d *DeliveryTimes) GetAllDeliveryTimes(ctx *gin.Context) {
	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	// Get the delivery time count from the delivery time application service.
	count, err := d.DeliveryTimeApp.CountDeliveryTimes()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the deliveryTimes from the delivery time application service.
	deliveryTimes, err := d.DeliveryTimeApp.GetAllDeliveryTimes(page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(deliveryTimes) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No delivery times found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(deliveryTimes) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var deliveryTimePublicData []interface{}

	for _, deliveryTime := range deliveryTimes {
		deliveryTimePublicData = append(deliveryTimePublicData, deliveryTime.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = deliveryTimePublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the deliveryTimes as a response.
	response.SendOK(ctx, data, "")
}

// GetDeliveryTimeByID retrieves a single delivery time by ID.
func (d *DeliveryTimes) GetDeliveryTimeByID(ctx *gin.Context) {
	// Parse the delivery time ID from the URL parameter.
	deliveryTimeID, err := strconv.ParseUint(ctx.Param("delivery_time_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid delivery time ID."))
		return
	}

	// Get the delivery time from the delivery time application service.
	deliveryTime, err := d.DeliveryTimeApp.GetDeliveryTimeByID(deliveryTimeID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Delivery time not found."))
		return
	}

	deliveryTimePublicData := deliveryTime.PublicData(language.GetLanguage(ctx))

	// Send the delivery time as a response.
	response.SendOK(ctx, deliveryTimePublicData, "")
}
