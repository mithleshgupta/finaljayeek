package entity

import "time"

// Device represent a device
type Device struct {
	ID         uint64     `gorm:"primary_key;auto_increment" json:"id"`
	UserID     uint64     `gorm:"default:null;index;" json:"user_id" validate:"numeric"`
	FCMToken   string     `gorm:"size:255;not null;unique" json:"fcm_token" validate:"required"`
	DeviceType DeviceType `gorm:"size:100;not null;index;" json:"device_type" validate:"oneof=android ios web"`
	DeviceInfo string     `gorm:"size:255" json:"device_info" validate:"required"`
	IsActive   bool       `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"default:null" json:"updated_at"`
	User       User       `gorm:"foreignKey:UserID" json:"user"`
}

type DevicePublicData struct {
	ID         uint64          `json:"id"`
	UserID     uint64          `json:"user_id"`
	FCMToken   string          `json:"fcm_token"`
	DeviceType DeviceType      `json:"device_type"`
	DeviceInfo string          `json:"device_info"`
	IsActive   bool            `json:"is_active"`
	User       *UserPublicData `json:"user"`
}

type DeviceType string

const (
	Android DeviceType = "android"
	IOS     DeviceType = "ios"
	Web     DeviceType = "web"
)

// PublicData returns a copy of the device's public information
func (d *Device) PublicData(languageCode string) interface{} {
	userPublicData := d.User.PublicData(languageCode).(*UserPublicData)

	return &DevicePublicData{
		ID:         d.ID,
		UserID:     d.UserID,
		FCMToken:   d.FCMToken,
		DeviceType: d.DeviceType,
		DeviceInfo: d.DeviceInfo,
		IsActive:   d.IsActive,
		User:       userPublicData,
	}
}
