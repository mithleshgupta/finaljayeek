package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// DriverRepository defines the methods for interacting with driver data
type DriverRepository interface {
	CreateDriver(driver *entity.Driver) (*entity.Driver, error)
	DriverWithFieldExists(field string, value string) (bool, error)
	CountDrivers() (int64, error)
	GetAllDrivers(page int, perPage int) ([]entity.Driver, error)
	GetDriverByID(uint64) (*entity.Driver, error)
	GetDriverByUserID(uint64) (*entity.Driver, error)
	CountDriversByUserLocationID(userLocationID uint64) (int64, error)
	GetDriversByUserLocationID(userLocationID uint64, page int, perPage int) ([]entity.Driver, error)
}
