package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// TruckModelRepository implements the repository.TruckModelRepository interface
type TruckModelRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewTruckModelRepository creates a new instance of the TruckModelRepository
func NewTruckModelRepository(db *gorm.DB) *TruckModelRepository {
	return &TruckModelRepository{db: db}
}

func (r *TruckModelRepository) CountTruckModels() (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.TruckModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TruckModelRepository) GetAllTruckModels(page int, perPage int) ([]entity.TruckModel, error) {
	var truckModels []entity.TruckModel
	if err := r.db.Debug().Model(&entity.TruckModel{}).Limit(perPage).Offset((page - 1) * perPage).Find(&truckModels).Error; err != nil {
		return nil, err
	}
	return truckModels, nil
}

// GetTruckModelByID retrieves a truck model by its ID
func (r *TruckModelRepository) GetTruckModelByID(id uint64) (*entity.TruckModel, error) {
	// TruckModel struct to store the retrieved truck model data
	var truckModel entity.TruckModel
	// Find the truck model by its ID and store the data in the truck model struct
	if err := r.db.Debug().Where("id = ?", id).Take(&truckModel).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the truck model data and nil error
	return &truckModel, nil
}
