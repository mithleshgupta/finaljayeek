package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// ShipmentContentRepository defines the methods for interacting with shipment content data
type ShipmentContentRepository interface {
	CountShipmentContents() (int64, error)
	GetAllShipmentContents(page int, perPage int) ([]entity.ShipmentContent, error)
	GetShipmentContentByID(uint64) (*entity.ShipmentContent, error)
}
