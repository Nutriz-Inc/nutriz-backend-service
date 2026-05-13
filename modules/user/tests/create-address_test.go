package tests

import (
	"fmt"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	"net/http"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestCreateAddress(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/user/address"
	headers := &utils.TestHeaders

	makeBody := func(zipcode string) dto.CreateAddressReq {
		return dto.CreateAddressReq{
			ZipCode:      zipcode,
			Street:       "Rua das Flores",
			Number:       utils.StringPtr("123"),
			City:         "São Paulo",
			State:        "SP",
			Neighborhood: "Centro",
			Complement:   utils.StringPtr("Apto 2"),
		}
	}

	t.Run("Success", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			body := makeBody("01001000")

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, headers)

			assert.Equal(t, http.StatusOK, status)
			assert.NotNil(t, resp["id_address"])
			assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfpcKE", resp["id_user"])
			assert.Equal(t, body.ZipCode, resp["zipcode"])
			assert.Equal(t, body.Street, resp["street"])
			assert.Equal(t, body.City, resp["city"])
			assert.Equal(t, body.State, resp["state"])
			assert.Equal(t, body.Neighborhood, resp["neighborhood"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Address with same zipcode already exists", func(t *testing.T) {
			body := makeBody("09415987")

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, headers)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Address with same zipcode already exists", resp["message"])
		})

		t.Run("User can have up to %d addresses", func(t *testing.T) {
			setupZipcodes := []string{"01002000", "01003000", "01004000"}

			for _, zipcode := range setupZipcodes {
				body := makeBody(zipcode)
				status, _ := fluxgo.RunTestRequest(app, "POST", endpoint, body, headers)
				assert.Equal(t, http.StatusOK, status)
			}

			body := makeBody("01005000")
			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, headers)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, fmt.Sprintf("User can have up to %d addresses", entities.MAX_ADDRESS_QUANTITY_PER_USER), resp["message"])
		})
	})
}
