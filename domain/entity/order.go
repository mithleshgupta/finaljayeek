package entity

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Order represent an order
type Order struct {
	ID                   uint64              `gorm:"primary_key;auto_increment" json:"id"`
	LocationID           uint64              `gorm:"index;" json:"location_id" validate:"required,numeric"`
	UserID               uint64              `gorm:"index;" json:"user_id" validate:"numeric"`
	DriverID             uint64              `gorm:"default:null;index;" json:"driver_id" validate:"numeric"`
	RecipientID          uint64              `gorm:"default:null;index;" json:"recipient_id" validate:"numeric"`
	CategoryID           uint64              `gorm:"index;" json:"category_id" validate:"required,numeric"`
	SizeID               uint64              `gorm:"default:null;index;" json:"size_id" validate:"omitempty,numeric"`
	TruckTypeID          uint64              `gorm:"default:null;index;" json:"truck_type_id" validate:"omitempty,numeric"`
	TruckModelID         uint64              `gorm:"default:null;index;" json:"truck_model_id" validate:"omitempty,numeric"`
	DeliveryTimeID       uint64              `gorm:"index;" json:"delivery_time_id" validate:"required,numeric"`
	ShipmentContentIDs   []uint64            `gorm:"-" json:"shipment_content_ids" validate:"required"`
	ExtraServiceIDs      *[]uint64           `gorm:"-" json:"extra_service_ids"`
	DestinationID        uint64              `gorm:"index;" json:"destination_id" validate:"required,numeric"`
	Quantity             uint64              `gorm:"default:1" json:"quantity" validate:"required,numeric"`
	RecipientPhoneNumber string              `gorm:"type:varchar(255)" json:"recipient_phone_number" validate:"required,e164"`
	Notes                *string             `gorm:"type:varchar(255);default:null" json:"notes"`
	Amount               *float64            `gorm:"default:null" json:"amount" validate:"omitempty,numeric"`
	Latitude             float64             `gorm:"type:decimal(10,8);not null;index;" json:"latitude" validate:"required"`
	Longitude            float64             `gorm:"type:decimal(11,8);not null;index;" json:"longitude" validate:"required"`
	PaymentMethod        *OrderPaymentMethod `gorm:"size:255;default:null" json:"payment_method" validate:"omitempty,oneof=cash"`
	Status               OrderStatus         `gorm:"size:255;default:order_created;index;" json:"status" validate:"oneof=order_created order_accepted pickup_in_progress shipment_picked_up in_transit at_destination_city out_for_delivery delivery_attempted delivery_rescheduled shipment_delivered order_completed order_canceled shipment_returned"`
	IsSender             bool                `gorm:"-" json:"is_sender"`
	IsReceiver           bool                `gorm:"-" json:"is_receiver"`
	CreatedAt            time.Time           `gorm:"default:CURRENT_TIMESTAMP;index;" json:"created_at"`
	UpdatedAt            time.Time           `gorm:"default:null" json:"updated_at"`
	Location             Location            `gorm:"foreignKey:LocationID" json:"location"`
	User                 User                `gorm:"foreignKey:UserID" json:"user"`
	Driver               Driver              `gorm:"foreignKey:DriverID" json:"driver"`
	Recipient            User                `gorm:"foreignKey:RecipientID" json:"recipient"`
	Category             Category            `gorm:"foreignKey:CategoryID" json:"category"`
	Size                 Size                `gorm:"foreignKey:SizeID" json:"size"`
	TruckType            TruckType           `gorm:"foreignKey:TruckTypeID" json:"truck_type"`
	TruckModel           TruckModel          `gorm:"foreignKey:TruckModelID" json:"truck_model"`
	DeliveryTime         DeliveryTime        `gorm:"foreignKey:DeliveryTimeID" json:"delivery_time"`
	ShipmentContents     []ShipmentContent   `gorm:"many2many:order_shipment_contents" json:"shipment_contents"`
	ExtraServices        *[]ExtraService     `gorm:"many2many:order_extra_services" json:"extra_services"`
	Destination          Location            `gorm:"foreignKey:DestinationID" json:"destination"`
	Rating               *Rating             `gorm:"-" json:"rating"`
}

