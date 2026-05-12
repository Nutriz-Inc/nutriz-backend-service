package location

import (
	c "context"
	"nutriz-backend-service/config"
	brasilapi "nutriz-backend-service/shared/provider/location/brasilapi"
)

type LocationProvider interface {
	GetAddressByZipCode(ctx c.Context, zipcode string) (*brasilapi.GetAddressByZipCodeRes, error)
}

func NewLocationProvider(cfg *config.Env) (LocationProvider, error) {
	if cfg.IsTest() {
		return &brasilapi.BrasilApiMock{}, nil
	}

	return brasilapi.NewBrasilApiService(), nil
}
