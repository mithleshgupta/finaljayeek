package interfaces

import (
	"strconv"

	"github.com/OmarBader7/web-service-jayeek/application"
	"github.com/OmarBader7/web-service-jayeek/pkg/language"
	"github.com/OmarBader7/web-service-jayeek/pkg/pagination"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
)

// Pages holds the page-related application interfaces
type Pages struct {
	PageApp application.PageApplicationInterface
}

// NewPages returns a new instance of Pages
func NewPages(pageApp application.PageApplicationInterface) *Pages {
	return &Pages{
		PageApp: pageApp,
	}
}

// GetAllPages retrieves a paginated list of all pages.
func (p *Pages) GetAllPages(ctx *gin.Context) {
	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	// Get the page count from the page application service.
	count, err := p.PageApp.CountPages()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the pages from the page application service.
	pages, err := p.PageApp.GetAllPages(page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(pages) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No pages found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(pages) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var pagePublicData []interface{}

	for _, page := range pages {
		pagePublicData = append(pagePublicData, page.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = pagePublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the pages as a response.
	response.SendOK(ctx, data, "")
}

// GetPageByID retrieves a single page by ID.
func (p *Pages) GetPageByID(ctx *gin.Context) {
	// Parse the page ID from the URL parameter.
	pageID, err := strconv.ParseUint(ctx.Param("page_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid page ID."))
		return
	}

	// Get the page from the page application service.
	page, err := p.PageApp.GetPageByID(pageID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	pagePublicData := page.PublicData(language.GetLanguage(ctx))

	// Send the page as a response.
	response.SendOK(ctx, pagePublicData, "")
}
