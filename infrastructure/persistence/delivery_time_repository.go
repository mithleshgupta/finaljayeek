package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// DeliveryTimeRepository implements the repository.DeliveryTimeRepository interface
type DeliveryTimeRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewDeliveryTimeRepository creates a new instance of the DeliveryTimeRepository
func NewDeliveryTimeRepository(db *gorm.DB) *DeliveryTimeRepository {
	return &DeliveryTimeRepository{db: db}
}

func (r *DeliveryTimeRepository) CountDeliveryTimes() (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.DeliveryTime{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *DeliveryTimeRepository) GetAllDeliveryTimes(page int, perPage int) ([]entity.DeliveryTime, error) {
	var deliveryTimes []entity.DeliveryTime
	if err := r.db.Debug().Model(&entity.DeliveryTime{}).Limit(perPage).Offset((page - 1) * perPage).Find(&deliveryTimes).Error; err != nil {
		return nil, err
	}
	return deliveryTimes, nil
}

// GetDeliveryTimeByID retrieves a delivery time by its ID
func (r *DeliveryTimeRepository) GetDeliveryTimeByID(id uint64) (*entity.DeliveryTime, error) {
	// DeliveryTime struct to store the retrieved delivery time data
	var deliveryTime entity.DeliveryTime
	// Find the delivery time by its ID and store the data in the delivery time struct
	if err := r.db.Debug().Where("id = ?", id).Take(&deliveryTime).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the delivery time data and nil error
	return &deliveryTime, nil
}
