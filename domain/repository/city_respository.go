package repository

import "github.com/OmarBader7/web-service-jayeek/domain/entity"

// CityRepository defines the methods for interacting with city data.
type CityRepository interface {
	CountCities() (int64, error)
	GetAllCities(page int, perPage int) ([]entity.City, error)
	GetCityByID(uint64) (*entity.City, error)
}
