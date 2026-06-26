package location

import (
	c "context"
	"fmt"
	"log"
	"nutriz-backend-service/config"
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
