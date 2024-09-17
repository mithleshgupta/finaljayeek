package entity

import (
	"encoding/json"
	"time"

	"github.com/OmarBader7/web-service-jayeek/infrastructure/config"
)

// TransportationMode represent a transportation mode
type TransportationMode struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"type:json;not null;" json:"name" validate:"required"`
	Marker    string    `gorm:"size:255;" json:"marker" validate:"required"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:null" json:"updated_at"`
}

// UnmarshalJSON custom unmarshal function for TransportationMode
func (tm *TransportationMode) UnmarshalJSON(data []byte) error {
	type Alias TransportationMode
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

// MarshalJSON custom marshal function for TransportationMode
func (tm *TransportationMode) MarshalJSON() ([]byte, error) {
	type Alias TransportationMode
	var nameTranslations map[string]string

	// Check if tm.Name is empty or not a valid JSON string
	if tm.Name == "" || !json.Valid([]byte(tm.Name)) {
		nameTranslations = make(map[string]string) // Return an empty map or handle it as needed
	} else {
		err := json.Unmarshal([]byte(tm.Name), &nameTranslations)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(&struct {
		Name map[string]string `json:"name"`
		*Alias
	}{
		Name:  nameTranslations,
		Alias: (*Alias)(tm),
	})
}

type TransportationModePublicData struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Marker string `json:"marker"`
}

// PublicData returns a copy of the transportation mode's public information
func (tm *TransportationMode) PublicData(languageCode string) interface{} {
	conf := config.NewConfig()

	baseImageURL := conf.BaseStorageURL

	// Get the translated name based on the language code
	var nameTranslations map[string]string
	err := json.Unmarshal([]byte(tm.Name), &nameTranslations)
	if err != nil {
		return nil
	}
	name, ok := nameTranslations[languageCode]
	if !ok {
		name = nameTranslations["en"] // Default to English if the translation is not found
	}

	return &TransportationModePublicData{
		ID:     tm.ID,
		Name:   name,
		Marker: baseImageURL + "/" + tm.Marker,
	}
}
