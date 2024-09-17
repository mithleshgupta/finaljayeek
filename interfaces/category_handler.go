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

// Categories holds the category-related application interfaces
type Categories struct {
	AuthService  auth.AuthServiceInterface
	TokenService auth.TokenInterface
	CategoryApp  application.CategoryApplicationInterface
}

// NewCategories returns a new instance of Categories
func NewCategories(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, categoryApp application.CategoryApplicationInterface) *Categories {
	return &Categories{
		AuthService:  authService,
		TokenService: tokenService,
		CategoryApp:  categoryApp,
	}
}

// GetAllCategories retrieves a paginated list of all categories.
func (c *Categories) GetAllCategories(ctx *gin.Context) {
	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	// Get the category count from the category application service.
	count, err := c.CategoryApp.CountCategories()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the categories from the category application service.
	categories, err := c.CategoryApp.GetAllCategories(page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(categories) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No categories found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(categories) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var categoryPublicData []interface{}

	for _, category := range categories {
		categoryPublicData = append(categoryPublicData, category.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = categoryPublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the categories as a response.
	response.SendOK(ctx, data, "")
}

// GetCategoryByID retrieves a single category by ID.
func (c *Categories) GetCategoryByID(ctx *gin.Context) {
	// Parse the category ID from the URL parameter.
	categoryID, err := strconv.ParseUint(ctx.Param("category_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid category ID."))
		return
	}

	// Get the category from the category application service.
	category, err := c.CategoryApp.GetCategoryByID(categoryID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Category not found."))
		return
	}

	categoryPublicData := category.PublicData(language.GetLanguage(ctx))

	// Send the category as a response.
	response.SendOK(ctx, categoryPublicData, "")
}
