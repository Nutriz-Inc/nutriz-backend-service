package utils

import (
	c "context"
	"fmt"
	"nutriz-backend-service/config"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/provider/location"
	"strings"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
)

type OptimizedRoute struct {
	StopOrders []int
	Duration   time.Duration
}

func GetOptimizedRoute(ctx c.Context, coordinates []location.Coordinate, config *config.Env) (*OptimizedRoute, error) {
	if len(coordinates) == 0 {
		return nil, fmt.Errorf("no coordinates provided")
	}

	if len(coordinates) == 1 {
		return &OptimizedRoute{StopOrders: []int{0}, Duration: 0}, nil
	}

	provider, err := location.NewLocationProvider(config)
	if err != nil {
		return nil, fmt.Errorf("error to initialize location provider: %v", err)
	}

	optimizedRoute, err := provider.GetOptimizedRoute(ctx, coordinates)
	if err != nil {
		return nil, fmt.Errorf("error getting optimized route: %v", err)
	}
	if optimizedRoute == nil || len(optimizedRoute.Waypoints) != len(coordinates) {
		return nil, fmt.Errorf("invalid optimized route response")
	}
	if len(optimizedRoute.Trips) == 0 {
		return nil, fmt.Errorf("optimized route without trips")
	}

	stopOrders := make([]int, len(coordinates))
	used := make(map[int]bool, len(coordinates))

	for index, waypoint := range optimizedRoute.Waypoints {
		if waypoint.WaypointIndex < 0 || waypoint.WaypointIndex >= len(coordinates) {
			return nil, fmt.Errorf("invalid waypoint index %d", waypoint.WaypointIndex)
		}
		if used[waypoint.WaypointIndex] {
			return nil, fmt.Errorf("duplicated waypoint index %d", waypoint.WaypointIndex)
		}

		used[waypoint.WaypointIndex] = true
		stopOrders[index] = waypoint.WaypointIndex
	}

	return &OptimizedRoute{
		StopOrders: stopOrders,
		Duration:   time.Duration(optimizedRoute.Trips[0].Duration * float64(time.Second)),
	}, nil
}

type StopCoordinates struct {
	Latitude  *float64
	Longitude *float64
}

func BuildOptimizedStopOrders(
	ctx c.Context,
	stops []StopCoordinates,
	config *config.Env,
) ([]int16, time.Duration, error) {
	if len(stops) == 0 {
		return []int16{}, 0, nil
	}

	coordinates := make([]location.Coordinate, 0, len(stops))
	for _, stop := range stops {
		latitude, longitude := FillMissingCoordinates(stop.Latitude, stop.Longitude)

		coordinates = append(coordinates, location.Coordinate{
			Latitude:  latitude,
			Longitude: longitude,
		})
	}

	optimizedRoute, err := GetOptimizedRoute(ctx, coordinates, config)
	if err != nil {
		return nil, 0, err
	}

	stopOrders := make([]int16, 0, len(optimizedRoute.StopOrders))
	for _, position := range optimizedRoute.StopOrders {
		stopOrders = append(stopOrders, int16(position))
	}

	return stopOrders, optimizedRoute.Duration, nil
}

func TotalRouteDuration(drivingTime time.Duration, stopCount int) time.Duration {
	return drivingTime + time.Duration(stopCount)*entities.ROUTE_STOP_SAFETY_TIME
}

func OptimizeStops(
	ctx c.Context,
	env *config.Env,
	stops []StopCoordinates,
) ([]int16, time.Duration, *fluxgo.GlobalError) {
	stopOrders, drivingTime, err := BuildOptimizedStopOrders(ctx, stops, env)
	if err != nil {
		return nil, 0, fluxgo.ErrorInternalError("Error to build the route: " + err.Error())
	}

	totalDuration := TotalRouteDuration(drivingTime, len(stops))

	if totalDuration > entities.MAX_ROUTE_DURATION {
		return nil, 0, fluxgo.ErrorBadRequest(
			fmt.Sprintf(
				"Route takes %.1f hours and the maximum allowed is %.0f hours",
				totalDuration.Hours(),
				entities.MAX_ROUTE_DURATION.Hours(),
			),
			"route.max_duration_exceeded",
		)
	}

	return stopOrders, totalDuration, nil
}

func MatchesAddressField(value *string, expected string) bool {
	if value == nil {
		return false
	}

	return strings.EqualFold(strings.TrimSpace(*value), strings.TrimSpace(expected))
}
