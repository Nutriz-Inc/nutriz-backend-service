package location

import (
	c "context"
	"nutriz-backend-service/config"
)

type LocationProvider interface {
	//BrasilApi
	GetAddressByZipCode(ctx c.Context, zipcode string) (*GetAddressByZipCodeRes, error)
	//Nominatim
	GetCoordinatesByAddress(ctx c.Context, address string) (*GetCoordinatesByAddressRes, error)
}

func NewLocationProvider(cfg *config.Env) (LocationProvider, error) {
	if cfg.IsTest() {
		return &LocationMock{}, nil
	}

	return NewLocationService(cfg), nil
}