type OrderPublicData struct {
	ID                 uint64                       `json:"id"`
	LocationID         uint64                       `json:"location_id"`
	UserID             uint64                       `json:"user_id"`
	DriverID           uint64                       `json:"driver_id"`
	RecipientID        uint64                       `json:"recipient_id"`
	CategoryID         uint64                       `json:"category_id"`
	SizeID             uint64                       `json:"size_id"`
	TruckTypeID        uint64                       `json:"truck_type_id"`
	TruckModelID       uint64                       `json:"truck_model_id"`
	DeliveryTimeID     uint64                       `json:"delivery_time_id"`
	ShipmentContentIDs []uint64                     `json:"shipment_content_ids"`
	ExtraServiceIDs    *[]uint64                    `json:"extra_service_ids"`
	DestinationID      uint64                       `json:"destination_id"`
	Quantity           uint64                       `json:"quantity"`
	Notes              *string                      `json:"notes"`
	Amount             *float64                     `json:"amount"`
	Latitude           float64                      `json:"latitude"`
	Longitude          float64                      `json:"longitude"`
	CreatedAt          time.Time                    `json:"created_at"`
	Status             OrderStatus                  `json:"status"`
	PaymentMethod      *OrderPaymentMethod          `json:"payment_method"`
	IsSender           bool                         `json:"is_sender"`
	IsReceiver         bool                         `json:"is_receiver"`
	Location           *LocationPublicData          `json:"location"`
	User               *UserPublicData              `json:"user"`
	Driver             *DriverPublicData            `json:"driver"`
	Recipient          *UserPublicData              `json:"recipient"`
	Category           *CategoryPublicData          `json:"category"`
	Size               *SizePublicData              `json:"size"`
	TruckType          *TruckTypePublicData         `json:"truck_type"`
	TruckModel         *TruckModelPublicData        `json:"truck_model"`
	DeliveryTime       *DeliveryTimePublicData      `json:"delivery_time"`
	ShipmentContents   []*ShipmentContentPublicData `json:"shipment_contents"`
	ExtraServices      []*ExtraServicePublicData    `json:"extra_services"`
	Destination        *LocationPublicData          `json:"destination"`
	Rating             *Rating                      `json:"rating"`
}

type OrderPaymentMethod string

const (
	OrderCashPaymentMethod OrderPaymentMethod = "cash"
)

type OrderStatus string

const (
	OrderCreatedStatus        OrderStatus = "order_created"
	OrderAcceptedStatus       OrderStatus = "order_accepted"
	PickupInProgressStatus    OrderStatus = "pickup_in_progress"
	ShipmentPickedUpStatus    OrderStatus = "shipment_picked_up"
	InTransitStatus           OrderStatus = "in_transit"
	AtDestinationCityStatus   OrderStatus = "at_destination_city"
	OutForDeliveryStatus      OrderStatus = "out_for_delivery"
	DeliveryAttemptedStatus   OrderStatus = "delivery_attempted"
	DeliveryRescheduledStatus OrderStatus = "delivery_rescheduled"
	ShipmentDeliveredStatus   OrderStatus = "shipment_delivered"
	OrderCompletedStatus      OrderStatus = "order_completed"
	OrderCanceledStatus       OrderStatus = "order_canceled"
	ShipmentReturnedStatus    OrderStatus = "shipment_returned"
)

// AfterFind is a gorm hook that sets the value of the IsDriver field
// based on whether a driver with the same user_id exists
func (o *Order) AfterFind(tx *gorm.DB) (err error) {
	// Rating struct to store the retrieved driver data
	var rating Rating

	// Retrieve the rating data using the order ID
	if err := tx.Table("ratings").Where("order_id = ?", o.ID).First(&rating).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// Handle the error appropriately, excluding the "record not found" error
			return err
		}

		// Handle the "record not found" case separately
		rating = Rating{}
	}

	if rating.ID != 0 {
		o.Rating = &rating
	} else {
		o.Rating = nil
	}

	return
}

