package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// PageApplication handles the business logic for pages
type PageApplication struct {
	pageRepo repository.PageRepository
}

var _ PageApplicationInterface = &PageApplication{}

// PageApplicationInterface defines the methods available for PageApplication
type PageApplicationInterface interface {
	CountPages() (int64, error)
	GetAllPages(page int, perPage int) ([]entity.Page, error)
	GetPageByID(uint64) (*entity.Page, error)
}

func (a *PageApplication) CountPages() (int64, error) {
	return a.pageRepo.CountPages()
}

func (a *PageApplication) GetAllPages(page int, perPage int) ([]entity.Page, error) {
	return a.pageRepo.GetAllPages(page, perPage)
}

// GetByID returns a page by its ID
func (a *PageApplication) GetPageByID(pageID uint64) (*entity.Page, error) {
	return a.pageRepo.GetPageByID(pageID)
}
