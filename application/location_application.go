package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// LocationApplication handles the business logic for locations
type LocationApplication struct {
	locationRepo repository.LocationRepository
}

var _ LocationApplicationInterface = &LocationApplication{}

// LocationApplicationInterface defines the methods available for LocationApplication
type LocationApplicationInterface interface {
	CountLocations() (int64, error)
	GetAllLocations(page int, perPage int) ([]entity.Location, error)
	GetLocationByID(uint64) (*entity.Location, error)
	GetLocationByCoordinates(float64, float64, *float64) (*entity.Location, error)
}

func (a *LocationApplication) CountLocations() (int64, error) {
	return a.locationRepo.CountLocations()
}

func (a *LocationApplication) GetAllLocations(page int, perPage int) ([]entity.Location, error) {
	return a.locationRepo.GetAllLocations(page, perPage)
}

// GetByID returns a location by its ID
func (a *LocationApplication) GetLocationByID(locationID uint64) (*entity.Location, error) {
	return a.locationRepo.GetLocationByID(locationID)
}

func (a *LocationApplication) GetLocationByCoordinates(longitude, latitude float64, radius *float64) (*entity.Location, error) {
	return a.locationRepo.GetLocationByCoordinates(longitude, latitude, radius)
}
