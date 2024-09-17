package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// TransportationModeRepository implements the repository.TransportationModeRepository interface
type TransportationModeRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewTransportationModeRepository creates a new instance of the TransportationModeRepository
func NewTransportationModeRepository(db *gorm.DB) *TransportationModeRepository {
	return &TransportationModeRepository{db: db}
}

func (r *TransportationModeRepository) CountTransportationModes() (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.TransportationMode{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TransportationModeRepository) GetAllTransportationModes(page int, perPage int) ([]entity.TransportationMode, error) {
	var transportationModes []entity.TransportationMode
	if err := r.db.Debug().Model(&entity.TransportationMode{}).Limit(perPage).Offset((page - 1) * perPage).Find(&transportationModes).Error; err != nil {
		return nil, err
	}
	return transportationModes, nil
}

// GetTransportationModeByID retrieves a transportation mode by its ID
func (r *TransportationModeRepository) GetTransportationModeByID(id uint64) (*entity.TransportationMode, error) {
	// TransportationMode struct to store the retrieved transportation mode data
	var transportationMode entity.TransportationMode
	// Find the transportation mode by its ID and store the data in the transportation mode struct
	if err := r.db.Debug().Where("id = ?", id).Take(&transportationMode).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the transportation mode data and nil error
	return &transportationMode, nil
}
