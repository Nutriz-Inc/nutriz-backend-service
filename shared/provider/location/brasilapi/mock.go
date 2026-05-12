package brasilapi

import c "context"

type BrasilApiMock struct{}

func (m *BrasilApiMock) GetAddressByZipCode(ctx c.Context, zipcode string) (*GetAddressByZipCodeRes, error) {
	mockCepResponse := GetAddressByZipCodeRes{
		Cep:          "01001-000",
		State:        "SP",
		City:         "São Paulo",
		Neighborhood: "Sé",
		Street:       "Praça da Sé",
		Service:      "open-cep",
		Location: Location{
			Type: "Point",
			Coordinates: Coordinates{
				Longitude: "-46.633309",
				Latitude:  "-23.550187",
			},
		},
	}

	return &mockCepResponse, nil
}
