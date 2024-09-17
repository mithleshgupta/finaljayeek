package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// SizeRepository defines the methods for interacting with size data
type SizeRepository interface {
	CountSizes() (int64, error)
	GetAllSizes(page int, perPage int) ([]entity.Size, error)
	GetSizeByID(uint64) (*entity.Size, error)
}
