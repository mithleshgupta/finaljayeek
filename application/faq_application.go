package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// FAQApplication handles the business logic for faqs
type FAQApplication struct {
	faqRepo repository.FAQRepository
}

var _ FAQApplicationInterface = &FAQApplication{}

// FAQApplicationInterface defines the methods available for FAQApplication
type FAQApplicationInterface interface {
	CountFAQs() (int64, error)
	GetAllFAQs(page int, perPage int) ([]entity.FAQ, error)
	GetFAQByID(uint64) (*entity.FAQ, error)
}

func (a *FAQApplication) CountFAQs() (int64, error) {
	return a.faqRepo.CountFAQs()
}

func (a *FAQApplication) GetAllFAQs(page int, perPage int) ([]entity.FAQ, error) {
	return a.faqRepo.GetAllFAQs(page, perPage)
}

// GetByID returns a faq by its ID
func (a *FAQApplication) GetFAQByID(faqID uint64) (*entity.FAQ, error) {
	return a.faqRepo.GetFAQByID(faqID)
}
