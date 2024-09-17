package entity

import (
	"time"
)

// User represents a user in the system
type Driver struct {
	ID                   uint64             `gorm:"primary_key;auto_increment" json:"id"`
	UserID               uint64             `gorm:"unique;index;" json:"user_id" validate:"required,numeric"`
	TransportationModeID uint64             `gorm:"index;" json:"transportation_mode_id" validate:"required,numeric"`
	Car                  string             `gorm:"size:255;not null;" json:"car" validate:"required"`
	IDNumber             string             `gorm:"unique,size:255;not null;" json:"id_number" validate:"required"`
	Latitude             float64            `gorm:"type:decimal(10,8);not null;" json:"latitude"`
	Longitude            float64            `gorm:"type:decimal(11,8);not null;" json:"longitude"`
	CreatedAt            time.Time          `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            time.Time          `gorm:"default:null" json:"updated_at"`
	Gender               Gender             `gorm:"size:255" json:"gender" validate:"oneof=male female"`
	User                 User               `gorm:"foreignKey:UserID" json:"user"`
	TransportationMode   TransportationMode `gorm:"foreignKey:TransportationModeID" json:"transportation_mode"`
}

type DriverPublicData struct {
	ID                   uint64                        `json:"id"`
	UserID               uint64                        `json:"user_id"`
	TransportationModeID uint64                        `json:"transportation_mode_id"`
	Car                  string                        `json:"car"`
	Latitude             float64                       `json:"latitude"`
	Longitude            float64                       `json:"longitude"`
	Gender               Gender                        `json:"gender"`
	User                 *UserPublicData               `json:"user"`
	TransportationMode   *TransportationModePublicData `json:"transportation_mode"`
}

type Gender string

const (
	MaleGender   Gender = "male"
	FemaleGender Gender = "female"
)

// PublicData returns a copy of the transportation mode's public information
func (d *Driver) PublicData(languageCode string) interface{} {
	userPublicData := d.User.PublicData(languageCode).(*UserPublicData)
	transportationModeData := d.TransportationMode.PublicData(languageCode)

	var transportationModePublicData *TransportationModePublicData
	if transportationModeData != nil {
		transportationModePublicData = transportationModeData.(*TransportationModePublicData)
	} else {
		transportationModePublicData = nil // Or set a default value as needed
	}

	return &DriverPublicData{
		ID:                   d.ID,
		UserID:               d.UserID,
		TransportationModeID: d.TransportationModeID,
		Car:                  d.Car,
		Latitude:             d.Latitude,
		Longitude:            d.Longitude,
		Gender:               d.Gender,
		User:                 userPublicData,
		TransportationMode:   transportationModePublicData,
	}
}
