package tests

import (
	"fmt"
	"net/http"
	"nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserBaby(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	userBabyID := "usb_01JTG8J5F6W9K2M4P7Q1X8Y3ZAL"
	endpoint := fmt.Sprintf("/v1/internal/user/baby/%s", userBabyID)
	headers := &utils.TestHeaders

	t.Run("Success", func(t *testing.T) {
		body := dtos.UpdateUserBabyReq{
			Name:      utils.StringPtr("Miguel Atualizado"),
			BirthDate: utils.StringPtr("2025-01-15"),
		}
		status, resp := fluxgo.RunTestRequest(app, "PUT", endpoint, body, headers)
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, userBabyID, resp["id_user_baby"])
		assert.Equal(t, "Miguel Atualizado", resp["name"])
		assert.NotNil(t, resp["updated_at"])
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("No permission", func(t *testing.T) {
			invalidHeader := &utils.TestHeadersAdmin
			body := dtos.UpdateUserBabyReq{
				Name: utils.StringPtr("Miguel Atualizado"),
			}
			status, resp := fluxgo.RunTestRequest(app, "PUT", endpoint, body, invalidHeader)
			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to update baby", resp["message"])
		})

		t.Run("Baby not found", func(t *testing.T) {
			body := dtos.UpdateUserBabyReq{
				Name: utils.StringPtr("Miguel Atualizado"),
			}
			route := "/v1/internal/user/baby/usb_01JTG8J5F6W9K2M4P7Q1X8Y3ZAE"
			status, resp := fluxgo.RunTestRequest(app, "PUT", route, body, headers)
			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "User baby not found", resp["message"])
		})

		t.Run("You don't have permission to access this resource", func(t *testing.T) {
			body := dtos.UpdateUserBabyReq{
				Name: utils.StringPtr("Miguel Atualizado"),
			}
			route := fmt.Sprintf("/v1/internal/user/baby/%s", "usb_01JTG8K8N4P2R6T9V1X3Y5Z7MIP")
			status, resp := fluxgo.RunTestRequest(app, "PUT", route, body, headers)
			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "You don't have permission to access this resource", resp["message"])
		})

		t.Run("Birthdate in the future", func(t *testing.T) {
			body := dtos.UpdateUserBabyReq{
				BirthDate: utils.StringPtr("2099-01-01"),
			}
			status, resp := fluxgo.RunTestRequest(app, "PUT", endpoint, body, headers)
			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Birth date cannot be in the future", resp["message"])
		})

		t.Run("At least one field must be different to update", func(t *testing.T) {
			body := dtos.UpdateUserBabyReq{
				Name:      utils.StringPtr("Miguel Atualizado"),
				BirthDate: utils.StringPtr("2025-01-15"),
			}
			status, resp := fluxgo.RunTestRequest(app, "PUT", endpoint, body, headers)
			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "At least one field must be different to update", resp["message"])
		})
	})
}
