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
	headers := &utils.TestHeadersCommon

	makeBody := func(zipcode string) dto.CreateAddressReq {
		return dto.CreateAddressReq{
			AddressCreateBase: dto.AddressCreateBase{
				ZipCode:    zipcode,
				Number:     utils.StringPtr("123"),
				Complement: utils.StringPtr("Apto 2"),
			},
		}
	}

	t.Run("Success", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			body := makeBody("01001000")

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, headers)

			assert.Equal(t, http.StatusCreated, status)
			assert.NotNil(t, resp["id_address"])
			assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfpcKE", resp["id_user"])
			assert.Equal(t, body.ZipCode, resp["zipcode"])
			assert.NotNil(t, resp["street"])
			assert.NotNil(t, resp["city"])
			assert.NotNil(t, resp["state"])
			assert.NotNil(t, resp["neighborhood"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("No permission", func(t *testing.T) {
			invalidHeader := &utils.TestHeadersAdmin
			body := makeBody("01001000")

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, invalidHeader)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to create address", resp["message"])
		})

		// t.Run("Address with same zipcode already exists", func(t *testing.T) {
		// 	body := makeBody("09415987")

		// 	status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, headers)

		// 	assert.Equal(t, http.StatusBadRequest, status)
		// 	assert.Equal(t, "Address with same zipcode already exists", resp["message"])
		// })

		t.Run("User can have up to %d addresses", func(t *testing.T) {
			setupZipcodes := []string{"01002000"}

			for _, zipcode := range setupZipcodes {
				body := makeBody(zipcode)
				status, _ := fluxgo.RunTestRequest(app, "POST", endpoint, body, headers)
				assert.Equal(t, http.StatusCreated, status)
			}

			body := makeBody("01005000")
			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, headers)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, fmt.Sprintf("User can have up to %d addresses", entities.MAX_ADDRESS_QUANTITY_PER_USER), resp["message"])
		})
	})
}
