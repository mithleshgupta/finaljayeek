package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// FAQRepository implements the repository.FAQRepository interface
type FAQRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewFAQRepository creates a new instance of the FAQRepository
func NewFAQRepository(db *gorm.DB) *FAQRepository {
	return &FAQRepository{db: db}
}

func (r *FAQRepository) CountFAQs() (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.FAQ{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *FAQRepository) GetAllFAQs(page int, perPage int) ([]entity.FAQ, error) {
	var faqs []entity.FAQ
	if err := r.db.Debug().Model(&entity.FAQ{}).Limit(perPage).Offset((page - 1) * perPage).Find(&faqs).Error; err != nil {
		return nil, err
	}
	return faqs, nil
}

// GetFAQByID retrieves a faq by its ID
func (r *FAQRepository) GetFAQByID(id uint64) (*entity.FAQ, error) {
	// FAQ struct to store the retrieved faq data
	var faq entity.FAQ
	// Find the faq by its ID and store the data in the faq struct
	if err := r.db.Debug().Where("id = ?", id).Take(&faq).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the faq data and nil error
	return &faq, nil
}
