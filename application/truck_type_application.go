package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// TruckTypeApplication handles the business logic for truck types
type TruckTypeApplication struct {
	truckTypeRepo repository.TruckTypeRepository
}

var _ TruckTypeApplicationInterface = &TruckTypeApplication{}

// TruckTypeApplicationInterface defines the methods available for TruckTypeApplication
type TruckTypeApplicationInterface interface {
	CountTruckTypes() (int64, error)
	GetAllTruckTypes(page int, perPage int) ([]entity.TruckType, error)
	GetTruckTypeByID(uint64) (*entity.TruckType, error)
}

func (a *TruckTypeApplication) CountTruckTypes() (int64, error) {
	return a.truckTypeRepo.CountTruckTypes()
}

func (a *TruckTypeApplication) GetAllTruckTypes(page int, perPage int) ([]entity.TruckType, error) {
	return a.truckTypeRepo.GetAllTruckTypes(page, perPage)
}

// GetByID returns a truckType by its ID
func (a *TruckTypeApplication) GetTruckTypeByID(truckTypeID uint64) (*entity.TruckType, error) {
	return a.truckTypeRepo.GetTruckTypeByID(truckTypeID)
}
