package tests

import (
	"fmt"
	"net/http"
	authDto "nutriz-backend-service/modules/auth/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestCreateDonation(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/donation"
	headers := &utils.TestHeaders

	t.Run("Success", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, nil, headers)

			assert.Equal(t, http.StatusCreated, status)
			assert.NotEmpty(t, body["id_donation"])
			assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfpcKE", body["created_by"])
			assert.Equal(t, true, body["is_active"])
			assert.Equal(t, false, body["is_recurrent"], "user has no valid blood exam")
		})

		t.Run("Recurrent donation when the user has a valid blood exam", func(t *testing.T) {
			loginData := authDto.LoginReq{
				Email:    "renata@email.com",
				Password: "12345678",
			}

			status, loginBody := fluxgo.RunTestRequest(app, "POST", "/public/auth/login", loginData, nil)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, true, loginBody["is_recurrent_donor"])

			token, ok := loginBody["token"].(string)
			assert.True(t, ok)

			recurrentHeaders := fluxgo.Headers{"Authorization": "Bearer " + token}

			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, nil, &recurrentHeaders)

			assert.Equal(t, http.StatusCreated, status)
			assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfpcBR", body["created_by"])
			assert.Equal(t, true, body["is_recurrent"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("No permission", func(t *testing.T) {
			invalidHeader := &utils.TestHeadersAdmin
			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, nil, invalidHeader)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to create donation", body["message"])
		})
		t.Run("User can have up to %d active donations", func(t *testing.T) {
			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, nil, headers)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(
				t,
				fmt.Sprintf("User can have up to %d active donations", entities.MAX_DONATION_ACTIVE_QUANTITY_PER_USER),
				body["message"],
			)
		})
	})
}
