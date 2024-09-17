package persistence

import (
	"fmt"

	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// DriverRepository implements the repository.DriverRepository interface
type DriverRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewDriverRepository creates a new instance of the DriverRepository
func NewDriverRepository(db *gorm.DB) *DriverRepository {
	return &DriverRepository{db: db}
}

// DriverWithFieldExists checks if a driver with the given field and value exists in the database
func (r *DriverRepository) DriverWithFieldExists(field string, value string) (bool, error) {
	var driver entity.Driver
	var driversCount int64
	if err := r.db.Debug().Model(&driver).Preload("User").Preload("User.Location").Preload("TransportationMode").Select("id").Where(fmt.Sprintf("%s = ?", field), value).Count(&driversCount).Error; err != nil {
		return false, err
	}
	return driversCount > 0, nil
}

// CreateDriver creates a new driver in the database
func (r *DriverRepository) CreateDriver(driver *entity.Driver) (*entity.Driver, error) {
	if err := r.db.Debug().Model(&driver).Create(&driver).Error; err != nil {
		return nil, err
	}
	if err := r.db.Debug().Model(&driver).Preload("User").Preload("User.Location").Preload("TransportationMode").Take(&driver).Error; err != nil {
		return nil, err
	}

	var orders []entity.Driver
	if err := r.db.Debug().Table("orders").Select("orders.*").
		Where("orders.status = ?", entity.OrderCreatedStatus).
		Where("ST_DWithin(ST_MakePoint(orders.longitude, orders.latitude)::geography, ST_MakePoint(?, ?)::geography, ?)", driver.Longitude, driver.Latitude, 50000).
		Order(fmt.Sprintf("ST_Distance(ST_MakePoint(orders.longitude, orders.latitude)::geography, ST_MakePoint(%.6f, %.6f)::geography)", driver.Longitude, driver.Latitude)).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	for _, order := range orders {
		if isAvailableSetting, err := driver.User.GetSettingByKey("is_available"); err == nil {
			isAvailable, ok := isAvailableSetting.(bool)
			if !ok {
				isAvailable = false // Default value if the setting is not a boolean
			}

			if isAvailable {
				var orderDriverPool entity.OrderDriverPool

				orderDriverPool.DriverID = driver.ID
				orderDriverPool.OrderID = order.ID

				if err := r.db.Debug().Model(&orderDriverPool).Create(&orderDriverPool).Error; err != nil {
					return nil, err
				}
			}
		}
	}

	return driver, nil
}

func (r *DriverRepository) CountDrivers() (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.Driver{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *DriverRepository) GetAllDrivers(page int, perPage int) ([]entity.Driver, error) {
	var drivers []entity.Driver
	if err := r.db.Debug().Model(&entity.Driver{}).Preload("User").Preload("User.Location").Preload("TransportationMode").Limit(perPage).Offset((page - 1) * perPage).Find(&drivers).Error; err != nil {
		return nil, err
	}
	return drivers, nil
}

// GetDriverByID retrieves a driver by its ID
func (r *DriverRepository) GetDriverByID(id uint64) (*entity.Driver, error) {
	// Driver struct to store the retrieved driver data
	var driver entity.Driver
	// Find the driver by its ID and store the data in the driver struct
	if err := r.db.Debug().Model(&entity.Driver{}).Preload("User").Preload("User.Location").Preload("TransportationMode").Where("id = ?", id).Take(&driver).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the driver data and nil error
	return &driver, nil
}

func (r *DriverRepository) GetDriverByUserID(userID uint64) (*entity.Driver, error) {
	// Driver struct to store the retrieved driver data
	var driver entity.Driver
	// Find the driver by its ID and store the data in the driver struct
	if err := r.db.Debug().Model(&entity.Driver{}).Preload("User").Preload("User.Location").Preload("TransportationMode").Where("user_id = ?", userID).Take(&driver).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the driver data and nil error
	return &driver, nil
}

func (r *DriverRepository) CountDriversByUserLocationID(userLocationID uint64) (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.Driver{}).Joins("JOIN users ON drivers.user_id = users.id").
		Where("users.location_id = ?", userLocationID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *DriverRepository) GetDriversByUserLocationID(userLocationID uint64, page int, perPage int) ([]entity.Driver, error) {
	var drivers []entity.Driver
	if err := r.db.Debug().Model(&entity.Driver{}).Preload("User").Preload("User.Location").Preload("TransportationMode").Joins("JOIN users ON drivers.user_id = users.id").
		Where("users.location_id = ?", userLocationID).Limit(perPage).Offset((page - 1) * perPage).Find(&drivers).Error; err != nil {
		return nil, err
	}
	return drivers, nil
}
