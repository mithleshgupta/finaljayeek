package interfaces

import (
	"fmt"
	"strconv"

	"github.com/OmarBader7/web-service-jayeek/application"
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/auth"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/chat"
	"github.com/OmarBader7/web-service-jayeek/pkg/geoutil"
	"github.com/OmarBader7/web-service-jayeek/pkg/language"
	"github.com/OmarBader7/web-service-jayeek/pkg/pagination"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	"github.com/OmarBader7/web-service-jayeek/pkg/validator"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
)

// Orders holds the order-related application interfaces
type Orders struct {
	AuthService        auth.AuthServiceInterface
	TokenService       auth.TokenInterface
	ChatService        chat.ChatServiceInterface
	OrderApp           application.OrderApplicationInterface
	UserApp            application.UserApplicationInterface
	CategoryApp        application.CategoryApplicationInterface
	LocationApp        application.LocationApplicationInterface
	DriverApp          application.DriverApplicationInterface
	SizeApp            application.SizeApplicationInterface
	TruckTypeApp       application.TruckTypeApplicationInterface
	TruckModelApp      application.TruckModelApplicationInterface
	DeliveryTimeApp    application.DeliveryTimeApplicationInterface
	ShipmentContentApp application.ShipmentContentApplicationInterface
	ExtraServiceApp    application.ExtraServiceApplicationInterface
	BalanceApp         application.BalanceApplicationInterface
	SettingApp         application.SettingApplicationInterface
	OfferApp           application.OfferApplicationInterface
	RatingApp          application.RatingApplicationInterface
}

// NewOrders returns a new instance of Orders
func NewOrders(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, chatService chat.ChatServiceInterface, orderApp application.OrderApplicationInterface, userApp application.UserApplicationInterface, categoryApp application.CategoryApplicationInterface, locationApp application.LocationApplicationInterface, driverApp application.DriverApplicationInterface, sizeApp application.SizeApplicationInterface, truckTypeApp application.TruckTypeApplicationInterface, truckModelApp application.TruckModelApplicationInterface, deliveryTimeApp application.DeliveryTimeApplicationInterface, shipmentContentApp application.ShipmentContentApplicationInterface, extraServiceApp application.ExtraServiceApplicationInterface, balanceApp application.BalanceApplicationInterface, settingApp application.SettingApplicationInterface, offerApp application.OfferApplicationInterface, ratingApp application.RatingApplicationInterface) *Orders {
	return &Orders{
		AuthService:        authService,
		TokenService:       tokenService,
		ChatService:        chatService,
		OrderApp:           orderApp,
		UserApp:            userApp,
		CategoryApp:        categoryApp,
		LocationApp:        locationApp,
		DriverApp:          driverApp,
		SizeApp:            sizeApp,
		TruckTypeApp:       truckTypeApp,
		TruckModelApp:      truckModelApp,
		DeliveryTimeApp:    deliveryTimeApp,
		ShipmentContentApp: shipmentContentApp,
		ExtraServiceApp:    extraServiceApp,
		BalanceApp:         balanceApp,
		SettingApp:         settingApp,
		OfferApp:           offerApp,
		RatingApp:          ratingApp,
	}
}

