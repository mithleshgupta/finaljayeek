package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// BalanceRepository defines the methods for interacting with category data
type BalanceRepository interface {
	CreateBalance(*entity.Balance) (*entity.Balance, error)
	GetBalanceByDriverID(uint64) (*entity.Balance, error)
}
