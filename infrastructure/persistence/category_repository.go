package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// CategoryRepository implements the repository.CategoryRepository interface
type CategoryRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewCategoryRepository creates a new instance of the CategoryRepository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CountCategories() (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.Category{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CategoryRepository) GetAllCategories(page int, perPage int) ([]entity.Category, error) {
	var categories []entity.Category
	if err := r.db.Debug().Model(&entity.Category{}).Limit(perPage).Offset((page - 1) * perPage).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryByID retrieves a category by its ID
func (r *CategoryRepository) GetCategoryByID(id uint64) (*entity.Category, error) {
	// Category struct to store the retrieved category data
	var category entity.Category
	// Find the category by its ID and store the data in the category struct
	if err := r.db.Debug().Where("id = ?", id).Take(&category).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the category data and nil error
	return &category, nil
}
