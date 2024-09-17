package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// ExtraServiceApplication handles the business logic for extra services
type ExtraServiceApplication struct {
	extraServiceRepo repository.ExtraServiceRepository
}

var _ ExtraServiceApplicationInterface = &ExtraServiceApplication{}

// ExtraServiceApplicationInterface defines the methods available for ExtraServiceApplication
type ExtraServiceApplicationInterface interface {
	CountExtraServices() (int64, error)
	GetAllExtraServices(page int, perPage int) ([]entity.ExtraService, error)
	GetExtraServiceByID(uint64) (*entity.ExtraService, error)
}

func (a *ExtraServiceApplication) CountExtraServices() (int64, error) {
	return a.extraServiceRepo.CountExtraServices()
}

func (a *ExtraServiceApplication) GetAllExtraServices(page int, perPage int) ([]entity.ExtraService, error) {
	return a.extraServiceRepo.GetAllExtraServices(page, perPage)
}

// GetByID returns a extraService by its ID
func (a *ExtraServiceApplication) GetExtraServiceByID(extraServiceID uint64) (*entity.ExtraService, error) {
	return a.extraServiceRepo.GetExtraServiceByID(extraServiceID)
}
