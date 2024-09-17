package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// DriverApplication handles the business logic for drivers
type DriverApplication struct {
	driverRepo repository.DriverRepository
}

var _ DriverApplicationInterface = &DriverApplication{}

// DriverApplicationInterface defines the methods available for DriverApplication
type DriverApplicationInterface interface {
	CreateDriver(driver *entity.Driver) (*entity.Driver, error)
	DriverWithFieldExists(string, string) (bool, error)
	CountDrivers() (int64, error)
	GetAllDrivers(page int, perPage int) ([]entity.Driver, error)
	GetDriverByID(uint64) (*entity.Driver, error)
	GetDriverByUserID(userID uint64) (*entity.Driver, error)
	CountDriversByUserLocationID(userLocationID uint64) (int64, error)
	GetDriversByUserLocationID(userLocationID uint64, page int, perPage int) ([]entity.Driver, error)
}

// CreateUser creates a new user in the database
func (a *DriverApplication) CreateDriver(driver *entity.Driver) (*entity.Driver, error) {
	return a.driverRepo.CreateDriver(driver)
}

// DriverWithFieldExists checks if a driver with the given field and value exists in the database
func (a *DriverApplication) DriverWithFieldExists(field string, value string) (bool, error) {
	return a.driverRepo.DriverWithFieldExists(field, value)
}

func (a *DriverApplication) CountDrivers() (int64, error) {
	return a.driverRepo.CountDrivers()
}

func (a *DriverApplication) GetAllDrivers(page int, perPage int) ([]entity.Driver, error) {
	return a.driverRepo.GetAllDrivers(page, perPage)
}

// GetByID returns a driver by its ID
func (a *DriverApplication) GetDriverByID(driverID uint64) (*entity.Driver, error) {
	return a.driverRepo.GetDriverByID(driverID)
}

func (a *DriverApplication) GetDriverByUserID(userID uint64) (*entity.Driver, error) {
	return a.driverRepo.GetDriverByUserID(userID)
}

func (a *DriverApplication) CountDriversByUserLocationID(userLocationID uint64) (int64, error) {
	return a.driverRepo.CountDriversByUserLocationID(userLocationID)
}

func (a *DriverApplication) GetDriversByUserLocationID(userLocationID uint64, page int, perPage int) ([]entity.Driver, error) {
	return a.driverRepo.GetDriversByUserLocationID(userLocationID, page, perPage)
}
