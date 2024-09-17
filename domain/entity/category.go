package entity

import (
	"encoding/json"
	"time"

	"github.com/OmarBader7/web-service-jayeek/infrastructure/config"
)

// Category represent a geographical category
type Category struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"type:json;not null;" json:"name" validate:"required"`
	MenuName  string    `gorm:"type:json;not null;" json:"menu_name" validate:"required"`
	Icon      string    `gorm:"size:255;" json:"icon" validate:"required"`
	MenuIcon  string    `gorm:"size:255;" json:"menu_icon" validate:"required"`
	IsTruck   *bool     `gorm:"type:boolean;default:null" json:"is_truck"`
	Order     *int64    `gorm:"default:0;size:255;" json:"order"`
	MenuOrder *int64    `gorm:"default:0;size:255;" json:"menu_order"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:null" json:"updated_at"`
}

// UnmarshalJSON custom unmarshal function for Category
func (c *Category) UnmarshalJSON(data []byte) error {
	type Alias Category
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

// MarshalJSON custom marshal function for Category
func (c *Category) MarshalJSON() ([]byte, error) {
	type Alias Category
	var nameTranslations map[string]string
	err := json.Unmarshal([]byte(c.Name), &nameTranslations)
	if err != nil {
		return nil, err
	}

	var menuNameTranslations map[string]string
	err = json.Unmarshal([]byte(c.MenuName), &menuNameTranslations)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		Name     map[string]string `json:"name"`
		MenuName map[string]string `json:"menu_name"`
		*Alias
	}{
		Name:     nameTranslations,
		MenuName: menuNameTranslations,
		Alias:    (*Alias)(c),
	})
}

type CategoryPublicData struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	MenuName  string `json:"menu_name"`
	Icon      string `json:"icon"`
	MenuIcon  string `json:"menu_icon"`
	IsTruck   *bool  `json:"is_truck"`
	Order     *int64 `json:"order"`
	MenuOrder *int64 `json:"menu_order"`
}

// PublicData returns a copy of the category's public information
func (c *Category) PublicData(languageCode string) interface{} {
	conf := config.NewConfig()

	baseImageURL := conf.BaseStorageURL

	// Get the translated name based on the language code
	var nameTranslations map[string]string
	err := json.Unmarshal([]byte(c.Name), &nameTranslations)
	if err != nil {
		return nil
	}
	name, ok := nameTranslations[languageCode]
	if !ok {
		name = nameTranslations["en"] // Default to English if the translation is not found
	}

	// Get the translated name based on the language code
	var menuNameTranslations map[string]string
	err = json.Unmarshal([]byte(c.MenuName), &menuNameTranslations)
	if err != nil {
		return nil
	}
	menuName, ok := menuNameTranslations[languageCode]
	if !ok {
		menuName = menuNameTranslations["en"] // Default to English if the translation is not found
	}

	return &CategoryPublicData{
		ID:        c.ID,
		Name:      name,
		MenuName:  menuName,
		Icon:      baseImageURL + "/" + c.Icon,
		MenuIcon:  baseImageURL + "/" + c.MenuIcon,
		IsTruck:   c.IsTruck,
		Order:     c.Order,
		MenuOrder: c.MenuOrder,
	}
}
