package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// ShipmentContentApplication handles the business logic for shipment contents
type ShipmentContentApplication struct {
	shipmentContentRepo repository.ShipmentContentRepository
}

var _ ShipmentContentApplicationInterface = &ShipmentContentApplication{}

// ShipmentContentApplicationInterface defines the methods available for ShipmentContentApplication
type ShipmentContentApplicationInterface interface {
	CountShipmentContents() (int64, error)
	GetAllShipmentContents(page int, perPage int) ([]entity.ShipmentContent, error)
	GetShipmentContentByID(uint64) (*entity.ShipmentContent, error)
}

func (a *ShipmentContentApplication) CountShipmentContents() (int64, error) {
	return a.shipmentContentRepo.CountShipmentContents()
}

func (a *ShipmentContentApplication) GetAllShipmentContents(page int, perPage int) ([]entity.ShipmentContent, error) {
	return a.shipmentContentRepo.GetAllShipmentContents(page, perPage)
}

// GetByID returns a shipmentContent by its ID
func (a *ShipmentContentApplication) GetShipmentContentByID(shipmentContentID uint64) (*entity.ShipmentContent, error) {
	return a.shipmentContentRepo.GetShipmentContentByID(shipmentContentID)
}
