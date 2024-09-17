package entity

import (
	"encoding/json"
	"time"
)

// Location represent a geographical location
type Location struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"type:json;not null;" json:"name" validate:"required"`
	Latitude  float64   `gorm:"type:decimal(10,8);not null;" json:"latitude" validate:"required"`
	Longitude float64   `gorm:"type:decimal(11,8);not null;" json:"longitude" validate:"required"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:null" json:"updated_at"`
}

// UnmarshalJSON custom unmarshal function for Location
func (tm *Location) UnmarshalJSON(data []byte) error {
	type Alias Location
	aux := &struct {
		Name map[string]string `json:"name"`
		*Alias
	}{
		Alias: (*Alias)(tm),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	nameJSON, err := json.Marshal(aux.Name)
	if err != nil {
		return err
	}

	tm.Name = string(nameJSON)
	return nil
}

func (l *Location) MarshalJSON() ([]byte, error) {
	type Alias Location
	var nameTranslations map[string]string

	// Check if l.Name is empty or not a valid JSON string
	if l.Name == "" || !json.Valid([]byte(l.Name)) {
		nameTranslations = make(map[string]string) // Return an empty map or handle it as needed
	} else {
		err := json.Unmarshal([]byte(l.Name), &nameTranslations)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(&struct {
		Name map[string]string `json:"name"`
		*Alias
	}{
		Name:  nameTranslations,
		Alias: (*Alias)(l),
	})
}

type LocationPublicData struct {
	ID        uint64  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// PublicData returns a copy of the location's public information
func (l *Location) PublicData(languageCode string) interface{} {
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

	return &LocationPublicData{
		ID:        l.ID,
		Name:      name,
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
	}
}
