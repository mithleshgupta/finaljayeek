package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// ExtraServiceRepository defines the methods for interacting with extra service data
type ExtraServiceRepository interface {
	CountExtraServices() (int64, error)
	GetAllExtraServices(page int, perPage int) ([]entity.ExtraService, error)
	GetExtraServiceByID(uint64) (*entity.ExtraService, error)
}
