package geoutil

import "math"

// EarthRadiusKm is an approximation of Earth's radius in kilometers.
const EarthRadiusKm = 6371.0

// Haversin returns the haversine of a number.
func Haversin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// CalculateDistance returns the distance (in kilometers) between two points on Earth
// specified by their latitude and longitude (in decimal degrees).
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180.0
	lon1Rad := lon1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0
	lon2Rad := lon2 * math.Pi / 180.0

	// Use Haversine formula
	a := Haversin(lat2Rad-lat1Rad) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*Haversin(lon2Rad-lon1Rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusKm * c
}
