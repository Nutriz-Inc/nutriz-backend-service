package location

import (
	c "context"
	"math"
)

type LocationMock struct{}

func (m *LocationMock) GetAddressByZipCode(ctx c.Context, zipcode string) (*GetAddressByZipCodeRes, error) {
	mockCepResponse := GetAddressByZipCodeRes{
		Cep:          "01001-000",
		State:        "SP",
		City:         "São Paulo",
		Neighborhood: "Sé",
		Street:       "Praça da Sé",
		Service:      "open-cep",
		Location: &Location{
			Type:        "Point",
			Coordinates: &Coordinates{},
		},
	}

	return &mockCepResponse, nil
}

func (m *LocationMock) GetCoordinatesByAddress(ctx c.Context, address string) (*GetCoordinatesByAddressRes, error) {
	mockNominatimResponse := []GetCoordinatesByAddressRes{
		{
			Lat: "-23.550187",
			Lon: "-46.633309",
		},
	}

	return &mockNominatimResponse[0], nil
}

const mockAverageSpeedKmh = 40.0

func (m *LocationMock) GetOptimizedRoute(ctx c.Context, coordinates []Coordinate) (*GetOptimizedRouteRes, error) {
	if len(coordinates) == 0 {
		return nil, nil
	}

	waypoints := make([]OptimizedWaypoint, 0, len(coordinates))
	for index := range coordinates {
		waypoints = append(waypoints, OptimizedWaypoint{
			WaypointIndex: index,
			TripsIndex:    0,
			Name:          "Mock street",
		})
	}

	if len(waypoints) > 1 {
		waypoints[0].WaypointIndex, waypoints[1].WaypointIndex = waypoints[1].WaypointIndex, waypoints[0].WaypointIndex
	}

	distance := 0.0
	for index := 1; index < len(coordinates); index++ {
		distance += haversineMeters(coordinates[index-1], coordinates[index])
	}

	duration := distance / (mockAverageSpeedKmh * 1000 / 3600)

	return &GetOptimizedRouteRes{
		Code:      "Ok",
		Waypoints: waypoints,
		Trips:     []OptimizedTrip{{Distance: distance, Duration: duration}},
	}, nil
}

func haversineMeters(from, to Coordinate) float64 {
	const earthRadiusMeters = 6371000.0

	fromLatitude := from.Latitude * math.Pi / 180
	toLatitude := to.Latitude * math.Pi / 180
	deltaLatitude := (to.Latitude - from.Latitude) * math.Pi / 180
	deltaLongitude := (to.Longitude - from.Longitude) * math.Pi / 180

	a := math.Sin(deltaLatitude/2)*math.Sin(deltaLatitude/2) +
		math.Cos(fromLatitude)*math.Cos(toLatitude)*math.Sin(deltaLongitude/2)*math.Sin(deltaLongitude/2)

	return earthRadiusMeters * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
