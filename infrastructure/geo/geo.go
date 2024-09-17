package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ctx = context.Background()
)

// GeoServiceInterface defines the methods that a geo service should implement.
type GeoServiceInterface interface {
	FetchGeoByIP(ip string) (*GeoDetails, error)
	FetchGeoByRequest(c *gin.Context) (*GeoDetails, error)
}

// GeoService represents the geo service implementation.
type GeoService struct {
}

// Ensure that GeoService implements GeoServiceInterface.
var _ GeoServiceInterface = &GeoService{}

// NewGeoService creates and returns a new instance of GeoService.
func NewGeoService() *GeoService {
	return &GeoService{}
}

// GeoDetails represents the details of a geo location.
type GeoDetails struct {
	ISO2 string
}

func (s *GeoService) FetchGeoByIP(ip string) (*GeoDetails, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		CountryCode string `json:"countryCode"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &GeoDetails{
		ISO2: data.CountryCode,
	}, nil
}

func (s *GeoService) FetchGeoByRequest(c *gin.Context) (*GeoDetails, error) {
	ip := c.Request.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = strings.Split(c.Request.RemoteAddr, ":")[0]
	}
	ip = "105.71.147.132"
	return s.FetchGeoByIP(ip)
}
