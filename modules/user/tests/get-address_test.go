package tests

import (
	"fmt"
	"net/http"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestGetAddress(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	headers := &utils.TestHeaders

	t.Run("Success", func(t *testing.T) {
		addressID := "adr_01JTX0H1V8N5Q3W7E2R4T6Y8ZXE"
		endpoint := fmt.Sprintf("/v1/internal/user/address/%s", addressID)

		status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, headers)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, addressID, body["id_address"])
		assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfpcKE", body["id_user"])
		assert.Equal(t, "01546022", body["zipcode"])
		assert.Equal(t, "Rua das Flores", body["street"])
		assert.Equal(t, "Belo Horizonte", body["city"])
		assert.Equal(t, "MG", body["state"])
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Address not found", func(t *testing.T) {
			endpoint := "/v1/internal/user/address/adr_01JTX0H1V8N5Q3W7E2R4T6Y8ZZZ"

			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, headers)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Address not found", body["message"])
		})
	})
}
