package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// BalanceApplication handles the business logic for balances
type BalanceApplication struct {
	balanceRepo repository.BalanceRepository
}

var _ BalanceApplicationInterface = &BalanceApplication{}

// BalanceApplicationInterface defines the methods available for BalanceApplication
type BalanceApplicationInterface interface {
	CreateBalance(balance *entity.Balance) (*entity.Balance, error)
	GetBalanceByDriverID(uint64) (*entity.Balance, error)
}

// CreateUser creates a new user in the database
func (a *BalanceApplication) CreateBalance(balance *entity.Balance) (*entity.Balance, error) {
	return a.balanceRepo.CreateBalance(balance)
}

// GetByID returns a balance by its ID
func (a *BalanceApplication) GetBalanceByDriverID(driverID uint64) (*entity.Balance, error) {
	return a.balanceRepo.GetBalanceByDriverID(driverID)
}
