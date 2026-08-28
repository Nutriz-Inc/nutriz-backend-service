package location

import (
	c "context"
	"fmt"
	"log"
	"nutriz-backend-service/config"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type LocationService struct {
	httpClient *resty.Client
	cfg        *config.Env
}

func NewLocationService(cfg *config.Env) *LocationService {
	return &LocationService{
		httpClient: resty.New().SetTimeout(10 * time.Second),
		cfg:        cfg,
	}
}

// BrasilApi
func (b *LocationService) GetAddressByZipCode(ctx c.Context, zipcode string) (*GetAddressByZipCodeRes, error) {
	var resp GetAddressByZipCodeRes
	endpoint := fmt.Sprintf("https://brasilapi.com.br/api/cep/v2/%s", zipcode)

	res, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(&resp).
		Get(endpoint)

	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf(
			"brasil api returned status %d",
			res.StatusCode(),
		)
	}

	return &resp, nil
}

// Nominatim
func (n *LocationService) GetCoordinatesByAddress(ctx c.Context, address string) (*GetCoordinatesByAddressRes, error) {
	var resp []GetCoordinatesByAddressRes

	userAgent := fmt.Sprintf("%s/%s (%s)", n.cfg.Service.Name, n.cfg.Service.Version, n.cfg.Service.Email)

	res, err := n.httpClient.R().
		SetContext(ctx).
		SetResult(&resp).
		SetQueryParams(map[string]string{
			"q":      address,
			"format": "json",
			"limit":  "1",
		}).
		SetHeader("User-Agent", userAgent).
		Get("https://nominatim.openstreetmap.org/search")

	if err != nil {
		return nil, err
	}

	log.Printf("[nominatim] user-agent=%q query=%q status=%d body=%s", userAgent, address, res.StatusCode(), res.Body())

	if res.IsError() || len(resp) == 0 {
		return nil, fmt.Errorf("no coordinates found")
	}

	return &resp[0], nil
}

// Osrm
const OSRM_MAX_COORDINATES = 100

func (o *LocationService) GetOptimizedRoute(ctx c.Context, coordinates []Coordinate) (*GetOptimizedRouteRes, error) {
	if len(coordinates) == 0 {
		return nil, fmt.Errorf("no coordinates provided")
	}
	if len(coordinates) > OSRM_MAX_COORDINATES {
		return nil, fmt.Errorf("too many coordinates, maximum is %d", OSRM_MAX_COORDINATES)
	}

	points := make([]string, 0, len(coordinates))
	for _, coordinate := range coordinates {
		points = append(points, fmt.Sprintf("%f,%f", coordinate.Longitude, coordinate.Latitude))
	}

	var resp GetOptimizedRouteRes

	endpoint := fmt.Sprintf(
		"https://router.project-osrm.org/trip/v1/driving/%s",
		strings.Join(points, ";"),
	)

	res, err := o.httpClient.R().
		SetContext(ctx).
		SetResult(&resp).
		SetQueryParams(map[string]string{
			"source":    "first",
			"roundtrip": "false",
			"overview":  "false",
		}).
		Get(endpoint)

	if err != nil {
		return nil, err
	}

	if res.IsError() || resp.Code != "Ok" {
		return nil, fmt.Errorf("osrm returned status %d and code %q", res.StatusCode(), resp.Code)
	}

	if len(resp.Waypoints) != len(coordinates) {
		return nil, fmt.Errorf("osrm returned %d waypoints for %d coordinates", len(resp.Waypoints), len(coordinates))
	}

	return &resp, nil
}
