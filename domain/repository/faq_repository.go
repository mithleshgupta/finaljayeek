package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// FAQRepository defines the methods for interacting with faq data
type FAQRepository interface {
	CountFAQs() (int64, error)
	GetAllFAQs(page int, perPage int) ([]entity.FAQ, error)
	GetFAQByID(uint64) (*entity.FAQ, error)
}
