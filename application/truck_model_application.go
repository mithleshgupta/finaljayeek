package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// TruckModelApplication handles the business logic for truck models
type TruckModelApplication struct {
	truckModelRepo repository.TruckModelRepository
}

var _ TruckModelApplicationInterface = &TruckModelApplication{}

// TruckModelApplicationInterface defines the methods available for TruckModelApplication
type TruckModelApplicationInterface interface {
	CountTruckModels() (int64, error)
	GetAllTruckModels(page int, perPage int) ([]entity.TruckModel, error)
	GetTruckModelByID(uint64) (*entity.TruckModel, error)
}

func (a *TruckModelApplication) CountTruckModels() (int64, error) {
	return a.truckModelRepo.CountTruckModels()
}

func (a *TruckModelApplication) GetAllTruckModels(page int, perPage int) ([]entity.TruckModel, error) {
	return a.truckModelRepo.GetAllTruckModels(page, perPage)
}

// GetByID returns a truckModel by its ID
func (a *TruckModelApplication) GetTruckModelByID(truckModelID uint64) (*entity.TruckModel, error) {
	return a.truckModelRepo.GetTruckModelByID(truckModelID)
}
