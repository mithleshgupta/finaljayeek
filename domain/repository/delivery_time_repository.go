package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// DeliveryTimeRepository defines the methods for interacting with delivery time data
type DeliveryTimeRepository interface {
	CountDeliveryTimes() (int64, error)
	GetAllDeliveryTimes(page int, perPage int) ([]entity.DeliveryTime, error)
	GetDeliveryTimeByID(uint64) (*entity.DeliveryTime, error)
}
