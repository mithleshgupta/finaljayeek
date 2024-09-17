package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// TruckTypeRepository defines the methods for interacting with transportation mode data
type TruckTypeRepository interface {
	CountTruckTypes() (int64, error)
	GetAllTruckTypes(page int, perPage int) ([]entity.TruckType, error)
	GetTruckTypeByID(uint64) (*entity.TruckType, error)
}
