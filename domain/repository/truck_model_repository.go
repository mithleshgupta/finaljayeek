package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// TruckModelRepository defines the methods for interacting with truck model data
type TruckModelRepository interface {
	CountTruckModels() (int64, error)
	GetAllTruckModels(page int, perPage int) ([]entity.TruckModel, error)
	GetTruckModelByID(uint64) (*entity.TruckModel, error)
}
