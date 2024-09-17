package validator

import (
	"github.com/OmarBader7/web-service-jayeek/pkg/language"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/ar"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ar_translations "github.com/go-playground/validator/v10/translations/ar"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/nyaruka/phonenumbers"
)

// ValidationError represents an error returned from the validation process
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var (
	// uni is a universal translator that supports arabic and english languages
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

func init() {
	en := en.New()
	uni = ut.New(en, ar.New())
	validate = validator.New()
}

// RegisterTranslations register the translations for the specified locale
func RegisterTranslations(locale string) error {
	trans, _ := uni.GetTranslator(locale)
	switch trans.Locale() {
	case "ar":
		ar_translations.RegisterDefaultTranslations(validate, trans)
	case "en":
		en_translations.RegisterDefaultTranslations(validate, trans)
	}
	return nil
}

// ValidateExcept validates a struct and returns an array of validation errors, excluding the fields passed as arguments.
func ValidateExcept(c *gin.Context, s interface{}, fields ...string) ([]*ValidationError, error) {
	locale := language.GetLanguage(c)
	if err := RegisterTranslations(locale); err != nil {
		return nil, err
	}
	if err := validate.StructExcept(s, fields...); err != nil {
		errs := GetValidationErrors(err, locale)
		return errs, err
	}
	return nil, nil
}

func ValidatePartial(c *gin.Context, s interface{}, fields ...string) ([]*ValidationError, error) {
	locale := language.GetLanguage(c)
	if err := RegisterTranslations(locale); err != nil {
		return nil, err
	}

	// Register the "phone" tag with the phone number validation function
	validate.RegisterValidation("phone", validatePhone)

	if err := validate.StructPartial(s, fields...); err != nil {
		errs := GetValidationErrors(err, locale)
		return errs, err
	}
	return nil, nil
}

// GetValidationErrors returns an array of validation errors from the validation process
func GetValidationErrors(err error, locale string) (errs []*ValidationError) {
	if err != nil {
		trans, _ := uni.GetTranslator(locale)
		for _, e := range err.(validator.ValidationErrors) {
			var validationError ValidationError
			validationError.Field = e.Field()
			validationError.Message = e.Translate(trans)
			errs = append(errs, &validationError)
		}
	}
	return errs
}

// validatePhone is a custom validation function for the "phone" tag
// on the RecipientPhoneNumber field in the Order struct.
// It checks that the value of RecipientPhoneNumber is a valid phone number.
func validatePhone(fl validator.FieldLevel) bool {
	// Parse the phone number string value from the field using the phonenumbers package
	phoneNumber, err := phonenumbers.Parse(fl.Field().String(), "")

	// Check if the parsing succeeded without error and if the resulting phone number is valid
	// using the IsValidNumber function provided by the phonenumbers package.
	return err == nil && phonenumbers.IsValidNumber(phoneNumber)
}
