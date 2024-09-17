package entity

// Order represent an order shipment content
type OrderShipmentContent struct {
	OrderID           uint64 `gorm:"primary_key" json:"order_id"`
	ShipmentContentID uint64 `gorm:"primary_key" json:"shipment_content_id"`
}
