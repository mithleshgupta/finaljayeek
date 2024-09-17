package language

import (
	"github.com/gin-gonic/gin"
)

// GetLanguage returns the user's preferred language based on the Accept-Language header or the "lng" query parameter.
// If the header or parameter is not present, the default language is returned. The default is "en" if not specified.
func GetLanguage(context *gin.Context, defaultLngs ...string) string {
	defaultLng := "en"

	if len(defaultLngs) > 0 {
		defaultLng = defaultLngs[0]
	}

	if context == nil || context.Request == nil {
		return defaultLng
	}

	// Check the Accept-Language header
	lng := context.GetHeader("Accept-Language")
	if lng != "" {
		return lng
	}

	// Check the "lng" query parameter
	lng = context.Query("lng")
	if lng == "" {
		return defaultLng
	}

	return lng
}
