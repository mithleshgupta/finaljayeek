package entity

// User represents a setting in the system
type Setting struct {
	Key       string    `gorm:"unique,size:255;not null;" json:"key" validate:"required"`
	Value     string    `json:"value"`
}

type SettingPublicData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// PublicData returns a copy of the setting's public information
func (d *Setting) PublicData() interface{} {
	return &SettingPublicData{
		Key:   d.Key,
		Value: d.Value,
	}
}
