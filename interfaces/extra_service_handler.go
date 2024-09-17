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

// ExtraServices holds the extra service-related application interfaces
type ExtraServices struct {
	AuthService     auth.AuthServiceInterface
	TokenService    auth.TokenInterface
	ExtraServiceApp application.ExtraServiceApplicationInterface
}

// NewExtraServices returns a new instance of ExtraServices
func NewExtraServices(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, extraServiceApp application.ExtraServiceApplicationInterface) *ExtraServices {
	return &ExtraServices{
		AuthService:     authService,
		TokenService:    tokenService,
		ExtraServiceApp: extraServiceApp,
	}
}

// GetAllExtraServices retrieves a paginated list of all extra services.
func (s *ExtraServices) GetAllExtraServices(ctx *gin.Context) {
	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	// Get the extra service count from the extra service application service.
	count, err := s.ExtraServiceApp.CountExtraServices()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the extraServices from the extra service application service.
	extraServices, err := s.ExtraServiceApp.GetAllExtraServices(page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(extraServices) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No extra services found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(extraServices) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var extraServicePublicData []interface{}

	for _, extraService := range extraServices {
		extraServicePublicData = append(extraServicePublicData, extraService.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = extraServicePublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the extraServices as a response.
	response.SendOK(ctx, data, "")
}

// GetExtraServiceByID retrieves a single extra service by ID.
func (s *ExtraServices) GetExtraServiceByID(ctx *gin.Context) {
	// Parse the extra service ID from the URL parameter.
	extraServiceID, err := strconv.ParseUint(ctx.Param("extra_service_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid extra service ID."))
		return
	}

	// Get the extra service from the extra service application service.
	extraService, err := s.ExtraServiceApp.GetExtraServiceByID(extraServiceID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Extra service not found."))
		return
	}

	extraServicePublicData := extraService.PublicData(language.GetLanguage(ctx))

	// Send the extra service as a response.
	response.SendOK(ctx, extraServicePublicData, "")
}
