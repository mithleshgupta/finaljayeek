package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// DeliveryTimeApplication handles the business logic for delivery times
type DeliveryTimeApplication struct {
	deliveryTimeRepo repository.DeliveryTimeRepository
}

var _ DeliveryTimeApplicationInterface = &DeliveryTimeApplication{}

// DeliveryTimeApplicationInterface defines the methods available for DeliveryTimeApplication
type DeliveryTimeApplicationInterface interface {
	CountDeliveryTimes() (int64, error)
	GetAllDeliveryTimes(page int, perPage int) ([]entity.DeliveryTime, error)
	GetDeliveryTimeByID(uint64) (*entity.DeliveryTime, error)
}

func (a *DeliveryTimeApplication) CountDeliveryTimes() (int64, error) {
	return a.deliveryTimeRepo.CountDeliveryTimes()
}

func (a *DeliveryTimeApplication) GetAllDeliveryTimes(page int, perPage int) ([]entity.DeliveryTime, error) {
	return a.deliveryTimeRepo.GetAllDeliveryTimes(page, perPage)
}

// GetByID returns a deliveryTime by its ID
func (a *DeliveryTimeApplication) GetDeliveryTimeByID(deliveryTimeID uint64) (*entity.DeliveryTime, error) {
	return a.deliveryTimeRepo.GetDeliveryTimeByID(deliveryTimeID)
}
