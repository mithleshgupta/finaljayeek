package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// PageRepository implements the repository.PageRepository interface
type PageRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewPageRepository creates a new instance of the PageRepository
func NewPageRepository(db *gorm.DB) *PageRepository {
	return &PageRepository{db: db}
}

func (r *PageRepository) CountPages() (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.Page{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PageRepository) GetAllPages(page int, perPage int) ([]entity.Page, error) {
	var pages []entity.Page
	if err := r.db.Debug().Model(&entity.Page{}).Limit(perPage).Offset((page - 1) * perPage).Find(&pages).Error; err != nil {
		return nil, err
	}
	return pages, nil
}

// GetPageByID retrieves a page by its ID
func (r *PageRepository) GetPageByID(id uint64) (*entity.Page, error) {
	// Page struct to store the retrieved page data
	var page entity.Page
	// Find the page by its ID and store the data in the page struct
	if err := r.db.Debug().Where("id = ?", id).Take(&page).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the page data and nil error
	return &page, nil
}
