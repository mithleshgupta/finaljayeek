package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// TransportationModeRepository defines the methods for interacting with transportation mode data
type TransportationModeRepository interface {
	CountTransportationModes() (int64, error)
	GetAllTransportationModes(page int, perPage int) ([]entity.TransportationMode, error)
	GetTransportationModeByID(uint64) (*entity.TransportationMode, error)
}
