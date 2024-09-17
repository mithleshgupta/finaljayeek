package entity

// Order represent an order extra service
type OrderExtraService struct {
	OrderID        uint64 `gorm:"primary_key" json:"order_id"`
	ExtraServiceID uint64 `gorm:"primary_key" json:"extra_service_id"`
}
