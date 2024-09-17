package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// SizeApplication handles the business logic for sizes
type SizeApplication struct {
	sizeRepo repository.SizeRepository
}

var _ SizeApplicationInterface = &SizeApplication{}

// SizeApplicationInterface defines the methods available for SizeApplication
type SizeApplicationInterface interface {
	CountSizes() (int64, error)
	GetAllSizes(page int, perPage int) ([]entity.Size, error)
	GetSizeByID(uint64) (*entity.Size, error)
}

func (a *SizeApplication) CountSizes() (int64, error) {
	return a.sizeRepo.CountSizes()
}

func (a *SizeApplication) GetAllSizes(page int, perPage int) ([]entity.Size, error) {
	return a.sizeRepo.GetAllSizes(page, perPage)
}

// GetByID returns a size by its ID
func (a *SizeApplication) GetSizeByID(sizeID uint64) (*entity.Size, error) {
	return a.sizeRepo.GetSizeByID(sizeID)
}
