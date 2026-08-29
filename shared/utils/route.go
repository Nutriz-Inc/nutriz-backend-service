package utils

import (
	c "context"
	"fmt"
	"nutriz-backend-service/config"
	"nutriz-backend-service/shared/provider/location"
	"time"
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
