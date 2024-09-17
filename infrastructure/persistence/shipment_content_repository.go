package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// ShipmentContentRepository implements the repository.ShipmentContentRepository interface
type ShipmentContentRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewShipmentContentRepository creates a new instance of the ShipmentContentRepository
func NewShipmentContentRepository(db *gorm.DB) *ShipmentContentRepository {
	return &ShipmentContentRepository{db: db}
}

func (r *ShipmentContentRepository) CountShipmentContents() (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.ShipmentContent{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ShipmentContentRepository) GetAllShipmentContents(page int, perPage int) ([]entity.ShipmentContent, error) {
	var shipmentContents []entity.ShipmentContent
	if err := r.db.Debug().Model(&entity.ShipmentContent{}).Limit(perPage).Offset((page - 1) * perPage).Find(&shipmentContents).Error; err != nil {
		return nil, err
	}
	return shipmentContents, nil
}

// GetShipmentContentByID retrieves a shipment content by its ID
func (r *ShipmentContentRepository) GetShipmentContentByID(id uint64) (*entity.ShipmentContent, error) {
	// ShipmentContent struct to store the retrieved shipment content data
	var shipmentContent entity.ShipmentContent
	// Find the shipment content by its ID and store the data in the shipment content struct
	if err := r.db.Debug().Where("id = ?", id).Take(&shipmentContent).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the shipment content data and nil error
	return &shipmentContent, nil
}
