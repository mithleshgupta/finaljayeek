package entity

import (
	"encoding/json"
	"time"
)

// TruckType represent a truck type
type TruckType struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"type:json;not null;" json:"name" validate:"required"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:null" json:"updated_at"`
}

// UnmarshalJSON custom unmarshal function for TruckType
func (tt *TruckType) UnmarshalJSON(data []byte) error {
	type Alias TruckType
	aux := &struct {
		Name map[string]string `json:"name"`
		*Alias
	}{
		Alias: (*Alias)(tt),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	nameJSON, err := json.Marshal(aux.Name)
	if err != nil {
		return err
	}

	tt.Name = string(nameJSON)
	return nil
}

// MarshalJSON custom marshal function for TruckType
func (tt *TruckType) MarshalJSON() ([]byte, error) {
	type Alias TruckType
	var nameTranslations map[string]string

	// Check if tt.Name is empty or not a valid JSON string
	if tt.Name == "" || !json.Valid([]byte(tt.Name)) {
		nameTranslations = make(map[string]string) // Return an empty map or handle it as needed
	} else {
		err := json.Unmarshal([]byte(tt.Name), &nameTranslations)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(&struct {
		Name map[string]string `json:"name"`
		*Alias
	}{
		Name:  nameTranslations,
		Alias: (*Alias)(tt),
	})
}

type TruckTypePublicData struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// PublicData returns a copy of the truck type's public information
func (tt *TruckType) PublicData(languageCode string) interface{} {
	// Get the translated name based on the language code
	var nameTranslations map[string]string
	err := json.Unmarshal([]byte(tt.Name), &nameTranslations)
	if err != nil {
		return nil
	}
	name, ok := nameTranslations[languageCode]
	if !ok {
		name = nameTranslations["en"] // Default to English if the translation is not found
	}

	return &TruckTypePublicData{
		ID:   tt.ID,
		Name: name,
	}
}
