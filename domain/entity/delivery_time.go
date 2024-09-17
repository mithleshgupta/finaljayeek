package entity

import (
	"encoding/json"
	"time"
)

// DeliveryTime represent a delivery time
type DeliveryTime struct {
	ID        uint64        `gorm:"primary_key;auto_increment" json:"id"`
	Name      string        `gorm:"type:json;not null;" json:"name" validate:"required"`
	Duration  time.Duration `gorm:"type:bigint;not null;index;" json:"duration" validate:"required"`
	CreatedAt time.Time     `gorm:"default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt time.Time     `gorm:"default:null;" json:"updated_at"`
}

// UnmarshalJSON custom unmarshal function for DeliveryTime
func (dt *DeliveryTime) UnmarshalJSON(data []byte) error {
	type Alias DeliveryTime
	aux := &struct {
		Name map[string]string `json:"name"`
		*Alias
	}{
		Alias: (*Alias)(dt),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	nameJSON, err := json.Marshal(aux.Name)
	if err != nil {
		return err
	}

	dt.Name = string(nameJSON)
	return nil
}

// MarshalJSON custom marshal function for DeliveryTime
func (dt *DeliveryTime) MarshalJSON() ([]byte, error) {
	type Alias DeliveryTime
	var nameTranslations map[string]string
	err := json.Unmarshal([]byte(dt.Name), &nameTranslations)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		Name map[string]string `json:"name"`
		*Alias
	}{
		Name:  nameTranslations,
		Alias: (*Alias)(dt),
	})
}

type DeliveryTimePublicData struct {
	ID       uint64        `json:"id"`
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
}

// PublicData returns a copy of the delivery time's public information
func (l *DeliveryTime) PublicData(languageCode string) interface{} {
	// Get the translated name based on the language code
	var nameTranslations map[string]string
	err := json.Unmarshal([]byte(l.Name), &nameTranslations)
	if err != nil {
		return nil
	}
	name, ok := nameTranslations[languageCode]
	if !ok {
		name = nameTranslations["en"] // Default to English if the translation is not found
	}

	return &DeliveryTimePublicData{
		ID:       l.ID,
		Name:     name,
		Duration: l.Duration,
	}
}
