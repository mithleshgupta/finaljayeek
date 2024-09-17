package entity

import "time"

// Balance represent an order
type Balance struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	DriverID  uint64    `gorm:"default:null;index;" json:"driver_id" validate:"numeric"`
	OrderID   uint64    `gorm:"default:null;index;" json:"order_id" validate:"numeric"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:null" json:"updated_at"`
	Balance   float64   `json:"balance" validate:"required,numeric"`
	Driver    Driver    `gorm:"foreignKey:DriverID" json:"driver"`
	Order     Order     `gorm:"foreignKey:OrderID" json:"order"`
}
