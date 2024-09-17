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

// Sizes holds the size-related application interfaces
type Sizes struct {
	AuthService  auth.AuthServiceInterface
	TokenService auth.TokenInterface
	SizeApp      application.SizeApplicationInterface
}

// NewSizes returns a new instance of Sizes
func NewSizes(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, sizeApp application.SizeApplicationInterface) *Sizes {
	return &Sizes{
		AuthService:  authService,
		TokenService: tokenService,
		SizeApp:      sizeApp,
	}
}

// GetAllSizes retrieves a paginated list of all sizes.
func (s *Sizes) GetAllSizes(ctx *gin.Context) {
	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	// Get the size count from the size application service.
	count, err := s.SizeApp.CountSizes()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the sizes from the size application service.
	sizes, err := s.SizeApp.GetAllSizes(page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(sizes) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No sizes found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(sizes) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var sizePublicData []interface{}

	for _, size := range sizes {
		sizePublicData = append(sizePublicData, size.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = sizePublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the sizes as a response.
	response.SendOK(ctx, data, "")
}

// GetSizeByID retrieves a single size by ID.
func (s *Sizes) GetSizeByID(ctx *gin.Context) {
	// Parse the size ID from the URL parameter.
	sizeID, err := strconv.ParseUint(ctx.Param("size_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid size ID."))
		return
	}

	// Get the size from the size application service.
	size, err := s.SizeApp.GetSizeByID(sizeID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Size not found."))
		return
	}

	sizePublicData := size.PublicData(language.GetLanguage(ctx))

	// Send the size as a response.
	response.SendOK(ctx, sizePublicData, "")
}
