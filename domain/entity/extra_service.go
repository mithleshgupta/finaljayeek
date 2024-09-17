package entity

import (
	"encoding/json"
	"time"

	"github.com/OmarBader7/web-service-jayeek/infrastructure/config"
)

// ExtraService represent a extra service
type ExtraService struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"type:json;not null;" json:"name" validate:"required"`
	Icon      string    `gorm:"size:255;" json:"icon" validate:"required"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:null" json:"updated_at"`
}

// UnmarshalJSON custom unmarshal function for ExtraService
func (c *ExtraService) UnmarshalJSON(data []byte) error {
	type Alias ExtraService
	aux := &struct {
		Name map[string]string `json:"name"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	nameJSON, err := json.Marshal(aux.Name)
	if err != nil {
		return err
	}

	c.Name = string(nameJSON)
	return nil
}

// MarshalJSON custom marshal function for ExtraService
func (c *ExtraService) MarshalJSON() ([]byte, error) {
	type Alias ExtraService
	var nameTranslations map[string]string
	err := json.Unmarshal([]byte(c.Name), &nameTranslations)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		Name map[string]string `json:"name"`
		*Alias
	}{
		Name:  nameTranslations,
		Alias: (*Alias)(c),
	})
}

type ExtraServicePublicData struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

// PublicData returns a copy of the extra service's public information
func (sc *ExtraService) PublicData(languageCode string) interface{} {
	conf := config.NewConfig()

	baseImageURL := conf.BaseStorageURL

	// Get the translated name based on the language code
	var nameTranslations map[string]string
	err := json.Unmarshal([]byte(sc.Name), &nameTranslations)
	if err != nil {
		return nil
	}
	name, ok := nameTranslations[languageCode]
	if !ok {
		name = nameTranslations["en"] // Default to English if the translation is not found
	}

	return &ExtraServicePublicData{
		ID:   sc.ID,
		Name: name,
		Icon: baseImageURL + "/" + sc.Icon,
	}
}
