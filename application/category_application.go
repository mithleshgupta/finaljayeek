package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// CategoryApplication handles the business logic for categories
type CategoryApplication struct {
	categoryRepo repository.CategoryRepository
}

var _ CategoryApplicationInterface = &CategoryApplication{}

// CategoryApplicationInterface defines the methods available for CategoryApplication
type CategoryApplicationInterface interface {
	CountCategories() (int64, error)
	GetAllCategories(page int, perPage int) ([]entity.Category, error)
	GetCategoryByID(uint64) (*entity.Category, error)
}

func (a *CategoryApplication) CountCategories() (int64, error) {
	return a.categoryRepo.CountCategories()
}

func (a *CategoryApplication) GetAllCategories(page int, perPage int) ([]entity.Category, error) {
	return a.categoryRepo.GetAllCategories(page, perPage)
}

// GetByID returns a category by its ID
func (a *CategoryApplication) GetCategoryByID(categoryID uint64) (*entity.Category, error) {
	return a.categoryRepo.GetCategoryByID(categoryID)
}
