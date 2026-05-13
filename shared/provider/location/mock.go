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
