package entity

import (
	"time"
)

type PhoneVerification struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Phone     string    `gorm:"size:45;not null;unique" json:"phone" validate:"required,e164"`
	Code      string    `gorm:"not null" json:"code"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
