package location

import (
	c "context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type LocationService struct {
	httpClient *resty.Client
}

func NewLocationService() *LocationService {
	return &LocationService{
		httpClient: resty.New().SetTimeout(5 * time.Second),
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

	res, err := n.httpClient.R().
		SetContext(ctx).
		SetResult(&resp).
		SetQueryParams(map[string]string{
			"q":      address,
			"format": "json",
			"limit":  "1",
		}).
		SetHeader("User-Agent", "location-service/1.0").
		Get("https://nominatim.openstreetmap.org/search")

	if err != nil {
		return nil, err
	}

	if res.IsError() || len(resp) == 0 {
		return nil, fmt.Errorf("no coordinates found")
	}

	return &resp[0], nil
}
