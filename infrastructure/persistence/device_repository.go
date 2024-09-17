package persistence

import (
	"fmt"

	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// DeviceRepository implements the repository.DeviceRepository interface
type DeviceRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewDeviceRepository creates a new instance of the DeviceRepository
func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// CreateDevice creates a new device in the database
func (r *DeviceRepository) CreateDevice(device *entity.Device) (*entity.Device, error) {
	if err := r.db.Debug().Model(&device).Create(&device).Error; err != nil {
		return nil, err
	}
	if err := r.db.Debug().Model(&device).Preload("User").Preload("User.Location").Take(&device).Error; err != nil {
		return nil, err
	}

	return device, nil
}

func (r *DeviceRepository) GetDeviceByUserID(userID uint64) (*entity.Device, error) {
	// Device struct to store the retrieved user data
	var device entity.Device
	// Find the user by its ID and store the data in the user struct
	if err := r.db.Debug().Model(&entity.Device{}).Preload("User").Preload("User.Location").Where("user_id = ?", userID).Take(&device).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the user data and nil error
	return &device, nil
}

// DeviceWithFieldExists checks if a user with the given field and value exists in the database
func (r *DeviceRepository) DeviceWithFieldExists(field string, value string) (bool, error) {
	var device entity.Device
	var devicesCount int64
	if err := r.db.Debug().Model(&device).Select("id").Where(fmt.Sprintf("%s = ?", field), value).Count(&devicesCount).Error; err != nil {
		return false, err
	}
	return devicesCount > 0, nil
}
