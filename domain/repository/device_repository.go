package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// DeviceRepository defines the methods for interacting with category data
type DeviceRepository interface {
	CreateDevice(*entity.Device) (*entity.Device, error)
	GetDeviceByUserID(uint64) (*entity.Device, error)
	DeviceWithFieldExists(field string, value string) (bool, error)
}
