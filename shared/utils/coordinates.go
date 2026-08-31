package utils

import (
	"math"
	"math/rand"
)

// Bounding box of the urban area of São Paulo city.
const (
	SAO_PAULO_MIN_LATITUDE  = -23.720
	SAO_PAULO_MAX_LATITUDE  = -23.400
	SAO_PAULO_MIN_LONGITUDE = -46.800
	SAO_PAULO_MAX_LONGITUDE = -46.400
)

const COORDINATE_PRECISION = 6

func GenerateSaoPauloCoordinates() (float64, float64) {
	latitude := randomCoordinate(SAO_PAULO_MIN_LATITUDE, SAO_PAULO_MAX_LATITUDE)
	longitude := randomCoordinate(SAO_PAULO_MIN_LONGITUDE, SAO_PAULO_MAX_LONGITUDE)

	return latitude, longitude
}

func FillMissingCoordinates(latitude, longitude *float64) (float64, float64) {
	if latitude != nil && longitude != nil {
		return *latitude, *longitude
	}

	return GenerateSaoPauloCoordinates()
}

func randomCoordinate(min, max float64) float64 {
	value := min + rand.Float64()*(max-min)
	factor := math.Pow(10, COORDINATE_PRECISION)

	return math.Round(value*factor) / factor
}
