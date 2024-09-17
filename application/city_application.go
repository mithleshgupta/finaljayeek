package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// CityApplication handles the business logic for cities.
type CityApplication struct {
	cityRepo repository.CityRepository
}

var _ CityApplicationInterface = &CityApplication{}

// CityApplicationInterface defines the methods available for CityApplication.
type CityApplicationInterface interface {
	CountCities() (int64, error)
	GetAllCities(page int, perPage int) ([]entity.City, error)
	GetCityByID(uint64) (*entity.City, error)
}

func (a *CityApplication) CountCities() (int64, error) {
	return a.cityRepo.CountCities()
}

func (a *CityApplication) GetAllCities(page int, perPage int) ([]entity.City, error) {
	return a.cityRepo.GetAllCities(page, perPage)
}

func (a *CityApplication) GetCityByID(cityID uint64) (*entity.City, error) {
	return a.cityRepo.GetCityByID(cityID)
}
