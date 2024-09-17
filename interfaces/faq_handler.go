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

// FAQs holds the faq-related application interfaces
type FAQs struct {
	FAQApp application.FAQApplicationInterface
}

// NewFAQs returns a new instance of FAQs
func NewFAQs(faqApp application.FAQApplicationInterface) *FAQs {
	return &FAQs{
		FAQApp: faqApp,
	}
}

// GetAllFAQs retrieves a paginated list of all faqs.
func (p *FAQs) GetAllFAQs(ctx *gin.Context) {
	// Get the desired faq number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per faq.
	perPage := 30

	// Get the faq count from the faq application service.
	count, err := p.FAQApp.CountFAQs()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the faqs from the faq application service.
	faqs, err := p.FAQApp.GetAllFAQs(page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(faqs) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No pages found."))
		return
	}

	// Check if the faq is valid
	if page <= 0 || (len(faqs) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("FAQ not found."))
		return
	}

	var faqPublicData []interface{}

	for _, faq := range faqs {
		faqPublicData = append(faqPublicData, faq.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = faqPublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the faqs as a response.
	response.SendOK(ctx, data, "")
}

// GetFAQByID retrieves a single faq by ID.
func (p *FAQs) GetFAQByID(ctx *gin.Context) {
	// Parse the faq ID from the URL parameter.
	faqID, err := strconv.ParseUint(ctx.Param("faq_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid faq ID."))
		return
	}

	// Get the faq from the faq application service.
	faq, err := p.FAQApp.GetFAQByID(faqID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("FAQ not found."))
		return
	}

	faqPublicData := faq.PublicData(language.GetLanguage(ctx))

	// Send the faq as a response.
	response.SendOK(ctx, faqPublicData, "")
}
