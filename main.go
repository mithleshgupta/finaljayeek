package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/OmarBader7/web-service-jayeek/infrastructure/auth"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/chat"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/config"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/geo"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/persistence"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/profile"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/user_setting"
	"github.com/OmarBader7/web-service-jayeek/interfaces"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/text/language"
)

func init() {
	// Loads environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
}

func main() {
	// Load the environment variables
	conf := config.NewConfig()

	Host := conf.Host
	Port := conf.Port
	PostgresHost := conf.PostgresHost
	PostgresPort := conf.PostgresPort
	PostgresDatabase := conf.PostgresDatabase
	PostgresUsername := conf.PostgresUsername
	PostgresPassword := conf.PostgresPassword
	PostgresSslMode := conf.PostgresSslMode
	PostgresTimeZone := conf.PostgresTimeZone
	RedisHost := conf.RedisHost
	RedisPassword := conf.RedisPassword
	RedisPort := conf.RedisPort
	streamApiKey := conf.StreamApiKey
	streamApiSecret := conf.StreamApiSecret

	// Create new Postgres repositories
	repositories, err := persistence.NewRepositories(PostgresHost, PostgresPort, PostgresUsername, PostgresPassword, PostgresDatabase, PostgresSslMode, PostgresTimeZone)
	if err != nil {
		log.Fatal("Error creating Postgres repositories: ", err)
	}
	repositories.AutoMigrate()

	var seed = flag.Bool("seed", false, "a bool")

	flag.Parse()
	if *seed {
		repositories.SeedCategories()
		repositories.SeedDeliveryTimes()
		repositories.SeedTruckTypes()
		repositories.SeedTruckModels()
		repositories.SeedLocations()
		repositories.SeedUsers()
		repositories.SeedShipmentContents()
		repositories.SeedExtraServices()
		repositories.SeedSizes()
		repositories.SeedTransportationModes()
		repositories.SeedPages()
		repositories.SeedFAQs()
		repositories.SeedSettings()
	}

	// Create new Redis service
	redisService, err := auth.NewRedisService(RedisHost, RedisPort, RedisPassword)
	if err != nil {
		log.Fatal("Error creating Redis service: ", err)
	}

	streamService, err := chat.NewStreamService(streamApiKey, streamApiSecret)
	if err != nil {
		log.Fatal("Error creating Stream service: ", err)
	}

	// Create new token generator
	tokenGenerator := auth.NewToken()

	// Create new authentication service
	authService := interfaces.NewAuth(redisService.AuthService, tokenGenerator, streamService.ChatService, repositories.User, repositories.Location, repositories.Order, repositories.PasswordReset, repositories.PhoneVerification)

	// Create new user service
	userService := interfaces.NewUsers(redisService.AuthService, tokenGenerator, repositories.User, repositories.Location)

	// Create new device service
	deviceService := interfaces.NewDevices(redisService.AuthService, tokenGenerator, repositories.Device, repositories.User)

	// Create new category service
	categoryService := interfaces.NewCategories(redisService.AuthService, tokenGenerator, repositories.Category)

	// Create new size service
	sizeService := interfaces.NewSizes(redisService.AuthService, tokenGenerator, repositories.Size)

	// Create new cities service

	cityService := interfaces.NewCities(redisService.AuthService, tokenGenerator, repositories.City)

	// Create new location service
	locationService := interfaces.NewLocations(redisService.AuthService, tokenGenerator, repositories.Location)

	// Create new shipment content service
	shipmentContentService := interfaces.NewShipmentContents(redisService.AuthService, tokenGenerator, repositories.ShipmentContent)

	// Create new extra service service
	extraServiceService := interfaces.NewExtraServices(redisService.AuthService, tokenGenerator, repositories.ExtraService)

	// Create new transportation mode service
	transportationModeService := interfaces.NewTransportationModes(redisService.AuthService, tokenGenerator, repositories.TransportationMode)

	// Create new delivery time service
	deliveryTimeService := interfaces.NewDeliveryTimes(redisService.AuthService, tokenGenerator, repositories.DeliveryTime)

	// Create new truck type service
	truckTypeService := interfaces.NewTruckTypes(redisService.AuthService, tokenGenerator, repositories.TruckType)

	// Create new truck model service
	truckModelService := interfaces.NewTruckModels(redisService.AuthService, tokenGenerator, repositories.TruckModel)

	// Create new driver service
	driverService := interfaces.NewDrivers(redisService.AuthService, tokenGenerator, repositories.Driver, repositories.User, repositories.TransportationMode, repositories.IdentityDocument)

	// Create new order service
	orderService := interfaces.NewOrders(redisService.AuthService, tokenGenerator, streamService.ChatService, repositories.Order, repositories.User, repositories.Category, repositories.Location, repositories.Driver, repositories.Size, repositories.TruckType, repositories.TruckModel, repositories.DeliveryTime, repositories.ShipmentContent, repositories.ExtraService, repositories.Balance, repositories.Setting, repositories.Offer, repositories.Rating)

	// Create new offer service
	offerService := interfaces.NewOffers(redisService.AuthService, tokenGenerator, streamService.ChatService, repositories.User, repositories.Offer, repositories.Order, repositories.Driver)

	// Create new page service
	pageService := interfaces.NewPages(repositories.Page)

	// Create new FAQ service
	faqService := interfaces.NewFAQs(repositories.FAQ)

	// Create new chat service
	chatService := interfaces.NewChat(redisService.AuthService, tokenGenerator, streamService.ChatService)

	// Create new profile service
	profileService := interfaces.NewProfile(redisService.AuthService, tokenGenerator, repositories.User, repositories.Driver, repositories.Location, profile.NewProfileService(userService.UserApp))

	// Create new profile service
	userSettingService := interfaces.NewUserSetting(redisService.AuthService, tokenGenerator, repositories.User, user_setting.NewUserSettingService(userService.UserApp))

	// Create new geo service
	geoService := interfaces.NewGeo(geo.NewGeoService())

	// Create new setting service
	settingService := interfaces.NewSettings(redisService.AuthService, tokenGenerator, repositories.Setting)

	// Create new router
	router := gin.Default()

	// Use i18n middleware
	router.Use(ginI18n.Localize(ginI18n.WithBundle(&ginI18n.BundleCfg{
		RootPath:         "./resources",
		AcceptLanguage:   []language.Tag{language.Arabic, language.English},
		DefaultLanguage:  language.English,
		UnmarshalFunc:    json.Unmarshal,
		FormatBundleFile: "json",
	})))

	// Create auth group for authentication routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authService.Register)
		authGroup.POST("/login", authService.Login)
		authGroup.POST("/logout", interfaces.AuthMiddleware(), authService.Logout)
		authGroup.POST("/refresh", authService.Refresh)
		authGroup.POST("/send", authService.SendPhoneVerificationCode)
		authGroup.POST("/verify", authService.VerifyPhoneVerificationCode)

		resetGroup := authGroup.Group("/reset")
		{
			resetGroup.POST("/send", authService.SendPasswordResetCode)
			resetGroup.POST("/verify", authService.VerifyPasswordResetCode)
			resetGroup.PUT("/", authService.Reset)
		}
	}

	userGroup := router.Group("/users")
	{
		userGroup.GET("/", interfaces.AuthMiddleware(), userService.GetAllUsers)
		userGroup.GET("/:user_id", interfaces.AuthMiddleware(), userService.GetUserByID)
		userGroup.GET("/by-phone/:phone", userService.GetUserByPhone)
	}

	deviceGroup := router.Group("/devices")
	{
		deviceGroup.POST("/", interfaces.AuthMiddleware(), deviceService.CreateDevice)
	}

	taxonomyGroup := router.Group("/taxonomies")
	{
		categoryGroup := taxonomyGroup.Group("/categories")
		{
			categoryGroup.GET("/", interfaces.AuthMiddleware(), categoryService.GetAllCategories)
			categoryGroup.GET("/:category_id", interfaces.AuthMiddleware(), categoryService.GetCategoryByID)
		}

		locationGroup := taxonomyGroup.Group("/locations")
		{
			locationGroup.GET("/", locationService.GetAllLocations)
			locationGroup.GET("/:location_id", interfaces.AuthMiddleware(), locationService.GetLocationByID)
			locationGroup.GET("/by-coordinates/:latitude/:longitude/:radius", interfaces.AuthMiddleware(), locationService.GetLocationByCoordinates)
		}

		sizeGroup := taxonomyGroup.Group("/sizes")
		{
			sizeGroup.GET("/", interfaces.AuthMiddleware(), sizeService.GetAllSizes)
			sizeGroup.GET("/:size_id", interfaces.AuthMiddleware(), sizeService.GetSizeByID)
		}

		cityGroup := taxonomyGroup.Group("/cities")
		{
			cityGroup.GET("/", cityService.GetAllCities)
			cityGroup.GET("/:city_id", interfaces.AuthMiddleware(), cityService.GetCityByID)
		}

		shipmentContentGroup := taxonomyGroup.Group("/shipment-contents")
		{
			shipmentContentGroup.GET("/", interfaces.AuthMiddleware(), shipmentContentService.GetAllShipmentContents)
			shipmentContentGroup.GET("/:shipment_content_id", interfaces.AuthMiddleware(), shipmentContentService.GetShipmentContentByID)
		}

		extraServiceGroup := taxonomyGroup.Group("/extra-services")
		{
			extraServiceGroup.GET("/", interfaces.AuthMiddleware(), extraServiceService.GetAllExtraServices)
			extraServiceGroup.GET("/:extra_service_id", interfaces.AuthMiddleware(), extraServiceService.GetExtraServiceByID)
		}

		transportationModeGroup := taxonomyGroup.Group("/transportation-modes")
		{
			transportationModeGroup.GET("/", interfaces.AuthMiddleware(), transportationModeService.GetAllTransportationModes)
			transportationModeGroup.GET("/:transportation_mode_id", interfaces.AuthMiddleware(), transportationModeService.GetTransportationModeByID)
		}

		deliveryTimeGroup := taxonomyGroup.Group("/delivery-times")
		{
			deliveryTimeGroup.GET("/", interfaces.AuthMiddleware(), deliveryTimeService.GetAllDeliveryTimes)
			deliveryTimeGroup.GET("/:delivery_time_id", interfaces.AuthMiddleware(), deliveryTimeService.GetDeliveryTimeByID)
		}

		truckTypeGroup := taxonomyGroup.Group("/truck-types")
		{
			truckTypeGroup.GET("/", interfaces.AuthMiddleware(), truckTypeService.GetAllTruckTypes)
			truckTypeGroup.GET("/:truck_type_id", interfaces.AuthMiddleware(), truckTypeService.GetTruckTypeByID)
		}

		truckModelGroup := taxonomyGroup.Group("/truck-models")
		{
			truckModelGroup.GET("/", interfaces.AuthMiddleware(), truckModelService.GetAllTruckModels)
			truckModelGroup.GET("/:truck_model_id", interfaces.AuthMiddleware(), truckModelService.GetTruckModelByID)
		}
	}

	driverGroup := router.Group("/drivers")
	{
		driverGroup.GET("/", interfaces.AuthMiddleware(), driverService.GetAllDrivers)
		driverGroup.POST("/", interfaces.AuthMiddleware(), driverService.CreateDriver)
		driverGroup.GET("/:driver_id", interfaces.AuthMiddleware(), driverService.GetDriverByID)
		driverGroup.GET("/by-location/:location_id", interfaces.AuthMiddleware(), driverService.GetDriversByUserLocationID)
	}

	orderGroup := router.Group("/orders")
	{
		driverPoolGroup := orderGroup.Group("/driver-pools")
		{
			driverPoolGroup.GET("/by-category/:category_id", interfaces.AuthMiddleware(), orderService.GetAllDriverPoolsByCategoryID)
			driverPoolGroup.PUT("/by-order/:order_id/accept", interfaces.AuthMiddleware(), orderService.AcceptDriverPoolByID)
			driverPoolGroup.PUT("/by-order/:order_id/decline", interfaces.AuthMiddleware(), orderService.DeclineDriverPoolByID)
		}

		orderGroup.GET("/", interfaces.AuthMiddleware(), orderService.GetAllOrders)
		orderGroup.GET("/purchases", interfaces.AuthMiddleware(), orderService.GetAllPurchases)
		orderGroup.POST("/", interfaces.AuthMiddleware(), orderService.CreateOrder)
		orderGroup.GET("/:order_id", interfaces.AuthMiddleware(), orderService.GetOrderByID)
		orderGroup.PUT("/:order_id/cancel", interfaces.AuthMiddleware(), orderService.CancelOrderByID)
		orderGroup.PUT("/:order_id/deliver", interfaces.AuthMiddleware(), orderService.DeliverOrderByID)
		orderGroup.PUT("/:order_id/pickup", interfaces.AuthMiddleware(), orderService.PickupOrderByID)
		orderGroup.POST("/:order_id/rate", interfaces.AuthMiddleware(), orderService.RateOrderByID)
		orderGroup.POST("/:order_id/return", interfaces.AuthMiddleware(), orderService.ReturnOrderByID)
	}

	offerGroup := router.Group("/offers")
	{
		offerGroup.GET("/", interfaces.AuthMiddleware(), offerService.GetAllOffers)
		offerGroup.PUT("/:offer_id/accept", interfaces.AuthMiddleware(), offerService.AcceptOfferByID)
		offerGroup.PUT("/:offer_id/decline", interfaces.AuthMiddleware(), offerService.DeclineOfferByID)
	}

	pageGroup := router.Group("/pages")
	{
		pageGroup.GET("/", pageService.GetAllPages)
		pageGroup.GET("/:page_id", pageService.GetPageByID)
	}

	faqGroup := router.Group("/faqs")
	{
		faqGroup.GET("/", faqService.GetAllFAQs)
		faqGroup.GET("/:faq_id", faqService.GetFAQByID)
	}

	profileGroup := router.Group("/profile")
	{
		profileGroup.GET("/", interfaces.AuthMiddleware(), profileService.GetProfileDetails)
		profileGroup.PUT("/", interfaces.AuthMiddleware(), profileService.UpdateProfileDetails)
		profileGroup.PUT("/phone", interfaces.AuthMiddleware(), profileService.UpdatePhoneNumber)
	}

	userSettingGroup := router.Group("/user-settings")
	{
		userSettingGroup.GET("/", interfaces.AuthMiddleware(), userSettingService.GetUserSettings)
		userSettingGroup.PUT("/", interfaces.AuthMiddleware(), userSettingService.UpdateUserSettings)
	}

	chatGroup := router.Group("/chat")
	{
		chatGroup.GET("/token/client", interfaces.AuthMiddleware(), chatService.GetClientStreamToken)
		chatGroup.GET("/token/driver", interfaces.AuthMiddleware(), chatService.GetDriverStreamToken)
	}

	router.GET("/geo/iso2", geoService.Iso2)

	settingGroup := router.Group("/settings")
	{
		settingGroup.GET("/", settingService.GetAllSettings)
	}

	router.Static("/uploads", "./uploads")

	// Start the router
	router.Run(fmt.Sprintf("%s:%s", Host, Port))
}
