package entity

import (
	"encoding/json"
	"time"
)

// Size represent a geographical size
type Size struct {
	ID          uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name        string    `gorm:"type:json;not null;" json:"name" validate:"required"`
	Description string    `gorm:"type:json;not null;" json:"description" validate:"required"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:null" json:"updated_at"`
}

// UnmarshalJSON custom unmarshal function for Size
func (s *Size) UnmarshalJSON(data []byte) error {
	type Alias Size
	aux := &struct {
		Name        map[string]string `json:"name"`
		Description map[string]string `json:"description"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	nameJSON, err := json.Marshal(aux.Name)
	if err != nil {
		return err
	}
	s.Name = string(nameJSON)

	descriptionJSON, err := json.Marshal(aux.Description)
	if err != nil {
		return err
	}
	s.Description = string(descriptionJSON)

	return nil
}

// MarshalJSON custom marshal function for Size
func (s *Size) MarshalJSON() ([]byte, error) {
	type Alias Size
	var nameTranslations map[string]string
	err := json.Unmarshal([]byte(s.Name), &nameTranslations)
	if err != nil {
		return nil, err
	}
	var descriptionTranslations map[string]string
	err = json.Unmarshal([]byte(s.Description), &descriptionTranslations)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		Name        map[string]string `json:"name"`
		Description map[string]string `json:"description"`
		*Alias
	}{
		Name:        nameTranslations,
		Description: descriptionTranslations,
		Alias:       (*Alias)(s),
	})
}

type SizePublicData struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PublicData returns a copy of the size's public information
func (s *Size) PublicData(languageCode string) interface{} {
	// Get the translated name based on the language code
	var nameTranslations map[string]string
	err := json.Unmarshal([]byte(s.Name), &nameTranslations)
	if err != nil {
		return nil
	}
	name, ok := nameTranslations[languageCode]
	if !ok {
		name = nameTranslations["en"] // Default to English if the translation is not found
	}

	// Get the translated description based on the language code
	var descriptionTranslations map[string]string
	err = json.Unmarshal([]byte(s.Description), &descriptionTranslations)
	if err != nil {
		return nil
	}
	description, ok := descriptionTranslations[languageCode]
	if !ok {
		description = descriptionTranslations["en"] // Default to English if the translation is not found
	}

	return &SizePublicData{
		ID:          s.ID,
		Name:        name,
		Description: description,
	}
}
