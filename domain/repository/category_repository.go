package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// CategoryRepository defines the methods for interacting with category data
type CategoryRepository interface {
	CountCategories() (int64, error)
	GetAllCategories(page int, perPage int) ([]entity.Category, error)
	GetCategoryByID(uint64) (*entity.Category, error)
}
