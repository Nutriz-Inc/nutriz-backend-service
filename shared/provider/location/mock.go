package location

import (
	c "context"
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

	return &GetOptimizedRouteRes{
		Code:      "Ok",
		Waypoints: waypoints,
		Trips:     []OptimizedTrip{{Distance: 1000, Duration: 600}},
	}, nil
}
