package utils

import (
	c "context"
	"fmt"
	"nutriz-backend-service/config"
	"nutriz-backend-service/shared/provider/location"
)

type AddressRes struct {
	State        string
	City         string
	Neighborhood string
	Street       string
	Longitude    *float64
	Latitude     *float64
}

func GetAddressByZipCode(ctx c.Context, zipcode string, config *config.Env) (*AddressRes, error) {
	provider, err := location.NewLocationProvider(config)
	if err != nil {
		return nil, fmt.Errorf("error to initialize location provider: %v", err)
	}

	addressData, err := provider.GetAddressByZipCode(ctx, zipcode)
	if err != nil {
		return nil, fmt.Errorf("error getting address by zipcode: %v", err)
	}

	res := &AddressRes{
		State:        addressData.State,
		City:         addressData.City,
		Neighborhood: addressData.Neighborhood,
		Street:       addressData.Street,
	}

	if addressData.Location != nil && addressData.Location.Coordinates != nil && addressData.Location.Coordinates.Latitude != nil && addressData.Location.Coordinates.Longitude != nil {
		res.Latitude = Float64Ptr(StringToFloat64(*addressData.Location.Coordinates.Latitude))
		res.Longitude = Float64Ptr(StringToFloat64(*addressData.Location.Coordinates.Longitude))

		return res, nil
	}

	query := fmt.Sprintf(
		"%s %s %s Brazil",
		addressData.Street,
		addressData.City,
		addressData.State,
	)

	coordinates, err := provider.GetCoordinatesByAddress(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting coordinates by address: %v", err)
	}

	res.Latitude = Float64Ptr(StringToFloat64(coordinates.Lat))
	res.Longitude = Float64Ptr(StringToFloat64(coordinates.Lon))

	return res, nil
}
