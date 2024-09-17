package entity

import (
	"time"
)

type PasswordReset struct {
	ID               uint64    `gorm:"primaryKey" json:"id"`
	UserID           uint64    `json:"user_id" validate:"required,numeric"`
	VerificationCode string    `gorm:"not null" json:"verification_code"`
	ExpiresAt        time.Time `gorm:"not null" json:"expires_at"`
	Used             bool      `gorm:"default:false" json:"used"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	User             User      `gorm:"foreignKey:UserID" json:"user"`
}
