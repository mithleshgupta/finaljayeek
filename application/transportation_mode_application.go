package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// TransportationModeApplication handles the business logic for transportation modes
type TransportationModeApplication struct {
	transportationModeRepo repository.TransportationModeRepository
}

var _ TransportationModeApplicationInterface = &TransportationModeApplication{}

// TransportationModeApplicationInterface defines the methods available for TransportationModeApplication
type TransportationModeApplicationInterface interface {
	CountTransportationModes() (int64, error)
	GetAllTransportationModes(page int, perPage int) ([]entity.TransportationMode, error)
	GetTransportationModeByID(uint64) (*entity.TransportationMode, error)
}

func (a *TransportationModeApplication) CountTransportationModes() (int64, error) {
	return a.transportationModeRepo.CountTransportationModes()
}

func (a *TransportationModeApplication) GetAllTransportationModes(page int, perPage int) ([]entity.TransportationMode, error) {
	return a.transportationModeRepo.GetAllTransportationModes(page, perPage)
}

// GetByID returns a transportationMode by its ID
func (a *TransportationModeApplication) GetTransportationModeByID(transportationModeID uint64) (*entity.TransportationMode, error) {
	return a.transportationModeRepo.GetTransportationModeByID(transportationModeID)
}
