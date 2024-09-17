package entity

import "time"

// Offer represent an order
type Offer struct {
	ID        uint64      `gorm:"primary_key;auto_increment" json:"id"`
	DriverID  uint64      `gorm:"index" json:"driver_id"`
	OrderID   uint64      `gorm:"index" json:"order_id"`
	Amount    float64     `json:"amount" validate:"required,numeric"`
	Status    OfferStatus `gorm:"size:255;default:pending;index;" validate:"oneof=pending accepted declined"`
	CreatedAt time.Time   `gorm:"default:CURRENT_TIMESTAMP;index;" json:"created_at"`
	UpdatedAt time.Time   `gorm:"default:null" json:"updated_at"`
	Driver    Driver      `gorm:"foreignkey:DriverID" json:"driver"`
	Order     Order       `gorm:"foreignkey:OrderID" json:"order"`
}

type OfferPublicData struct {
	ID       uint64            `json:"id"`
	DriverID uint64            `json:"driver_id"`
	OrderID  uint64            `json:"order_id"`
	Amount   float64           `json:"amount"`
	Status   OfferStatus       `json:"status"`
	Driver   *DriverPublicData `json:"driver"`
	Order    *OrderPublicData  `json:"order"`
}

type OfferStatus string

const (
	OfferStatusPending  OfferStatus = "pending"
	OfferStatusAccepted OfferStatus = "accepted"
	OfferStatusDeclined OfferStatus = "declined"
)

// PublicData returns a copy of the offer's public information
func (o *Offer) PublicData(languageCode string) interface{} {
	driverPublicData := o.Driver.PublicData(languageCode).(*DriverPublicData)
	orderPublicData := o.Order.PublicData(languageCode).(*OrderPublicData)

	return &OfferPublicData{
		ID:       o.ID,
		DriverID: o.DriverID,
		OrderID:  o.OrderID,
		Amount:   o.Amount,
		Status:   o.Status,
		Driver:   driverPublicData,
		Order:    orderPublicData,
	}
}
