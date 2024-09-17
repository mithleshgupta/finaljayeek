package entity

import (
	"encoding/json"
	"time"
)

// FAQ represent a geographical faq
type FAQ struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Question  string    `gorm:"type:json;not null;" json:"question" validate:"required"`
	Answer    string    `gorm:"type:json;not null;" json:"answer" validate:"required"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:null" json:"updated_at"`
}

// UnmarshalJSON custom unmarshal function for FAQ
func (tm *FAQ) UnmarshalJSON(data []byte) error {
	type Alias FAQ
	aux := &struct {
		Question map[string]string `json:"question"`
		Answer   map[string]string `json:"answer"`
		*Alias
	}{
		Alias: (*Alias)(tm),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	questionJSON, err := json.Marshal(aux.Question)
	if err != nil {
		return err
	}

	tm.Question = string(questionJSON)

	answerJSON, err := json.Marshal(aux.Answer)
	if err != nil {
		return err
	}

	tm.Answer = string(answerJSON)

	return nil
}

func (f *FAQ) MarshalJSON() ([]byte, error) {
	type Alias FAQ
	var questionTranslations map[string]string

	// Check if f.Question is empty or not a valid JSON string
	if f.Question == "" || !json.Valid([]byte(f.Question)) {
		questionTranslations = make(map[string]string) // Return an empty map or handle it as needed
	} else {
		err := json.Unmarshal([]byte(f.Question), &questionTranslations)
		if err != nil {
			return nil, err
		}
	}

	var answerTranslations map[string]string

	// Check if f.Answer is empty or not a valid JSON string
	if f.Answer == "" || !json.Valid([]byte(f.Answer)) {
		answerTranslations = make(map[string]string) // Return an empty map or handle it as needed
	} else {
		err := json.Unmarshal([]byte(f.Answer), &answerTranslations)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(&struct {
		Question map[string]string `json:"question"`
		Answer   map[string]string `json:"answer"`
		*Alias
	}{
		Question: questionTranslations,
		Answer:   answerTranslations,
		Alias:    (*Alias)(f),
	})
}

type FAQPublicData struct {
	ID       uint64 `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// PublicData returns a copy of the faq's public information
func (f *FAQ) PublicData(languageCode string) interface{} {
	// Get the translated question based on the language code
	var questionTranslations map[string]string
	err := json.Unmarshal([]byte(f.Question), &questionTranslations)
	if err != nil {
		return nil
	}
	question, ok := questionTranslations[languageCode]
	if !ok {
		question = questionTranslations["en"] // Default to English if the translation is not found
	}

	// Get the translated answer based on the language code
	var answerTranslations map[string]string
	err = json.Unmarshal([]byte(f.Answer), &answerTranslations)
	if err != nil {
		return nil
	}
	answer, ok := answerTranslations[languageCode]
	if !ok {
		answer = answerTranslations["en"] // Default to English if the translation is not found
	}

	return &FAQPublicData{
		ID:       f.ID,
		Question: question,
		Answer:   answer,
	}
}