// PublicData returns a copy of the order's public information
func (o *Order) PublicData(languageCode string) interface{} {
	locationPublicData := o.Location.PublicData(languageCode).(*LocationPublicData)
	userPublicData := o.User.PublicData(languageCode).(*UserPublicData)

	var driverPublicData *DriverPublicData
	if o.Driver.ID != 0 {
		driverPublicData = o.Driver.PublicData(languageCode).(*DriverPublicData)
	} else {
		driverPublicData = nil
	}

	var recipientPublicData *UserPublicData
	if o.Recipient.ID != 0 {
		recipientPublicData = o.Recipient.PublicData(languageCode).(*UserPublicData)
	} else {
		recipientPublicData = nil
	}

	categoryPublicData := o.Category.PublicData(languageCode).(*CategoryPublicData)

	var sizePublicData *SizePublicData
	if o.Size.ID != 0 {
		sizePublicData = o.Size.PublicData(languageCode).(*SizePublicData)
	} else {
		sizePublicData = nil
	}

	var truckTypePublicData *TruckTypePublicData
	if o.TruckType.ID != 0 {
		truckTypePublicData = o.TruckType.PublicData(languageCode).(*TruckTypePublicData)
	} else {
		truckTypePublicData = nil
	}

	var truckModelPublicData *TruckModelPublicData
	if o.TruckModel.ID != 0 {
		truckModelPublicData = o.TruckModel.PublicData(languageCode).(*TruckModelPublicData)
	} else {
		truckModelPublicData = nil
	}

	deliveryTimePublicData := o.DeliveryTime.PublicData(languageCode).(*DeliveryTimePublicData)

	shipmentContentPublicDataList := make([]*ShipmentContentPublicData, len(o.ShipmentContents))
	for i, shipmentContent := range o.ShipmentContents {
		shipmentContentPublicDataList[i] = shipmentContent.PublicData(languageCode).(*ShipmentContentPublicData)
	}

	extraServicePublicDataList := make([]*ExtraServicePublicData, len(*o.ExtraServices))
	for i, extraService := range *o.ExtraServices {
		extraServicePublicDataList[i] = extraService.PublicData(languageCode).(*ExtraServicePublicData)
	}

	destinationPublicData := o.Destination.PublicData(languageCode).(*LocationPublicData)

	return &OrderPublicData{
		ID:                 o.ID,
		LocationID:         o.LocationID,
		UserID:             o.UserID,
		DriverID:           o.DriverID,
		RecipientID:        o.RecipientID,
		CategoryID:         o.CategoryID,
		SizeID:             o.SizeID,
		TruckTypeID:        o.TruckTypeID,
		TruckModelID:       o.TruckModelID,
		DeliveryTimeID:     o.DeliveryTimeID,
		ShipmentContentIDs: o.ShipmentContentIDs,
		ExtraServiceIDs:    o.ExtraServiceIDs,
		DestinationID:      o.DestinationID,
		Quantity:           o.Quantity,
		Notes:              o.Notes,
		Amount:             o.Amount,
		Latitude:           o.Latitude,
		Longitude:          o.Longitude,
		PaymentMethod:      o.PaymentMethod,
		IsSender:           o.IsSender,
		IsReceiver:         o.IsReceiver,
		Status:             o.Status,
		CreatedAt:          o.CreatedAt,
		Location:           locationPublicData,
		User:               userPublicData,
		Driver:             driverPublicData,
		Recipient:          recipientPublicData,
		Category:           categoryPublicData,
		Size:               sizePublicData,
		TruckType:          truckTypePublicData,
		TruckModel:         truckModelPublicData,
		DeliveryTime:       deliveryTimePublicData,
		ShipmentContents:   shipmentContentPublicDataList,
		ExtraServices:      extraServicePublicDataList,
		Destination:        destinationPublicData,
		Rating:             o.Rating,
	}
}
