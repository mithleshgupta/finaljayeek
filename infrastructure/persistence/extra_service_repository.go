package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// ExtraServiceRepository implements the repository.ExtraServiceRepository interface
type ExtraServiceRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewExtraServiceRepository creates a new instance of the ExtraServiceRepository
func NewExtraServiceRepository(db *gorm.DB) *ExtraServiceRepository {
	return &ExtraServiceRepository{db: db}
}

func (r *ExtraServiceRepository) CountExtraServices() (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.ExtraService{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ExtraServiceRepository) GetAllExtraServices(page int, perPage int) ([]entity.ExtraService, error) {
	var extraServices []entity.ExtraService
	if err := r.db.Debug().Model(&entity.ExtraService{}).Limit(perPage).Offset((page - 1) * perPage).Find(&extraServices).Error; err != nil {
		return nil, err
	}
	return extraServices, nil
}

// GetExtraServiceByID retrieves a extra service by its ID
func (r *ExtraServiceRepository) GetExtraServiceByID(id uint64) (*entity.ExtraService, error) {
	// ExtraService struct to store the retrieved extra service data
	var extraService entity.ExtraService
	// Find the extra service by its ID and store the data in the extra service struct
	if err := r.db.Debug().Where("id = ?", id).Take(&extraService).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the extra service data and nil error
	return &extraService, nil
}
