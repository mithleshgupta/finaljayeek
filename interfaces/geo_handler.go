package interfaces

import (
	"github.com/OmarBader7/web-service-jayeek/infrastructure/geo"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
)

type Geo struct {
	GeoService geo.GeoServiceInterface
}

// NewGeo creates and returns a new instance of Geo.
func NewGeo(geoService geo.GeoServiceInterface) *Geo {
	return &Geo{
		GeoService: geoService,
	}
}

// Iso2 is a handler that retrieves the iso2 code for a given geo location.
func (g *Geo) Iso2(c *gin.Context) {
	// Retrieve the geo location using the GeoService.
	geo, err := g.GeoService.FetchGeoByRequest(c)

	if err != nil {
		// Handle the error if the geo location couldn't be retrieved.
		response.SendInternalServerError(c, ginI18n.MustGetMessage("Failed to retrieve geo location."))
		return
	}

	// Prepare the response data.
	geoData := make(map[string]interface{})
	geoData["iso2"] = geo.ISO2

	// Send the response.
	response.SendOK(c, geoData, "")
}
