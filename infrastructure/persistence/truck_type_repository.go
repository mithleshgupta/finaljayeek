package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// TruckTypeRepository implements the repository.TruckTypeRepository interface
type TruckTypeRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewTruckTypeRepository creates a new instance of the TruckTypeRepository
func NewTruckTypeRepository(db *gorm.DB) *TruckTypeRepository {
	return &TruckTypeRepository{db: db}
}

func (r *TruckTypeRepository) CountTruckTypes() (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.TruckType{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TruckTypeRepository) GetAllTruckTypes(page int, perPage int) ([]entity.TruckType, error) {
	var truckTypes []entity.TruckType
	if err := r.db.Debug().Model(&entity.TruckType{}).Limit(perPage).Offset((page - 1) * perPage).Find(&truckTypes).Error; err != nil {
		return nil, err
	}
	return truckTypes, nil
}

// GetTruckTypeByID retrieves a truck type by its ID
func (r *TruckTypeRepository) GetTruckTypeByID(id uint64) (*entity.TruckType, error) {
	// TruckType struct to store the retrieved truck type data
	var truckType entity.TruckType
	// Find the truck type by its ID and store the data in the truck type struct
	if err := r.db.Debug().Where("id = ?", id).Take(&truckType).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the truck type data and nil error
	return &truckType, nil
}
