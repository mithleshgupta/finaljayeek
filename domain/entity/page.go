package entity

import (
	"encoding/json"
	"time"
)

// Page represent a geographical page
type Page struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"type:json;not null;" json:"name" validate:"required"`
	Body      string    `gorm:"type:json;not null;" json:"body" validate:"required"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:null" json:"updated_at"`
}

// UnmarshalJSON custom unmarshal function for Page
func (tm *Page) UnmarshalJSON(data []byte) error {
	type Alias Page
	aux := &struct {
		Name map[string]string `json:"name"`
		Body map[string]string `json:"body"`
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

	bodyJSON, err := json.Marshal(aux.Body)
	if err != nil {
		return err
	}

	tm.Body = string(bodyJSON)

	return nil
}

func (p *Page) MarshalJSON() ([]byte, error) {
	type Alias Page
	var nameTranslations map[string]string

	// Check if p.Name is empty or not a valid JSON string
	if p.Name == "" || !json.Valid([]byte(p.Name)) {
		nameTranslations = make(map[string]string) // Return an empty map or handle it as needed
	} else {
		err := json.Unmarshal([]byte(p.Name), &nameTranslations)
		if err != nil {
			return nil, err
		}
	}

	var bodyTranslations map[string]string

	// Check if p.Body is empty or not a valid JSON string
	if p.Body == "" || !json.Valid([]byte(p.Body)) {
		bodyTranslations = make(map[string]string) // Return an empty map or handle it as needed
	} else {
		err := json.Unmarshal([]byte(p.Body), &bodyTranslations)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(&struct {
		Name map[string]string `json:"name"`
		Body map[string]string `json:"body"`
		*Alias
	}{
		Name:  nameTranslations,
		Body:  bodyTranslations,
		Alias: (*Alias)(p),
	})
}

type PagePublicData struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Body string `json:"body"`
}

// PublicData returns a copy of the page's public information
func (p *Page) PublicData(languageCode string) interface{} {
	// Get the translated name based on the language code
	var nameTranslations map[string]string
	err := json.Unmarshal([]byte(p.Name), &nameTranslations)
	if err != nil {
		return nil
	}
	name, ok := nameTranslations[languageCode]
	if !ok {
		name = nameTranslations["en"] // Default to English if the translation is not found
	}

	// Get the translated body based on the language code
	var bodyTranslations map[string]string
	err = json.Unmarshal([]byte(p.Body), &bodyTranslations)
	if err != nil {
		return nil
	}
	body, ok := bodyTranslations[languageCode]
	if !ok {
		body = bodyTranslations["en"] // Default to English if the translation is not found
	}

	return &PagePublicData{
		ID:   p.ID,
		Name: name,
		Body: body,
	}
}