// CreateOrder handles the creation of a new order
func (d *Orders) CreateOrder(c *gin.Context) {
	var order entity.Order

	// Bind the JSON body of the request to the Order struct
	if err := c.ShouldBindJSON(&order); err != nil {
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

	// Validate all fields of the order struct except for UserID, LocationID, Amount, Status, Location, User, Driver, Recipient, Category, Size, TruckType, TruckModel, DeliveryTime, Destination, and ShipmentContents
	validationErrors, _ := validator.ValidateExcept(c, &order, "UserID", "LocationID", "Amount", "Status", "Location", "User", "Driver", "Recipient", "Category", "Size", "TruckType", "TruckModel", "DeliveryTime", "Destination", "ShipmentContents", "ExtraServices")
	if validationErrors != nil {
		response.SendUnprocessableEntity(c, validationErrors, "")
		return
	}

	// Get the category by its ID
	category, err := d.CategoryApp.GetCategoryByID(order.CategoryID)
	if err != nil {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Category not found."))
		return
	}

	// Get the destination by its ID
	destination, err := d.LocationApp.GetLocationByID(order.DestinationID)
	if err != nil {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Destination not found."))
		return
	}

	var driver *entity.Driver

	if order.DriverID != 0 {
		// Get the driver by its ID
		driver, err = d.DriverApp.GetDriverByID(order.DriverID)
		if err != nil {
			response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Driver not found."))
			return
		}
	}

	var size *entity.Size
	var truckType *entity.TruckType
	var truckModel *entity.TruckModel

	if category.IsTruck != nil && *category.IsTruck {
		// Get the truck type by its ID
		truckType, err = d.TruckTypeApp.GetTruckTypeByID(order.TruckTypeID)
		if err != nil {
			response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Truck type not found."))
			return
		}

		// Get the truck type by its ID
		truckModel, err = d.TruckModelApp.GetTruckModelByID(order.TruckModelID)
		if err != nil {
			response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Truck model not found."))
			return
		}
	} else {
		// Get the size by its ID
		size, err = d.SizeApp.GetSizeByID(order.SizeID)
		if err != nil {
			response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Size not found."))
			return
		}
	}

	// Get the delivery time by its ID
	deliveryTime, err := d.DeliveryTimeApp.GetDeliveryTimeByID(order.DeliveryTimeID)
	if err != nil {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Delivery time not found."))
		return
	}

	shipmentContents := make([]entity.ShipmentContent, len(order.ShipmentContentIDs))
	for i, shipmentContentID := range order.ShipmentContentIDs {
		shipmentContent, err := d.ShipmentContentApp.GetShipmentContentByID(shipmentContentID)
		if err != nil {
			response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Shipment content not found."))
			return
		}
		shipmentContents[i] = *shipmentContent
	}
	order.ShipmentContents = shipmentContents

	extraServices := make([]entity.ExtraService, len(*order.ExtraServiceIDs))
	for i, extraServiceID := range *order.ExtraServiceIDs {
		extraService, err := d.ExtraServiceApp.GetExtraServiceByID(extraServiceID)
		if err != nil {
			response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Extra service not found."))
			return
		}
		extraServices[i] = *extraService
	}
	order.ExtraServices = &extraServices

	location, err := d.LocationApp.GetLocationByCoordinates(order.Longitude, order.Latitude, nil)
	if err != nil {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Location not found."))
		return
	}

	// Set the LocationID, UserID, DriverID, SizeID, DeliveryTimeID, DestinationID and Status for the order struct
	order.LocationID = location.ID
	order.UserID = user.ID
	if order.DriverID != 0 {
		order.DriverID = driver.ID
	}

	if category.IsTruck != nil && *category.IsTruck {
		if order.TruckTypeID != 0 {
			order.TruckTypeID = truckType.ID
		}
		if order.TruckModelID != 0 {
			order.TruckModelID = truckModel.ID
		}
	} else {
		if order.SizeID != 0 {
			order.SizeID = size.ID
		}
	}

	order.DeliveryTimeID = deliveryTime.ID
	order.DestinationID = destination.ID
	order.Status = entity.OrderCreatedStatus

	// Create the new order
	createdOrder, err := d.OrderApp.CreateOrder(&order)
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	response.SendOK(c, createdOrder.PublicData(language.GetLanguage(c)), "")
}

// GetAllOrders retrieves a paginated list of all orders.
func (o *Orders) GetAllOrders(ctx *gin.Context) {
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

	// Get the driver from the driver application service
	driver, err := o.DriverApp.GetDriverByUserID(user.ID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Driver not found."))
		return
	}

	statusesToExclude := []entity.OrderStatus{
		entity.OrderCreatedStatus,
	}

	// Get the order count from the order application service.
	count, err := o.OrderApp.CountOrdersByDriverIDExcludingStatus(driver.ID, statusesToExclude)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the orders from the order application service.
	orders, err := o.OrderApp.GetAllOrdersByDriverIDExcludingStatus(driver.ID, statusesToExclude, page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(orders) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No orders found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(orders) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var orderPublicData []interface{}

	for _, order := range orders {
		orderPublicData = append(orderPublicData, order.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = orderPublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the orders as a response.
	response.SendOK(ctx, data, "")
}

// GetAllOrders retrieves a paginated list of all purchases.
func (o *Orders) GetAllPurchases(ctx *gin.Context) {
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

	statusesToExclude := []entity.OrderStatus{
		entity.OrderCreatedStatus,
	}

	// Get the order count from the order application service.
	count, err := o.OrderApp.CountOrdersByUserIDAndRecipientIDExcludingStatus(user.ID, user.ID, statusesToExclude)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the orders from the order application service.
	orders, err := o.OrderApp.GetAllOrdersByUserIDAndRecipientIDExcludingStatus(user.ID, user.ID, statusesToExclude, page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(orders) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No orders found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(orders) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var orderPublicData []interface{}

	for _, order := range orders {
		if order.User.ID == user.ID {
			order.IsSender = true
		}
		if order.Recipient.ID == user.ID {
			order.IsReceiver = true
		}
		orderPublicData = append(orderPublicData, order.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = orderPublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the orders as a response.
	response.SendOK(ctx, data, "")
}

// GetAllDriverPools retrieves a paginated list of all orders.
func (o *Orders) GetAllDriverPoolsByCategoryID(ctx *gin.Context) {
	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	orderBy := ctx.Query("order_by")
	byArrival := ctx.Query("by_arrival")

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

	// Parse the order ID from the URL parameter.
	categoryID, err := strconv.ParseUint(ctx.Param("category_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid category ID."))
		return
	}

	// Get the user from the user application service
	user, err := o.UserApp.GetUserByID(userID)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Get the driver from the driver application service
	driver, err := o.DriverApp.GetDriverByUserID(user.ID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Driver not found."))
		return
	}

	// Get the order count from the order application service.
	count, err := o.OrderApp.CountDriverPoolsByDriverIDAndCategoryID(driver.ID, categoryID, &orderBy)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the orders from the order application service.
	orders, err := o.OrderApp.GetAllDriverPoolsByDriverIDAndCategoryID(driver.ID, categoryID, page, perPage, &orderBy, &byArrival)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(orders) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No orders found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(orders) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var orderPublicData []interface{}

	for _, order := range orders {
		orderPublicData = append(orderPublicData, order.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = orderPublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the orders as a response.
	response.SendOK(ctx, data, "")
}

func (o *Orders) AcceptDriverPoolByID(ctx *gin.Context) {
	var offer entity.Offer

	// Bind the JSON body of the request to the Order struct
	if err := ctx.ShouldBindJSON(&offer); err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid request body"))
		return
	}

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

	// Validate all fields of the offer struct except for DriverID, OrderID, Status, Driver, and Order
	validationErrors, _ := validator.ValidateExcept(ctx, &offer, "DriverID", "OrderID", "Status", "Driver", "Order")
	if validationErrors != nil {
		response.SendUnprocessableEntity(ctx, validationErrors, "")
		return
	}

	// Get the driver from the driver application service
	driver, err := o.DriverApp.GetDriverByUserID(user.ID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Driver not found."))
		return
	}

	// Parse the order ID from the URL parameter.
	orderID, err := strconv.ParseUint(ctx.Param("order_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid order ID."))
		return
	}

	// Get the order from the order application service.
	order, err := o.OrderApp.GetOrderByID(orderID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Order not found."))
		return
	}

	// Get the order from the order application service.
	orderDriverPool, err := o.OrderApp.GetOrderDriverPoolByOrderIDAndDriverID(order.ID, driver.ID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Order driver pool not found."))
		return
	}

	offersCount, _ := o.OfferApp.CountOffersByDriverIDAndStatus(driver.ID, entity.OfferStatusPending)

	statusesToExclude := []entity.OrderStatus{
		entity.OrderCompletedStatus,
		entity.OrderCanceledStatus,
		entity.ShipmentReturnedStatus,
	}

	ordersCount, _ := o.OrderApp.CountOrdersByDriverIDExcludingStatus(driver.ID, statusesToExclude)

	maxOrdersPerTrip, err := o.getMaxOrdersPerTrip()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
	}

	if offersCount+ordersCount >= maxOrdersPerTrip {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Max order limit reached."))
		return
	}

	dist := geoutil.CalculateDistance(driver.Latitude, driver.Longitude, order.Latitude, order.Longitude)
	if dist > 10.0 {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Driver out of range."))
		return
	}

	orderDriverPool.Status = entity.AcceptedStatus

	if _, err := o.OrderApp.UpdateOrderDriverPoolByOrderIDAndDriverID(order.ID, driver.ID, orderDriverPool); err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	offer.DriverID = driver.ID
	offer.OrderID = order.ID
	offer.Status = entity.OfferStatusPending

	// Create the new offer
	_, err = o.OfferApp.CreateOffer(&offer)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	response.SendOK(ctx, order.PublicData(language.GetLanguage(ctx)), "")
}

func (o *Orders) DeclineDriverPoolByID(ctx *gin.Context) {
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

	// Get the driver from the driver application service
	driver, err := o.DriverApp.GetDriverByUserID(user.ID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Driver not found."))
		return
	}

	// Parse the order ID from the URL parameter.
	orderID, err := strconv.ParseUint(ctx.Param("order_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid order ID."))
		return
	}

	// Get the order from the order application service.
	order, err := o.OrderApp.GetOrderByID(orderID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Order not found."))
		return
	}

	// Get the order from the order application service.
	orderDriverPool, err := o.OrderApp.GetOrderDriverPoolByOrderIDAndDriverID(order.ID, driver.ID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Order driver pool not found."))
		return
	}

	orderDriverPool.Status = entity.RejectedStatus

	if _, err := o.OrderApp.UpdateOrderDriverPoolByOrderIDAndDriverID(order.ID, driver.ID, orderDriverPool); err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	response.SendOK(ctx, nil, "")
}

func (o *Orders) CancelOrderByID(ctx *gin.Context) {
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

	// Parse the order ID from the URL parameter.
	orderID, err := strconv.ParseUint(ctx.Param("order_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid order ID."))
		return
	}

	// Get the order from the order application service.
	order, err := o.OrderApp.GetOrderByIDAndUserID(orderID, user.ID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Order not found."))
		return
	}

	order.Status = entity.OrderCanceledStatus

	updatedOrder, err := o.OrderApp.UpdateOrderByID(order.ID, order)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if order.RecipientID != 0 {
		// Get the recipient from the user application service
		recipient, err := o.UserApp.GetUserByID(order.RecipientID)
		if err != nil {
			response.SendNotFound(ctx, ginI18n.MustGetMessage("Recipient not found."))
			return
		}

		err = o.ChatService.RemoveMember(fmt.Sprintf("order-%d", order.ID), fmt.Sprintf("client-%d", recipient.ID))

		if err != nil {
			response.SendInternalServerError(ctx, ginI18n.MustGetMessage("Failed to remove members from the chat channel."))
			return
		}
	}

	err = o.ChatService.RemoveMember(fmt.Sprintf("order-%d", order.ID), fmt.Sprintf("client-%d", userID))

	if err != nil {
		response.SendInternalServerError(ctx, ginI18n.MustGetMessage("Failed to remove members from the chat channel."))
		return
	}

	response.SendOK(ctx, updatedOrder.PublicData(language.GetLanguage(ctx)), "")
}

func (o *Orders) DeliverOrderByID(ctx *gin.Context) {
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

	// Get the driver from the driver application service
	driver, err := o.DriverApp.GetDriverByUserID(user.ID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Driver not found."))
		return
	}

	// Parse the order ID from the URL parameter.
	orderID, err := strconv.ParseUint(ctx.Param("order_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid order ID."))
		return
	}

	// Get the order from the order application service.
	order, err := o.OrderApp.GetOrderByIDAndDriverID(orderID, driver.ID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Order not found."))
		return
	}

	order.Status = entity.ShipmentDeliveredStatus

	updatedOrder, err := o.OrderApp.UpdateOrderByID(order.ID, order)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	var balance entity.Balance

	balance.OrderID = updatedOrder.ID
	balance.DriverID = updatedOrder.Driver.ID

	// Create the new balance
	_, err = o.BalanceApp.CreateBalance(&balance)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	err = o.ChatService.RemoveMember(fmt.Sprintf("order-%d", order.ID), fmt.Sprintf("driver-%d", driver.User.ID))

	if err != nil {
		response.SendInternalServerError(ctx, ginI18n.MustGetMessage("Failed to remove members from the chat channel."))
		return
	}

	if order.RecipientID != 0 {
		// Get the recipient from the user application service
		recipient, err := o.UserApp.GetUserByID(order.RecipientID)
		if err != nil {
			response.SendNotFound(ctx, ginI18n.MustGetMessage("Recipient not found."))
			return
		}

		err = o.ChatService.RemoveMember(fmt.Sprintf("order-%d", order.ID), fmt.Sprintf("client-%d", recipient.ID))

		if err != nil {
			response.SendInternalServerError(ctx, ginI18n.MustGetMessage("Failed to remove members from the chat channel."))
			return
		}
	}

	err = o.ChatService.RemoveMember(fmt.Sprintf("order-%d", order.ID), fmt.Sprintf("driver-%d", driver.User.ID))

	if err != nil {
		response.SendInternalServerError(ctx, ginI18n.MustGetMessage("Failed to remove members from the chat channel."))
		return
	}

	response.SendOK(ctx, updatedOrder.PublicData(language.GetLanguage(ctx)), "")
}

func (o *Orders) PickupOrderByID(ctx *gin.Context) {
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

	// Get the driver from the driver application service
	driver, err := o.DriverApp.GetDriverByUserID(user.ID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Driver not found."))
		return
	}

	// Parse the order ID from the URL parameter.
	orderID, err := strconv.ParseUint(ctx.Param("order_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid order ID."))
		return
	}

	// Get the order from the order application service.
	order, err := o.OrderApp.GetOrderByIDAndDriverID(orderID, driver.ID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Order not found."))
		return
	}

	order.Status = entity.ShipmentPickedUpStatus

	updatedOrder, err := o.OrderApp.UpdateOrderByID(order.ID, order)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	response.SendOK(ctx, updatedOrder.PublicData(language.GetLanguage(ctx)), "")
}

// GetOrderByID retrieves a single order by ID.
func (l *Orders) GetOrderByID(ctx *gin.Context) {
	// Parse the order ID from the URL parameter.
	orderID, err := strconv.ParseUint(ctx.Param("order_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid order ID."))
		return
	}

	// Get the order from the order application service.
	order, err := l.OrderApp.GetOrderByID(orderID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Order not found."))
		return
	}

	// Send the order as a response.
	response.SendOK(ctx, order.PublicData(language.GetLanguage(ctx)).(*entity.OrderPublicData), "")
}

