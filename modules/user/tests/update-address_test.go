package tests

import (
	"fmt"
	"net/http"
	httpCode "net/http"
	"nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAddress(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	addressID := "adr_01JTG8K8N4P2R6T9V1X3Y5Z7MIP"
	duplicateZipcode := "01545000"
	endpoint := fmt.Sprintf("/internal/user/address/%s", addressID)

	headers := &utils.TestHeaders

	t.Run("Success", func(t *testing.T) {
		body := dtos.UpdateAddressReq{
			ZipCode:    utils.StringPtr("22250040"),
			Complement: utils.StringPtr("Apto 17"),
			Number:     utils.StringPtr("17"),
		}

		status, resp := fluxgo.RunTestRequest(app, "PUT", endpoint, body, headers)

		assert.Equal(t, httpCode.StatusOK, status)
		assert.Equal(t, addressID, resp["id_address"])
		assert.Equal(t, "22250040", resp["zipcode"])
		assert.NotNil(t, resp["updated_at"])
		assert.NotNil(t, resp["street"])
		assert.NotNil(t, resp["city"])
		assert.NotNil(t, resp["state"])
		assert.NotNil(t, resp["neighborhood"])
		assert.NotNil(t, resp["latitude"])
		assert.NotNil(t, resp["longitude"])
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("No permission", func(t *testing.T) {
			invalidHeader := &utils.TestHeadersAdmin
			body := dtos.UpdateAddressReq{
				ZipCode:    utils.StringPtr("22250040"),
				Complement: utils.StringPtr("Apto 17"),
				Number:     utils.StringPtr("17"),
			}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpoint, body, invalidHeader)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to update address", resp["message"])
		})
		t.Run("Address not found", func(t *testing.T) {
			body := dtos.UpdateAddressReq{
				ZipCode: utils.StringPtr("22250040"),
			}
			route := "/internal/user/address/adr_01JTX0H1V8N5Q3W7E2R4T6Y8ZZZ"

			status, resp := fluxgo.RunTestRequest(app, "PUT", route, body, headers)

			assert.Equal(t, httpCode.StatusNotFound, status)
			assert.Equal(t, "Address not found", resp["message"])
		})

		t.Run("You don't have permission to access this resource", func(t *testing.T) {
			body := dtos.UpdateAddressReq{
				ZipCode: utils.StringPtr("22250040"),
			}
			route := "/internal/user/address/adr_01JTX0H1V8N5Q3W7E2R4T6Y8MBL"

			status, resp := fluxgo.RunTestRequest(app, "PUT", route, body, headers)

			assert.Equal(t, httpCode.StatusForbidden, status)
			assert.Equal(t, "You don't have permission to access this resource", resp["message"])
		})

		t.Run("Address with same zipcode already exists", func(t *testing.T) {
			body := dtos.UpdateAddressReq{
				ZipCode: &duplicateZipcode,
			}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpoint, body, headers)

			assert.Equal(t, httpCode.StatusBadRequest, status)
			assert.Equal(t, "Address with same zipcode already exists", resp["message"])
		})

		t.Run("At least one field must be different to update", func(t *testing.T) {
			body := dtos.UpdateAddressReq{
				ZipCode:    utils.StringPtr("22250040"),
				Complement: utils.StringPtr("Apto 17"),
				Number:     utils.StringPtr("17"),
			}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpoint, body, headers)

			assert.Equal(t, httpCode.StatusBadRequest, status)
			assert.Equal(t, "At least one field must be different to update", resp["message"])
		})
	})
}
