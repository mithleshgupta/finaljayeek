package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// LocationRepository defines the methods for interacting with location data
type LocationRepository interface {
	CountLocations() (int64, error)
	GetAllLocations(page int, perPage int) ([]entity.Location, error)
	GetLocationByID(uint64) (*entity.Location, error)
	GetLocationByCoordinates(float64, float64, *float64) (*entity.Location, error)
}