// RateOrderByID handles the creation of a new order
func (d *Orders) RateOrderByID(c *gin.Context) {
	var rating entity.Rating

	// Parse the order ID from the URL parameter.
	orderID, err := strconv.ParseUint(c.Param("order_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Invalid order ID."))
		return
	}

	// Bind the JSON body of the request to the Order struct
	if err := c.ShouldBindJSON(&rating); err != nil {
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

	// Get the order from the order application service.
	order, err := d.OrderApp.GetOrderByIDAndUserID(orderID, user.ID)
	if err != nil {
		response.SendNotFound(c, ginI18n.MustGetMessage("Order not found."))
		return
	}

	validationErrors, _ := validator.ValidateExcept(c, &rating, "OrderID", "UserID")
	if validationErrors != nil {
		response.SendUnprocessableEntity(c, validationErrors, "")
		return
	}

	rating.OrderID = order.ID
	rating.UserID = user.ID

	// Create the new order
	_, err = d.RatingApp.CreateRating(&rating)
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	err = d.ChatService.RemoveMember(fmt.Sprintf("order-%d", order.ID), fmt.Sprintf("client-%d", userID))

	if err != nil {
		response.SendInternalServerError(c, ginI18n.MustGetMessage("Failed to remove members from the chat channel."))
		return
	}

	response.SendOK(c, order.PublicData(language.GetLanguage(c)), "")
}

