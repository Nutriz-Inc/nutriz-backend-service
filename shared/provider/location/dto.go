package location

// BrasilApi
type GetAddressByZipCodeRes struct {
	Cep          string    `json:"cep"`
	State        string    `json:"state"`
	City         string    `json:"city"`
	Neighborhood string    `json:"neighborhood"`
	Street       string    `json:"street"`
	Service      string    `json:"service"`
	Location     *Location `json:"location"`
}

type Location struct {
	Type        string       `json:"type"`
	Coordinates *Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Longitude *string `json:"longitude"`
	Latitude  *string `json:"latitude"`
}

// Nominatim
type GetCoordinatesByAddressRes struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// OSRM
type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type GetOptimizedRouteRes struct {
	Code      string              `json:"code"`
	Waypoints []OptimizedWaypoint `json:"waypoints"`
	Trips     []OptimizedTrip     `json:"trips"`
}

type OptimizedWaypoint struct {
	// WaypointIndex is the position of the input coordinate inside the optimized trip
	WaypointIndex int    `json:"waypoint_index"`
	TripsIndex    int    `json:"trips_index"`
	Name          string `json:"name"`
}

type OptimizedTrip struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
}
