package brasilapi

import (
	c "context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type BrasilApiService struct {
	httpClient *resty.Client
}

func NewBrasilApiService() *BrasilApiService {
	return &BrasilApiService{
		httpClient: resty.New().SetTimeout(5 * time.Second),
	}
}

func (b *BrasilApiService) GetAddressByZipCode(ctx c.Context, zipcode string) (*GetAddressByZipCodeRes, error) {
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