// ReturnOrderByID handles the creation of a new order
func (d *Orders) ReturnOrderByID(c *gin.Context) {
	var newOrder entity.Order

	// Parse the order ID from the URL parameter.
	orderID, err := strconv.ParseUint(c.Param("order_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Invalid order ID."))
		return
	}

	// Bind the JSON body of the request to the Order struct
	if err := c.ShouldBindJSON(&newOrder); err != nil {
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

	// Get the order from the order application service.
	order, err := d.OrderApp.GetOrderByIDAndRecipientID(orderID, user.ID)
	if err != nil {
		response.SendNotFound(c, ginI18n.MustGetMessage("Order not found."))
		return
	}

	if validationErrors, _ := validator.ValidatePartial(c, &newOrder, "Longitude", "Latitude"); validationErrors != nil {
		response.SendUnprocessableEntity(c, validationErrors, "")
		return
	}

	// Get the category by its ID
	category, err := d.CategoryApp.GetCategoryByID(order.CategoryID)
	if err != nil {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Category not found."))
		return
	}

	// Get the destination by its ID
	destination, err := d.LocationApp.GetLocationByID(order.LocationID)
	if err != nil {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Destination not found."))
		return
	}

	var size *entity.Size
	var truckType *entity.TruckType
	var truckModel *entity.TruckModel

	if category.IsTruck != nil && *category.IsTruck {
		// Get the truck type by its ID
		truckType, err = d.TruckTypeApp.GetTruckTypeByID(order.TruckTypeID)
		if err != nil {
			response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Truck type not found."))
			return
		}

		// Get the truck type by its ID
		truckModel, err = d.TruckModelApp.GetTruckModelByID(order.TruckModelID)
		if err != nil {
			response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Truck model not found."))
			return
		}
	} else {
		// Get the size by its ID
		size, err = d.SizeApp.GetSizeByID(order.SizeID)
		if err != nil {
			response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Size not found."))
			return
		}
	}

	// Get the delivery time by its ID
	deliveryTime, err := d.DeliveryTimeApp.GetDeliveryTimeByID(order.DeliveryTimeID)
	if err != nil {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Delivery time not found."))
		return
	}

	shipmentContents := make([]entity.ShipmentContent, len(order.ShipmentContents))
	for i, shipmentContent := range order.ShipmentContents {
		shipmentContent, err := d.ShipmentContentApp.GetShipmentContentByID(shipmentContent.ID)
		if err != nil {
			response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Shipment content not found."))
			return
		}
		shipmentContents[i] = *shipmentContent
	}
	newOrder.ShipmentContents = shipmentContents

	extraServices := make([]entity.ExtraService, len(*order.ExtraServices))
	for i, extraService := range *order.ExtraServices {
		extraService, err := d.ExtraServiceApp.GetExtraServiceByID(extraService.ID)
		if err != nil {
			response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Extra service not found."))
			return
		}
		extraServices[i] = *extraService
	}
	newOrder.ExtraServices = &extraServices

	location, err := d.LocationApp.GetLocationByCoordinates(newOrder.Longitude, newOrder.Latitude, nil)
	if err != nil {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Location not found."))
		return
	}

	// Set the LocationID, UserID, DriverID, SizeID, DeliveryTimeID, DestinationID and Status for the order struct
	newOrder.CategoryID = category.ID
	newOrder.LocationID = location.ID
	newOrder.UserID = user.ID

	if category.IsTruck != nil && *category.IsTruck {
		if order.TruckTypeID != 0 {
			newOrder.TruckTypeID = truckType.ID
		}
		if order.TruckModelID != 0 {
			newOrder.TruckModelID = truckModel.ID
		}
	} else {
		if order.SizeID != 0 {
			newOrder.SizeID = size.ID
		}
	}

	newOrder.DeliveryTimeID = deliveryTime.ID
	newOrder.DestinationID = destination.ID
	newOrder.Quantity = order.Quantity
	newOrder.RecipientPhoneNumber = order.User.Phone
	newOrder.Status = entity.OrderCreatedStatus

	// Create the new order
	createdOrder, err := d.OrderApp.CreateOrder(&newOrder)
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	order.Status = entity.ShipmentReturnedStatus

	_, err = d.OrderApp.UpdateOrderByID(order.ID, order)
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	response.SendOK(c, createdOrder.PublicData(language.GetLanguage(c)), "")
}

func (o *Orders) getMaxOrdersPerTrip() (int64, error) {
	maxOrdersPerTripStr, err := o.SettingApp.GetSettingByKey("max_orders_per_trip")
	if err != nil {
		return 0, err
	}
	maxOrdersPerTrip, err := strconv.ParseInt(maxOrdersPerTripStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return maxOrdersPerTrip, nil
}
