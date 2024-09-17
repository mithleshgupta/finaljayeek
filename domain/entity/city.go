package entity

import (
	"time"
)

// City represents a city with relevant fields.
type City struct {
	ID       uint64 `gorm:"primary_key;auto_increment" json:"id"`
	// Name     string `gorm:"not null;" json:"name"`
	RegionID uint64 `gorm:"not null;" json:"region_id"`
	NameAr   string `gorm:"not null;" json:"name_ar"`
	NameEn   string `gorm:"not null;" json:"name_en"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:null" json:"updated_at"`
}

// CityPublicData defines the public data structure for a city.
type CityPublicData struct {
	ID       uint64 `json:"city_id"`
	// Name     string `json:"name"`
	RegionID uint64 `json:"region_id"`
	NameAr   string `json:"name_ar"`
	NameEn   string `json:"name_en"`
}

// PublicData returns a copy of the city's public information.
func (c *City) PublicData() interface{} {
	return &CityPublicData{
		ID:       c.ID,
		// Name:     c.Name,
		RegionID: c.RegionID,
		NameAr:   c.NameAr,
		NameEn:   c.NameEn,
	}
}
