package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// DeviceApplication handles the business logic for devices
type DeviceApplication struct {
	deviceRepo repository.DeviceRepository
}

var _ DeviceApplicationInterface = &DeviceApplication{}

// DeviceApplicationInterface defines the methods available for DeviceApplication
type DeviceApplicationInterface interface {
	CreateDevice(device *entity.Device) (*entity.Device, error)
	GetDeviceByUserID(uint64) (*entity.Device, error)
	DeviceWithFieldExists(field string, value string) (bool, error)
}

// CreateUser creates a new user in the database
func (a *DeviceApplication) CreateDevice(device *entity.Device) (*entity.Device, error) {
	return a.deviceRepo.CreateDevice(device)
}

// GetByID returns a device by user ID
func (a *DeviceApplication) GetDeviceByUserID(userID uint64) (*entity.Device, error) {
	return a.deviceRepo.GetDeviceByUserID(userID)
}

func (a *DeviceApplication) DeviceWithFieldExists(field string, value string) (bool, error) {
	return a.deviceRepo.DeviceWithFieldExists(field, value)
}
