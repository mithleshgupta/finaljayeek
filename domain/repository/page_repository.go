package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// PageRepository defines the methods for interacting with page data
type PageRepository interface {
	CountPages() (int64, error)
	GetAllPages(page int, perPage int) ([]entity.Page, error)
	GetPageByID(uint64) (*entity.Page, error)
}
