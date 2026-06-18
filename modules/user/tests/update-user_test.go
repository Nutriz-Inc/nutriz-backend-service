package tests

import (
	"net/http"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUser(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()	
	
	endpoint := "/internal/user/usr_2veL1FPpuXxUaZcFaEC57BfpcKE"

	t.Run("Success", func(t *testing.T) {
		t.Run("Update name", func(t *testing.T) {
			data := dto.UpdateUserReq{
				Name: utils.StringPtr("Maria Silva Atualizada"),
			}
			headers := &utils.TestHeaders
			status, body := fluxgo.RunTestRequest(app, "PATCH", endpoint, data, headers)
			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, "Maria Silva Atualizada", body["name"])
		})

		t.Run("Update email", func(t *testing.T) {
			data := dto.UpdateUserReq{
				Email: utils.StringPtr("maria.nova@email.com"),
			}
			headers := &utils.TestHeaders
			status, body := fluxgo.RunTestRequest(app, "PATCH", endpoint, data, headers)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, "maria.nova@email.com", body["email"])
		})

		t.Run("Admin updates user type", func(t *testing.T) {
			nurseEndpoint := "/internal/user/usr_2veL1FPpuXxUaZcFaEC57BfpxWS"
			data := dto.UpdateUserReq{
				Type: (*entities.EnumUserType)(utils.StringPtr("adm")),
			}
			headers := &utils.TestHeadersAdmin
			status, body := fluxgo.RunTestRequest(app, "PATCH", nurseEndpoint, data, headers)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, "adm", body["type"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("User not found", func(t *testing.T) {
			data := dto.UpdateUserReq{
				Name: utils.StringPtr("Ninguem"),
			}
			headers := &utils.TestHeaders
			status, body := fluxgo.RunTestRequest(app, "PATCH", "/internal/user/usr_2veL1FPpuXxUaZcAbEC57BfpcMM", data, headers)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "User not found", body["message"])
		})

		t.Run("Common user cannot receive internal_identifier", func(t *testing.T) {
			data := dto.UpdateUserReq{
				InternalIdentifier: utils.StringPtr("ID-001"),
			}
			headers := &utils.TestHeaders
			status, body := fluxgo.RunTestRequest(app, "PATCH", endpoint, data, headers)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Common users cannot have internal identifier", body["message"])
			assert.Equal(t, "user.invalid_field", body["code"])
		})

		t.Run("Non-admin cannot change type", func(t *testing.T) {
			data := dto.UpdateUserReq{
				Type: (*entities.EnumUserType)(utils.StringPtr("nurse")),
			}
			headers := &utils.TestHeaders
			status, body := fluxgo.RunTestRequest(app, "PATCH", endpoint, data, headers)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "Only admins can change user type", body["message"])
			assert.Equal(t, "user.forbidden_type", body["code"])
		})

		t.Run("Email already in use", func(t *testing.T) {
			data := dto.UpdateUserReq{
				Email: utils.StringPtr("marta@email.com"),
			}
			headers := &utils.TestHeaders
			status, body := fluxgo.RunTestRequest(app, "PATCH", endpoint, data, headers)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Email already in use", body["message"])
			assert.Equal(t, "user.duplicate_email", body["code"])
		})

		t.Run("Forbidden - different user trying to update", func(t *testing.T) {
			data := dto.UpdateUserReq{
				Name: utils.StringPtr("Hacker"),
			}
			headers := &utils.TestHeadersAdmin
			status, body := fluxgo.RunTestRequest(app, "PATCH", endpoint, data, headers)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "You don't have permission to update this user", body["message"])
			assert.Equal(t, "user.forbidden", body["code"])
		})
	})
}