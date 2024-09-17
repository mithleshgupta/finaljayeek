package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// BalanceRepository implements the repository.BalanceRepository interface
type BalanceRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewBalanceRepository creates a new instance of the BalanceRepository
func NewBalanceRepository(db *gorm.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}

// CreateBalance creates a new balance in the database
func (r *BalanceRepository) CreateBalance(balance *entity.Balance) (*entity.Balance, error) {
	if err := r.db.Debug().Model(&balance).Create(&balance).Error; err != nil {
		return nil, err
	}
	if err := r.db.Debug().Model(&balance).Preload("Driver").Preload("Driver.User").Preload("Driver.User.Location").Preload("Order").
		Preload("Order.Location").
		Preload("Order.User").
		Preload("Order.User.Location").
		Preload("Order.Driver").
		Preload("Order.Driver.User").
		Preload("Order.Driver.User.Location").
		Preload("Order.Driver.TransportationMode").
		Preload("Order.Category").
		Preload("Order.Size").
		Preload("Order.DeliveryTime").
		Preload("Order.ShipmentContents").
		Preload("Order.ExtraServices").
		Preload("Order.Destination").Take(&balance).Error; err != nil {
		return nil, err
	}

	return balance, nil
}

func (r *BalanceRepository) GetBalanceByDriverID(driverID uint64) (*entity.Balance, error) {
	// Balance struct to store the retrieved driver data
	var balance entity.Balance
	// Find the driver by its ID and store the data in the driver struct
	if err := r.db.Debug().Model(&entity.Balance{}).Preload("Driver").Preload("Driver.User").Preload("Driver.User.Location").Preload("Order").
		Preload("Order.Location").
		Preload("Order.User").
		Preload("Order.User.Location").
		Preload("Order.Driver").
		Preload("Order.Driver.User").
		Preload("Order.Driver.User.Location").
		Preload("Order.Driver.TransportationMode").
		Preload("Order.Category").
		Preload("Order.Size").
		Preload("Order.DeliveryTime").
		Preload("Order.ShipmentContents").
		Preload("Order.ExtraServices").
		Preload("Order.Destination").Where("driver_id = ?", driverID).Take(&balance).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the driver data and nil error
	return &balance, nil
}
