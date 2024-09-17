package entity

// OrderDriverPool represent an order driver pool
type OrderDriverPool struct {
	OrderID  uint64                `gorm:"index;" validate:"numeric"`
	DriverID uint64                `gorm:"default:null;index;" validate:"numeric"`
	Status   OrderDriverPoolStatus `gorm:"size:255;default:pending;index;" validate:"oneof=pending accepted rejected expired"`
}

type OrderDriverPoolStatus string

const (
	PendingStatus  OrderDriverPoolStatus = "pending"
	AcceptedStatus OrderDriverPoolStatus = "accepted"
	RejectedStatus OrderDriverPoolStatus = "rejected"
	EexpiredStatus OrderDriverPoolStatus = "expired"
)
