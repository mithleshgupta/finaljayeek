package interfaces

import (
	"fmt"
	"strconv"

	"github.com/OmarBader7/web-service-jayeek/application"
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/auth"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/chat"
	"github.com/OmarBader7/web-service-jayeek/pkg/language"
	"github.com/OmarBader7/web-service-jayeek/pkg/pagination"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
)

// Offers holds the offer-related application interfaces
type Offers struct {
	AuthService  auth.AuthServiceInterface
	TokenService auth.TokenInterface
	ChatService  chat.ChatServiceInterface
	UserApp      application.UserApplicationInterface
	OfferApp     application.OfferApplicationInterface
	OrderApp     application.OrderApplicationInterface
	DriverApp    application.DriverApplicationInterface
}

// NewOffers returns a new instance of Offers
func NewOffers(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, chatService chat.ChatServiceInterface, userApp application.UserApplicationInterface, offerApp application.OfferApplicationInterface, orderApp application.OrderApplicationInterface, driverApp application.DriverApplicationInterface) *Offers {
	return &Offers{
		AuthService:  authService,
		TokenService: tokenService,
		ChatService:  chatService,
		UserApp:      userApp,
		OfferApp:     offerApp,
		OrderApp:     orderApp,
		DriverApp:    driverApp,
	}
}

// GetAllOffers retrieves a paginated list of all offers.
func (o *Offers) GetAllOffers(ctx *gin.Context) {
	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	// Extract the token metadata from the request
	metadata, err := o.TokenService.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Fetch the authenticated user's ID from the auth service
	userID, err := o.AuthService.FetchAuth(metadata.AccessTokenUUID)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Get the user from the user application service
	user, err := o.UserApp.GetUserByID(userID)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Get the offer count from the offer application service.
	count, err := o.OfferApp.CountOffersByStatusAndUserID(entity.OfferStatusPending, user.ID)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the offers from the offer application service.
	offers, err := o.OfferApp.GetAllOffersByStatusAndUserID(entity.OfferStatusPending, user.ID, page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(offers) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No offers found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(offers) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var offerPublicData []interface{}

	for _, offer := range offers {
		offerPublicData = append(offerPublicData, offer.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = offerPublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the offers as a response.
	response.SendOK(ctx, data, "")
}

func (o *Offers) AcceptOfferByID(ctx *gin.Context) {
	// Extract the token metadata from the request
	metadata, err := o.TokenService.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Fetch the authenticated user's ID from the auth service
	userID, err := o.AuthService.FetchAuth(metadata.AccessTokenUUID)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Get the user from the user application service
	user, err := o.UserApp.GetUserByID(userID)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Parse the offer ID from the URL parameter.
	offerID, err := strconv.ParseUint(ctx.Param("offer_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid offer ID."))
		return
	}

	// Get the offer from the offer application service.
	offer, err := o.OfferApp.GetOfferByIDAndUserID(offerID, user.ID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Offer not found."))
		return
	}

	// Get the driver from the driver application service
	driver, err := o.DriverApp.GetDriverByID(offer.DriverID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Driver not found."))
		return
	}

	// Get the order from the order application service
	order, err := o.OrderApp.GetOrderByIDAndUserID(offer.OrderID, user.ID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Order not found."))
		return
	}

	offer.Status = entity.OfferStatusAccepted

	if _, err := o.OfferApp.UpdateOfferByID(offer.ID, offer); err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	order.DriverID = driver.ID
	order.Amount = &offer.Amount
	order.Status = entity.OrderAcceptedStatus

	_, err = o.OrderApp.UpdateOrderByID(order.ID, order)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	offers, err := o.OfferApp.GetAllOffersByStatusAndOrderID(entity.OfferStatusPending, order.ID)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	for _, offer := range offers {
		offer.Status = entity.OfferStatusDeclined
		if _, err := o.OfferApp.UpdateOfferByID(offer.ID, &offer); err != nil {
			response.SendInternalServerError(ctx, err.Error())
			return
		}
	}

	err = o.ChatService.CreateChannel(fmt.Sprintf("order-%d", order.ID), fmt.Sprintf("client-%d", user.ID))

	if err != nil {
		response.SendInternalServerError(ctx, ginI18n.MustGetMessage("Failed to create chat channel."))
		return
	}

	err = o.ChatService.AddMember(fmt.Sprintf("order-%d", order.ID), fmt.Sprintf("driver-%d", driver.User.ID))

	if err != nil {
		response.SendInternalServerError(ctx, ginI18n.MustGetMessage("Failed to add members to the chat channel."))
		return
	}

	// Get the recipient from the user application service
	if recipient, err := o.UserApp.GetUserByPhone(order.RecipientPhoneNumber); err == nil {
		order.RecipientID = recipient.ID

		_, err = o.OrderApp.UpdateOrderByID(order.ID, order)
		if err != nil {
			response.SendInternalServerError(ctx, err.Error())
			return
		}

		err = o.ChatService.AddMember(fmt.Sprintf("order-%d", order.ID), fmt.Sprintf("client-%d", recipient.ID))

		if err != nil {
			response.SendInternalServerError(ctx, ginI18n.MustGetMessage("Failed to remove members from the chat channel."))
			return
		}
	}

	response.SendOK(ctx, offer.PublicData(language.GetLanguage(ctx)), "")
}

func (o *Offers) DeclineOfferByID(ctx *gin.Context) {
	// Extract the token metadata from the request
	metadata, err := o.TokenService.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Fetch the authenticated user's ID from the auth service
	userID, err := o.AuthService.FetchAuth(metadata.AccessTokenUUID)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Get the user from the user application service
	user, err := o.UserApp.GetUserByID(userID)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Parse the offer ID from the URL parameter.
	offerID, err := strconv.ParseUint(ctx.Param("offer_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid offer ID."))
		return
	}

	// Get the offer from the offer application service.
	offer, err := o.OfferApp.GetOfferByIDAndUserID(offerID, user.ID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Offer not found."))
		return
	}

	offer.Status = entity.OfferStatusDeclined

	if _, err := o.OfferApp.UpdateOfferByID(offer.ID, offer); err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	response.SendOK(ctx, offer.PublicData(language.GetLanguage(ctx)), "")
}
