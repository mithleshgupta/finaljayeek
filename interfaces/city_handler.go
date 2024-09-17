package interfaces

import (
	"strconv"

	"github.com/OmarBader7/web-service-jayeek/application"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/auth"
	"github.com/OmarBader7/web-service-jayeek/pkg/pagination"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	"github.com/gin-gonic/gin"
)

// Cities holds the city-related application interfaces.
type Cities struct {
	AuthService  auth.AuthServiceInterface
	TokenService auth.TokenInterface
	CityApp      application.CityApplicationInterface
}

// NewCities returns a new instance of Cities.
func NewCities(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, cityApp application.CityApplicationInterface) *Cities {
	return &Cities{
		AuthService:  authService,
		TokenService: tokenService,
		CityApp:      cityApp,
	}
}

// GetAllCities retrieves a paginated list of all cities.
func (c *Cities) GetAllCities(ctx *gin.Context) {
	page := pagination.GetPage(ctx)
	perPage := 30
		
	count, err := c.CityApp.CountCities()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	cities, err := c.CityApp.GetAllCities(page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(cities) == 0 {
		response.SendOK(ctx, nil, "No cities found.")
		return
	}

	if page <= 0 || (len(cities) == 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, "Page not found.")
		return
	}

	var cityPublicData []interface{}
	for _, city := range cities {
		cityPublicData = append(cityPublicData, city.PublicData())
	}

	data := make(map[string]interface{})
	data["data"] = cityPublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	response.SendOK(ctx, data, "")
}

// GetCityByID retrieves a single city by ID.
func (c *Cities) GetCityByID(ctx *gin.Context) {
	cityID, err := strconv.ParseUint(ctx.Param("city_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, "Invalid city ID.")
		return
	}

	city, err := c.CityApp.GetCityByID(cityID)
	if err != nil {
		response.SendNotFound(ctx, "City not found.")
		return
	}

	cityPublicData := city.PublicData()
	response.SendOK(ctx, cityPublicData, "")
}
