package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// SizeRepository implements the repository.SizeRepository interface
type SizeRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewSizeRepository creates a new instance of the SizeRepository
func NewSizeRepository(db *gorm.DB) *SizeRepository {
	return &SizeRepository{db: db}
}

func (r *SizeRepository) CountSizes() (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.Size{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *SizeRepository) GetAllSizes(page int, perPage int) ([]entity.Size, error) {
	var sizes []entity.Size
	if err := r.db.Debug().Model(&entity.Size{}).Limit(perPage).Offset((page - 1) * perPage).Find(&sizes).Error; err != nil {
		return nil, err
	}
	return sizes, nil
}

// GetSizeByID retrieves a size by its ID
func (r *SizeRepository) GetSizeByID(id uint64) (*entity.Size, error) {
	// Size struct to store the retrieved size data
	var size entity.Size
	// Find the size by its ID and store the data in the size struct
	if err := r.db.Debug().Where("id = ?", id).Take(&size).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the size data and nil error
	return &size, nil
}
